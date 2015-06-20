package usecases

import "github.com/erraroo/erraroo/logger"

func RateExceeded(token string) error {
	logger.Info("SHOULD SEND EMAIL TO PROJECT FOLKS", "id", "token", token)
	return nil
}
