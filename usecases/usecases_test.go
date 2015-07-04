package usecases

import (
	"os"
	"testing"

	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/emailer"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/models"
	"github.com/nerdyworm/rsq"
)

var (
	emailSender = &mockSender{}
	queue       = rsq.NewMemoryAdapter()
)

func TestMain(m *testing.M) {
	config.Env = "test"
	config.Postgres = "dbname=erraroo_test sslmode=disable"

	models.Setup(config.Postgres)
	models.SetupForTesting()
	defer models.Shutdown()

	emailer.Use(emailSender)
	emailSender.Clear()

	jobs.Use(queue)

	ret := m.Run()
	os.Exit(ret)
}
