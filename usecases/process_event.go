package usecases

import (
	"fmt"

	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
	"github.com/nerdyworm/puller"
)

func ProcessEvent(eventID int64) error {
	e, err := models.Events.FindByID(eventID)
	if err != nil {
		return err
	}

	err = e.PostProcess()
	if err != nil {
		logger.Error("event.PostProcess", "err", err, "event.ID", e.ID)
		return err
	}

	if e.Kind == "js.error" {
		return afterJsErrorProcessed(e)
	}

	return nil
}

type Event struct {
	Name    string
	Payload interface{}
}

func afterJsErrorProcessed(e *models.Event) error {
	p, err := models.Projects.FindByID(e.ProjectID)
	if err != nil {
		return err
	}

	group, err := models.Errors.FindOrCreate(p, e)
	if err != nil {
		logger.Error("finding or creating error group", "err", err)
		return err
	}

	if group.ShouldNotify() {
		err = notifyUsersOfNewError(p, group)
		if err != nil {
			logger.Error("notifying users", "err", err, "project", p.ID, "group", group.ID)
			return err
		}
	}

	err = models.Errors.Touch(group)
	if err != nil {
		logger.Error("touching group", "err", err, "group", group.ID)
		return err
	}

	if !group.Muted {
		project, err := models.Projects.FindByID(group.ProjectID)
		if err != nil {
			logger.Error("finding project", "err", err, "group", group.ID, "project", group.ProjectID)
			return err
		}

		puller.Publish(accountChannel(project.AccountID), Event{
			Name:    "errors.update",
			Payload: serializers.NewUpdateError(project, group),
		})
	}

	return nil
}

func accountChannel(accountID int64) string {
	return fmt.Sprintf("accounts.%d", accountID)
}

func notifyUsersOfNewError(project *models.Project, group *models.Error) error {
	users, err := models.Users.ByAccountID(project.AccountID)
	if err != nil {
		return err
	}

	for _, user := range users {
		pref, err := models.Prefs.Get(user)
		if err != nil {
			logger.Error("getting prefs for user", "err", err, "user", user.ID, "email", user.Email)
			continue
		}

		if pref.EmailOnError {
			err = mailers.DeliverNewErrorNotification(user, group)
			if err != nil {
				logger.Error("deliver new group notifcation", "err", err, "user", user.ID, "email", user.Email)
				continue
			}
		}

	}

	return nil
}
