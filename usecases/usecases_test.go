package usecases

import (
	"os"
	"testing"

	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/emailer"
	"github.com/erraroo/erraroo/models"
)

var (
	emailSender = &mockSender{}
)

func TestMain(m *testing.M) {
	config.Env = "test"
	config.Postgres = "dbname=erraroo_test sslmode=disable"

	models.Setup(config.Postgres)
	models.SetupForTesting()
	defer models.Shutdown()

	emailer.Use(emailSender)
	emailSender.Clear()

	ret := m.Run()
	os.Exit(ret)
}
