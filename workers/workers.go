package workers

import (
	"encoding/json"

	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/usecases"
	"github.com/nerdyworm/rsq"
)

func CreateJsError(job *rsq.Job) error {
	payload := map[string]interface{}{}
	err := json.Unmarshal(job.Payload, &payload)
	if err != nil {
		return err
	}

	id := int64(payload["projectID"].(float64))
	project, err := models.Projects.FindByID(id)
	if err != nil {
		return err
	}

	raw, err := json.Marshal(payload["data"])
	if err != nil {
		logger.Error("marshalling payload", "payload", payload)
		return err
	}

	e := models.NewEvent(project, "js.error", string(raw))
	err = models.Events.Insert(e)
	if err != nil {
		logger.Error("inserting event", "err", err)
		return err
	}

	return usecases.AfterErrorEventCreated(e)
}

func InvitationDeliver(job *rsq.Job) error {
	var token string
	err := json.Unmarshal(job.Payload, &token)
	if err != nil {
		return err
	}

	return usecases.InvitationDeliver(token)
}

func CheckEmberDependencies(job *rsq.Job) error {
	projectID := 0
	err := json.Unmarshal(job.Payload, &projectID)
	if err != nil {
		return err
	}

	return usecases.CheckEmberDependencies(int64(projectID), nil)
}

// TODO
// - clean up data that is expired
func Cron(job *rsq.Job) error {
	logger.Info("cron job running", "t", string(job.Payload))

	if err := enqueueDependencyChecking(); err != nil {
		logger.Error("could not enqueueDependencyChecking", "err", err)
		return err
	}

	return nil
}

func enqueueDependencyChecking() err {
	repositories, err := models.AllRepositories()
	if err != nil {
		return err
	}

	for _, r := range repositories {
		if r.GithubOK() {
			jobs.Push("CheckEmberDependencies", r.ProjectID)
		}
	}

	return nil
}
