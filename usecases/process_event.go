package usecases

import (
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"
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

	return nil
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
