package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/usecases"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/nerdyworm/rsq"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

// App is the main application for erraroo
type App struct {
	Router     *mux.Router
	JobRouter  *rsq.JobRouter
	HTTPServer http.Handler
}

func New() *App {
	a := &App{}
	a.setupMux()
	a.setupMiddleware()
	a.setupQueue()
	return a
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.HTTPServer.ServeHTTP(w, r)
}

func (a *App) Run(job *rsq.Job) error {
	return a.JobRouter.Run(job)
}

func (a *App) setupMiddleware() {
	a.HTTPServer = alice.New(
		api.LoggingMiddleware,
		api.CorsMiddleware,
		api.GzipMiddleware,
	).Then(a.Router)
}

func (a *App) setupMux() {
	a.Router = mux.NewRouter()
	a.Router.NotFoundHandler = a.Handler(api.NotFoundHandler)
	a.Router.Handle("/api/v1/signups", a.Handler(api.SignupsCreate)).Methods("POST")
	a.Router.Handle("/api/v1/sessions", a.Handler(api.SessionsCreate)).Methods("POST")
	a.Router.Handle("/api/v1/sessions", a.Handler(api.SessionsDestroy)).Methods("DELETE")
	a.Router.Handle("/api/v1/me", a.AuthroziedHandler(api.MeHandler)).Methods("GET")
	a.Router.Handle("/api/v1/events", a.Handler(api.EventsCreate)).Methods("POST")
	a.Router.Handle("/api/v1/projects", a.AuthroziedHandler(api.ProjectsIndex)).Methods("GET")
	a.Router.Handle("/api/v1/projects", a.AuthroziedHandler(api.ProjectsCreate)).Methods("POST")
	a.Router.Handle("/api/v1/projects/{id}", a.AuthroziedHandler(api.ProjectsShow)).Methods("GET")
	a.Router.Handle("/api/v1/projects/{id}", a.AuthroziedHandler(api.ProjectsUpdate)).Methods("PUT")
	a.Router.Handle("/api/v1/projects/{id}", a.AuthroziedHandler(api.ProjectsDelete)).Methods("DELETE")
	a.Router.Handle("/api/v1/projects/{id}/regenerate-token", a.AuthroziedHandler(api.ProjectsRegenerateToken)).Methods("POST")
	a.Router.Handle("/api/v1/events/{id}", a.AuthroziedHandler(api.EventsShow)).Methods("GET")
	a.Router.Handle("/api/v1/events", a.AuthroziedHandler(api.EventsIndex)).Methods("GET")
	a.Router.Handle("/api/v1/errors", a.AuthroziedHandler(api.ErrorsIndex)).Methods("GET")
	a.Router.Handle("/api/v1/errors/{id}", a.AuthroziedHandler(api.ErrorsShow)).Methods("GET")
	a.Router.Handle("/api/v1/errors/{id}", a.AuthroziedHandler(api.ErrorsUpdate)).Methods("PUT")
	a.Router.Handle("/api/v1/invitations", a.AuthroziedHandler(api.InvitationsIndex)).Methods("GET")
	a.Router.Handle("/api/v1/invitations", a.AuthroziedHandler(api.InvitationsCreate)).Methods("POST")
	a.Router.Handle("/api/v1/invitations/{token}", a.Handler(api.InvitationsShow)).Methods("GET")
	a.Router.Handle("/api/v1/invitations/{token}", a.AuthroziedHandler(api.InvitationsDelete)).Methods("DELETE")
	a.Router.Handle("/api/v1/users/{id}", a.AuthroziedHandler(api.UsersShow)).Methods("GET")
	a.Router.Handle("/api/v1/prefs/{id}", a.AuthroziedHandler(api.PrefsUpdate)).Methods("PUT")
	a.Router.Handle("/api/v1/passwords", a.AuthroziedHandler(api.PasswordsCreate)).Methods("POST")
	a.Router.Handle("/api/v1/passwordRecovers", a.Handler(api.PasswordRecoversCreate)).Methods("POST")
	a.Router.Handle("/api/v1/passwordRecovers/{token}", a.Handler(api.PasswordRecoversUpdate)).Methods("PUT")
	a.Router.Handle("/api/v1/passwordRecovers/{token}", a.Handler(api.PasswordRecoversShow)).Methods("GET")
	a.Router.Handle("/api/v1/timings", a.AuthroziedHandler(api.TimingsIndex)).Methods("GET")
	a.Router.Handle("/api/v1/backlog", a.AuthroziedHandler(api.Backlog)).Methods("POST")
	a.Router.HandleFunc("/healthcheck", api.Healthcheck).Methods("GET")
}

func (a *App) setupQueue() {
	a.JobRouter = rsq.NewJobRouter()
	a.JobRouter.Handle("invitation.deliver", a.JobHandler(func(job *rsq.Job, c *cx.Context) error {
		var token string
		err := json.Unmarshal(job.Payload, &token)
		if err != nil {
			return err
		}

		return usecases.InvitationDeliver(token)
	}))

	a.JobRouter.Handle("create.js.error", a.JobHandler(func(job *rsq.Job, c *cx.Context) error {
		payload := map[string]interface{}{}
		err := json.Unmarshal(job.Payload, &payload)
		if err != nil {
			return err
		}

		id := int64(payload["projectID"].(float64))
		project, err := models.Projects.FindByID(id)
		if err != nil {
			return err
		}

		raw, err := json.Marshal(payload["data"])
		if err != nil {
			logger.Error("marshalling payload", "payload", payload)
			return err
		}

		e := models.NewEvent(project, "js.error", string(raw))
		err = models.Events.Insert(e)
		if err != nil {
			logger.Error("inserting event", "err", err)
			return err
		}

		return usecases.AfterErrorEventCreated(e)
	}))

}

func (a *App) newContext(r *http.Request) (*cx.Context, error) {
	var err error
	ctx := &cx.Context{}

	if r != nil {
		id, err := getCurrentUserID(r)
		if err == ErrInvalidToken {
			return ctx, err
		}

		if err != nil {
			logger.Error("getting current user", "err", err)
			return ctx, err
		}

		if id == 0 {
			return ctx, err
		}

		ctx.User, err = models.Users.FindByID(id)
		if err != nil {
			logger.Error("finding user from token", "err", err)
			return ctx, err
		}
	}

	return ctx, err
}

func getCurrentUserID(r *http.Request) (int64, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return 0, nil
	}

	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return config.TokenSigningKey, nil
	})

	if err != nil {
		return 0, ErrInvalidToken
	}

	id := token.Claims["user_id"].(float64)
	return int64(id), nil
}

// SetupForTesting sets the app up to be tested
func (a *App) SetupForTesting() {
	models.SetupForTesting()
	models.Migrate()
	jobs.Use(rsq.NewMemoryAdapter())
}

// Shutdown shutsdown all the things
func (a *App) Shutdown() {
	defer models.Shutdown()

	err := jobs.Shutdown()
	if err != nil {
		logger.Error("jobs.Shutdown()", "err", err)
	}
}

// AppHandler is the fn we use
type AppHandler func(http.ResponseWriter, *http.Request, *cx.Context) error

func (a *App) Handler(fn AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := a.newContext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = fn(w, r, ctx)
		if err == nil {
			return
		}

		switch err.(type) {
		case models.ValidationErrors:
			api.JSON(w, http.StatusBadRequest, err)
		default:
			handleError(err, w)
		}
	})
}

func (a *App) AuthroziedHandler(fn AppHandler) http.Handler {
	return a.Handler(func(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
		if ctx.User == nil {
			return cx.ErrLoginRequired
		}

		return fn(w, r, ctx)
	})
}

type AppJobHandler func(*rsq.Job, *cx.Context) error

func (a *App) JobHandler(fn AppJobHandler) rsq.JobHandlerFunc {
	return func(job *rsq.Job) error {
		start := time.Now()

		ctx, err := a.newContext(nil)
		if err != nil {
			return err
		}

		err = fn(job, ctx)
		if err != nil {
			logger.Error(err.Error(), "name", job.Name, "payload", fmt.Sprintf("%s", job.Payload), "runtime", time.Since(start))
		} else {
			//logger.Info("ran", "name", job.Name, "payload", fmt.Sprintf("%s", job.Payload), "runtime", time.Since(start))
			logger.Info("ran", "name", job.Name, "runtime", time.Since(start))
		}

		return err
	}
}

func handleError(err error, w http.ResponseWriter) {
	switch err {
	case models.ErrNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case cx.ErrLoginRequired:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case ErrInvalidToken:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	default:
		logger.Error("executing app handler", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
