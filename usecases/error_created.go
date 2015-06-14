package usecases

import (
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"
)

func ErrorCreated(p *models.Project, e *models.Error) error {
	group, err := models.Groups.FindOrCreate(p, e)
	if err != nil {
		logger.Error("finding or creating group", "err", err)
		return err
	}

	if group.ShouldNotify() {
		err = groupNotifcations(p, group)
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

func groupNotifcations(project *models.Project, group *models.Group) error {
	users, err := models.Users.ByAccountID(project.AccountID)
	if err != nil {
		return err
	}

	for _, user := range users {
		err := mailers.DeliverNewGroupNotification(user, group)
		if err != nil {
			logger.Error("deliver new group notifcation", "err", err, "user", user.ID, "email", user.Email)
			continue
		}
	}

	return nil
}
