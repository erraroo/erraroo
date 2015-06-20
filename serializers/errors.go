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

type UpdateError struct {
	Error    Error
	Projects []Project
}

func NewShowError(g *models.Error) ShowError {
	return ShowError{
		Error: Error{g},
	}
}

func NewUpdateError(p *models.Project, g *models.Error) UpdateError {
	return UpdateError{
		Error: Error{g},
		Projects: []Project{
			Project{p},
		},
	}
}

type Errors struct {
	Errors []Error
	Meta   struct {
		Pagination Pagination
	}
}

type Pagination struct {
	Pages int
	Page  int
	Limit int
	Total int
}

func NewErrors(results models.ErrorResults) Errors {
	groups := Errors{}
	groups.Errors = make([]Error, len(results.Errors))

	for i, p := range results.Errors {
		groups.Errors[i] = Error{p}
	}

	pages := math.Ceil(float64(results.Total) / float64(results.Query.PerPageOrDefault()))
	groups.Meta.Pagination.Pages = int(pages)
	groups.Meta.Pagination.Page = results.Query.PageOrDefault()
	groups.Meta.Pagination.Limit = results.Query.PerPageOrDefault()
	groups.Meta.Pagination.Total = int(results.Total)

	return groups
}
