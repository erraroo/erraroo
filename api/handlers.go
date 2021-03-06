package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/erraroo/erraroo/api/bus"
	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/gorilla/mux"
	"github.com/nerdyworm/puller"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	http.Error(w, "not found", http.StatusNotFound)
	return nil
}

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "ok", http.StatusOK)
}

func MeHandler(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	if ctx.User == nil {
		w.WriteHeader(http.StatusForbidden)
	} else {
		prefs, err := models.Prefs.Get(ctx.User)
		if err != nil {
			logger.Error("getting prefs", "err", err, "user", ctx.User.ID, "email", ctx.User.Email)
			return err
		}

		return JSON(w, http.StatusOK, serializers.NewShowUser(ctx.User, prefs))
	}

	return nil
}

func UsersShow(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	id := mux.Vars(r)["id"]

	user := ctx.User

	if id != "me" {
		userID, err := StrToID(id)
		if err != nil {
			return err
		}

		user, err = models.Users.FindByID(userID)
		if err != nil {
			return err
		}
	}

	if !ctx.User.CanSee(user) {
		return models.ErrNotFound
	}

	prefs, err := models.Prefs.Get(user)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowUser(user, prefs))
}

func Backlog(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	channels := puller.Channels{}
	accountKey := fmt.Sprintf("accounts.%d", ctx.User.AccountID)

	channels["global"] = lastID(r.Form, "global")
	channels[accountKey] = lastID(r.Form, accountKey)

	backlog, err := bus.Pull(channels, 60*time.Second)
	if err != nil {
		logger.Error("pulling backlog", "err", err)
		return err
	}

	return JSON(w, http.StatusOK, backlog)
}

func lastID(values url.Values, key string) int64 {
	if value, ok := values[key]; ok {
		if len(value) > 0 {
			if i, err := strconv.Atoi(value[0]); err == nil {
				return int64(i)
			}

		}
	}

	return 0
}
