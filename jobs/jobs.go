package jobs

import (
	"encoding/json"
	"log"

	"github.com/erraroo/erraroo/cx"
	"github.com/erraroo/erraroo/models"
	"github.com/nerdyworm/rsq"
)

var Queue rsq.Queue

func Setup(q rsq.Queue) {
	Queue = q
}

func AfterCreateError(job *rsq.Job, ctx *cx.Context) error {
	var id int64
	err := json.Unmarshal(job.Payload, &id)
	if err != nil {
		return err
	}

	return AfterCreateErrorFn(id)
}

func AfterCreateErrorFn(id int64) error {
	e, err := models.Errors.FindByID(id)
	if err != nil {
		return err
	}

	project, err := models.Projects.FindByID(e.ProjectID)
	if err != nil {
		return err
	}

	resources := models.NewResourceStore()
	err = e.PopulateStackContext(resources)
	if err != nil {
		log.Println(err)
		return err
	}

	err = models.Errors.Update(e)
	if err != nil {
		log.Println(err)
		return err
	}

	group, err := models.Groups.FindOrCreate(project, e)
	if err != nil {
		log.Printf("error creating group: %v\n", err)
		return err
	}

	return models.Groups.Touch(group)
}