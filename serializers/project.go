package serializers

import "github.com/erraroo/erraroo/models"

type Project struct {
	*models.Project
}

type ShowProject struct {
	Project Project
}

func NewShowProject(p *models.Project) ShowProject {
	return ShowProject{
		Project: Project{p},
	}
}

type Projects struct {
	Projects []Project
}

func NewProjects(ps []*models.Project) Projects {
	projects := Projects{}
	projects.Projects = make([]Project, len(ps))

	for i, p := range ps {
		projects.Projects[i] = Project{p}
	}

	return projects
}
