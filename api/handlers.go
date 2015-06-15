package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/serializers"
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
		return JSON(w, http.StatusOK, serializers.NewShowUser(ctx.User))
	}

	return nil
}
