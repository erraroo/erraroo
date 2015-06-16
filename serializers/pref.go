package serializers

import "github.com/erraroo/erraroo/models"

type Pref struct {
	*models.Pref
	ID int64
}

type ShowPref struct {
	Pref Pref
}

func NewPref(p *models.Pref) Pref {
	return Pref{p, p.UserID}
}

func NewShowPref(p *models.Pref) ShowPref {
	return ShowPref{NewPref(p)}
}
