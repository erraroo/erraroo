package serializers

import (
	"fmt"
	"math"

	"github.com/erraroo/erraroo/models"
)

type Error struct {
	*models.Error
	Links ErrorLinks `json:"links"`
}

type ErrorLinks struct {
	Tags string `json:"tags"`
}

type ShowError struct {
	Error Error
}

type UpdateError struct {
	Error    Error
	Projects []Project
}

func NewError(e *models.Error) Error {
	links := ErrorLinks{
		Tags: fmt.Sprintf("/api/v1/errors/%d/tags", e.ID),
	}

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

type Tag struct {
	models.TagValue
}

type Errors struct {
	Errors []Error
	//Tags   []Tag
	Meta struct {
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

	//mapping := make(map[int64][]int64)
	//for i, t := range results.Tags {
	//groups.Tags[i] = Tag{t}

	//if ids, ok := mapping[t.ErrorID]; ok {
	//mapping[t.ErrorID] = append(ids, t.ID)
	//} else {
	//mapping[t.ErrorID] = []int64{t.ID}
	//}
	//}

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
