package serializers

import (
	"github.com/erraroo/erraroo/models"
)

type Library struct {
	*models.Library
}

type Libraries struct {
	Libraries []Library
}

func NewLibraries(libraries []*models.Library) Libraries {
	s := Libraries{}
	s.Libraries = make([]Library, len(libraries))

	for i, l := range libraries {
		s.Libraries[i] = Library{l}
	}

	return s
}
