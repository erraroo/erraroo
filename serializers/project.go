package serializers

import (
	"fmt"

	"github.com/erraroo/erraroo/models"
)

type Project struct {
	*models.Project
	Links ProjectLinks `json:"links"`
}

type ShowProject struct {
	Project Project
}

func NewShowProject(p *models.Project) ShowProject {
	return ShowProject{
		Project: NewProject(p),
	}
}

type Projects struct {
	Projects []Project
}

func NewProject(p *models.Project) Project {
	links := ProjectLinks{
		Repository: fmt.Sprintf("/api/v1/projects/%d/repository", p.ID),
		Revisions:  fmt.Sprintf("/api/v1/projects/%d/revisions", p.ID),
	}

	return Project{p, links}
}

type ProjectLinks struct {
	Repository string `json:"repository"`
	Revisions  string `json:"revisions"`
}

func NewProjects(ps []*models.Project) Projects {
	projects := Projects{}
	projects.Projects = make([]Project, len(ps))

	for i, p := range ps {
		projects.Projects[i] = NewProject(p)
	}

	return projects
}
