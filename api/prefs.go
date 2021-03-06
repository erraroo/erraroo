package api

import (
	"net/http"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

type PrefsUpdateRequest struct {
	Pref PrefsParams
}

type PrefsParams struct {
	EmailOnError bool
}

func PrefsUpdate(w http.ResponseWriter, r *http.Request, ctx *cx.Context) error {
	pref, err := models.Prefs.Get(ctx.User)
	if err != nil {
		return err
	}

	request := PrefsUpdateRequest{}
	err = Decode(r, &request)
	if err != nil {
		return err
	}

	pref.EmailOnError = request.Pref.EmailOnError
	err = models.Prefs.Update(pref)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, serializers.NewShowPref(pref))
}
