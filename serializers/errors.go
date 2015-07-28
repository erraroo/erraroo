package serializers

import (
	"math"

	"github.com/erraroo/erraroo/models"
)

type Error struct {
	*models.Error
	Links ErrorLinks `json:"links"`
}

type ErrorLinks struct{}

type ShowError struct {
	Error Error
}

type UpdateError struct {
	Error    Error
	Projects []Project
}

func NewError(e *models.Error) Error {
	links := ErrorLinks{}

	return Error{e, links}
}

func NewShowError(e *models.Error) ShowError {
	return ShowError{
		Error: NewError(e),
	}
}

func NewUpdateError(p *models.Project, e *models.Error) UpdateError {
	return UpdateError{
		Error: NewError(e),
		Projects: []Project{
			NewProject(p),
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

	for i, e := range results.Errors {
		groups.Errors[i] = NewError(e)
	}

	pages := math.Ceil(float64(results.Total) / float64(results.Query.PerPageOrDefault()))
	groups.Meta.Pagination.Pages = int(pages)
	groups.Meta.Pagination.Page = results.Query.PageOrDefault()
	groups.Meta.Pagination.Limit = results.Query.PerPageOrDefault()
	groups.Meta.Pagination.Total = int(results.Total)

	return groups
}
