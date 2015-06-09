package serializers

import (
	"math"

	"github.com/erraroo/erraroo/models"
)

type Error struct {
	*models.Error
}

type ShowError struct {
	Error Error
}

func NewShowError(p *models.Error) ShowError {
	return ShowError{
		Error: Error{p},
	}
}

type Errors struct {
	Errors []Error
	Meta   struct {
		Pagination Pagination
	}
}

func NewErrors(results models.ErrorResults) Errors {
	e := Errors{}
	e.Errors = make([]Error, len(results.Errors))

	for i, p := range results.Errors {
		e.Errors[i] = Error{p}
	}

	pages := math.Ceil(float64(results.Total) / float64(results.Query.PerPageOrDefault()))
	e.Meta.Pagination.Pages = int(pages)
	e.Meta.Pagination.Page = results.Query.PageOrDefault()
	e.Meta.Pagination.Limit = results.Query.PerPageOrDefault()
	e.Meta.Pagination.Total = int(results.Total)

	return e
}
