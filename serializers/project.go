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
		Libraries: fmt.Sprintf("/api/v1/projects/%d/libraries", p.ID),
	}

	return Project{p, links}
}

type ProjectLinks struct {
	Libraries string `json:"libraries"`
}

func NewProjects(ps []*models.Project) Projects {
	projects := Projects{}
	projects.Projects = make([]Project, len(ps))

	for i, p := range ps {
		projects.Projects[i] = NewProject(p)
	}

	return projects
}
