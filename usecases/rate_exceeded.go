package usecases

import (
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/mailers"
	"github.com/erraroo/erraroo/models"
)

func RateExceeded(token string) error {
	project, err := models.Projects.FindByToken(token)
	if err != nil {
		return err
	}

	account := &models.Account{ID: project.AccountID}

	recently, err := models.RateLimitNotifcations.WasRecentlyNotified(account)
	if err != nil {
		logger.Error("checking if account was recently notified", "account", account.ID, "token", token, "err", err)
		return err
	}

	if recently {
		logger.Info("account was recently notified, not delivering notifcations", "account", account, "token", token)
		return nil
	}

	err = models.RateLimitNotifcations.Insert(account)
	if err != nil {
		logger.Error("inserting rate limit notifcation", "account", account.ID, "token", token, "err", err)
		return err
	}

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

		// TODO: pref.EmailOnRateLimit
		if pref.EmailOnError {
			err = mailers.DeliverRateLimitNotifcation(user, project)
			if err != nil {
				logger.Error("delivering rate limite notifcation", "err", err, "user", user.ID, "email", user.Email)
				continue
			}
		}
	}

	return nil
}
