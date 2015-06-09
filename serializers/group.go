package serializers

import (
	"math"

	"github.com/erraroo/erraroo/models"
)

type Group struct {
	*models.Group
}

type ShowGroup struct {
	Group Group
}

type UpdateGroup struct {
	Group    Group
	Projects []Project
}

func NewShowGroup(g *models.Group) ShowGroup {
	return ShowGroup{
		Group: Group{g},
	}
}

func NewUpdateGroup(p *models.Project, g *models.Group) UpdateGroup {
	return UpdateGroup{
		Group: Group{g},
		Projects: []Project{
			Project{p},
		},
	}
}

type Groups struct {
	Groups []Group
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

func NewGroups(results models.GroupResults) Groups {
	groups := Groups{}
	groups.Groups = make([]Group, len(results.Groups))

	for i, p := range results.Groups {
		groups.Groups[i] = Group{p}
	}

	pages := math.Ceil(float64(results.Total) / float64(results.Query.PerPageOrDefault()))
	groups.Meta.Pagination.Pages = int(pages)
	groups.Meta.Pagination.Page = results.Query.PageOrDefault()
	groups.Meta.Pagination.Limit = results.Query.PerPageOrDefault()
	groups.Meta.Pagination.Total = int(results.Total)

	return groups
}
