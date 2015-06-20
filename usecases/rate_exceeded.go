package usecases

import (
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
)

func RateExceeded(token string) error {
	project, err := models.Projects.FindByToken(token)
	if err != nil {
		return err
	}

	logger.Error("rate limited exceeded", "account", project.AccountID)
	return nil
}
