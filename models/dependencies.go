package models

import "fmt"

type OutdatedRevision struct {
	ID           int64
	ProjectID    int64
	SHA          string
	Dependencies []Dependency `json:"dependencies"`
}

func (o *OutdatedRevision) Empty() bool {
	return len(o.Dependencies) == 0
}

type Dependency struct {
	Latest   string `json:"latest"`
	Location string `json:"location"`
	Name     string `json:"name"`
	Target   string `json:"target"`
}

func (d Dependency) String() string {
	return fmt.Sprintf("[%s] %s %s -> %s", d.Location, d.Name, d.Target, d.Latest)
}
