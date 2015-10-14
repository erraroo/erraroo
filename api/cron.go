package api

import (
	"net/http"
	"time"

	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/logger"
	"github.com/gorilla/mux"
)

func Cron(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	logger.Info("cron", "token", token)

	if token != config.CronToken {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	jobs.Push("cron", time.Now().UTC().Unix())
}
