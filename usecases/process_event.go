package usecases

import (
	"github.com/erraroo/erraroo/api/bus"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"
	"github.com/erraroo/erraroo/serializers"
)

func ProcessEvent(eventID int64) error {
	event, err := models.Events.FindByID(eventID)
	if err != nil {
		return err
	}

	if event.Kind == "js.error" {
		return AfterErrorEventCreated(event)
	}

	return nil
}

func AfterErrorEventCreated(event *models.Event) error {
	p, err := models.Projects.FindByID(event.ProjectID)
	if err != nil {
		return err
	}

	e, err := models.Errors.FindOrCreate(p, event)
	if err != nil {
		logger.Error("finding or creating error e", "err", err)
		return err
	}

	if e.ShouldNotify() {
		err = notifyUsersOfNewError(p, e)
		if err != nil {
			logger.Error("notifying users", "err", err, "project", p.ID, "e", e.ID)
			return err
		}
	}

	err = models.Errors.Touch(e)
	if err != nil {
		logger.Error("touching e", "err", err, "e", e.ID)
		return err
	}

	if !e.Muted {
		bus.Push(p.Channel(), bus.Notifcation{
			Name:    "errors.update",
			Payload: serializers.NewUpdateError(p, e),
		})
	}

	return nil
}

func notifyUsersOfNewError(project *models.Project, group *models.Error) error {
	users, err := models.Users.FindByAccountID(project.AccountID)
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
				logger.Error("deliver new error notifcation", "err", err, "user", user.ID, "email", user.Email)
				continue
			}
		}
	}

	return nil
}
