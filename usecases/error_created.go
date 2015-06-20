package usecases

import (
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"
)

func ErrorCreated(eventID int64) error {
	e, err := models.Events.FindByID(eventID)
	if err != nil {
		return err
	}

	err = processEvent(e)
	if err != nil {
		return err
	}

	return afterEventProcessed(e)
}

func processEvent(e *models.Event) error {
	resources := models.NewResourceStore()

	err := e.PopulateStackContext(resources)
	if err != nil {
		logger.Error("populating stack context", "err", err, "error.ID", e.ID)
	} else {
		err = models.Events.Update(e)
		if err != nil {
			logger.Error("updating error", "err", err, "error", e.ID)
			return err
		}
	}

	return nil
}

func afterEventProcessed(e *models.Event) error {
	p, err := models.Projects.FindByID(e.ProjectID)
	if err != nil {
		return err
	}

	group, err := models.Groups.FindOrCreate(p, e)
	if err != nil {
		logger.Error("finding or creating group", "err", err)
		return err
	}

	if group.ShouldNotify() {
		err = notifyUsersOfNewError(p, group)
		if err != nil {
			logger.Error("group notifcations", "err", err, "project", p.ID, "group", group.ID)
			return err
		}
	}

	err = models.Groups.Touch(group)
	if err != nil {
		logger.Error("touching group", "err", err, "group", group.ID)
		return err
	}

	return nil

}

func notifyUsersOfNewError(project *models.Project, group *models.Group) error {
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
			err = mailers.DeliverNewGroupNotification(user, group)
			if err != nil {
				logger.Error("deliver new group notifcation", "err", err, "user", user.ID, "email", user.Email)
				continue
			}
		}

	}

	return nil
}
