package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/dgrijalva/jwt-go"
	"github.com/erraroo/erraroo/api/errors"
	"github.com/erraroo/erraroo/api/events"
	"github.com/erraroo/erraroo/api/groups"
	"github.com/erraroo/erraroo/api/projects"
	"github.com/erraroo/erraroo/api/sessions"
	"github.com/erraroo/erraroo/api/signups"
	"github.com/erraroo/erraroo/api/timings"
	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/nerdyworm/rsq"
	"github.com/rs/cors"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

}

// Context is the context that is passed into each request
type Context struct {
	Store *models.Store
	User  *models.User
	Queue rsq.Queue
}

// App is the main application for erraroo
type App struct {
	Store      *models.Store
	Router     *mux.Router
	JobRouter  *rsq.JobRouter
	HTTPServer http.Handler
	Queue      rsq.Queue
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.HTTPServer.ServeHTTP(w, r)
}

func (a *App) Run(job *rsq.Job) error {
	return a.JobRouter.Run(job)
}

func (a *App) newContext(r *http.Request) (*cx.Context, error) {
	var err error
	ctx := &cx.Context{
		Queue: a.Queue,
	}

	if r != nil {
		id, err := getCurrentUserID(r)
		if err != nil {
			log.Printf("[error] getting current user: `%v`\n", err)
			return ctx, err
		}

		if id == 0 {
			return ctx, err
		}

		ctx.User, err = models.Users.FindByID(id)
		if err != nil {
			log.Printf("[error] finding user from token: `%v`\n", err)
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
		log.Printf("[error] parsing token: `%v`\n", err)
		return 0, err
	}

	id := token.Claims["user_id"].(float64)
	return int64(id), nil
}

// SetupForTesting sets the app up to be tested
func (a *App) SetupForTesting() {
	models.SetupForTesting()
	models.Migrate()
	a.Queue = rsq.NewMemoryAdapter()
}

// Shutdown shutsdown all the things
func (a *App) Shutdown() {
	defer a.Queue.Shutdown()
	models.Shutdown()
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
			cx.JSON(w, http.StatusBadRequest, err)
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
		ctx, err := a.newContext(nil)
		if err != nil {
			return err
		}

		return fn(job, ctx)
	}
}

func handleError(err error, w http.ResponseWriter) {
	switch err {
	case models.ErrNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case cx.ErrLoginRequired:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	default:
		log.Printf("Error executing app handler `%v`\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func New() *App {
	a := &App{}
	a.Router = mux.NewRouter()
	a.Router.NotFoundHandler = a.Handler(notFoundHandler)
	a.Router.Handle("/api/v1/signups", a.Handler(signups.Create)).Methods("POST")
	a.Router.Handle("/api/v1/sessions", a.Handler(sessions.Create)).Methods("POST")
	a.Router.Handle("/api/v1/sessions", a.AuthroziedHandler(sessions.Destroy)).Methods("DELETE")
	a.Router.Handle("/api/v1/me", a.AuthroziedHandler(showMe)).Methods("GET")
	a.Router.Handle("/api/v1/events", a.Handler(events.Create)).Methods("POST")
	a.Router.Handle("/api/v1/projects", a.AuthroziedHandler(projects.Index)).Methods("GET")
	a.Router.Handle("/api/v1/projects", a.AuthroziedHandler(projects.Create)).Methods("POST")
	a.Router.Handle("/api/v1/projects/{id}", a.AuthroziedHandler(projects.Show)).Methods("GET")
	a.Router.Handle("/api/v1/projects/{id}", a.AuthroziedHandler(projects.Update)).Methods("PUT")
	a.Router.Handle("/api/v1/errors/{id}", a.AuthroziedHandler(errors.Show)).Methods("GET")
	a.Router.Handle("/api/v1/errors", a.AuthroziedHandler(errors.Index)).Methods("GET")
	a.Router.Handle("/api/v1/groups", a.AuthroziedHandler(groups.Index)).Methods("GET")
	a.Router.Handle("/api/v1/groups/{id}", a.AuthroziedHandler(groups.Show)).Methods("GET")
	a.Router.Handle("/api/v1/groups/{id}", a.AuthroziedHandler(groups.Update)).Methods("PUT")
	a.Router.Handle("/api/v1/users/{id}", a.AuthroziedHandler(showMe)).Methods("GET")
	a.Router.Handle("/api/v1/timings", a.AuthroziedHandler(timings.Index)).Methods("GET")

	c := cors.New(cors.Options{
		//Debug:          true,
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders: []string{"Accept", "Authorization"},
	})

	a.HTTPServer = alice.New(
		loggingHandler,
		c.Handler,
	).Then(a.Router)

	a.Queue = rsq.NewSqsAdapter(rsq.SqsOptions{
		AwsConfig: &aws.Config{
			Region: config.AwsRegion,
			Credentials: credentials.NewStaticCredentials(
				config.AwsAccessKeyID,
				config.AwsSecretAccessKey,
				""),
		},
		LongPollTimeout:   10,
		MessagesPerWorker: 1,
		QueueURL:          config.SqsQueueURL,
	})
	a.JobRouter = rsq.NewJobRouter()
	a.JobRouter.Handle("create.error", a.JobHandler(jobs.AfterCreateError))

	return a
}

func notFoundHandler(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	http.Error(w, "not found", http.StatusNotFound)
	return nil
}

func showMe(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	if ctx.User == nil {
		w.WriteHeader(http.StatusForbidden)
	} else {
		return cx.JSON(w, http.StatusOK, serializers.NewShowUser(ctx.User))
	}

	return nil
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("[http] method=%s url=%s runtime=%v", r.Method, r.URL, time.Since(start))
	})
}
