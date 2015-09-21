package models

import (
	"os"
	"testing"

	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/jobs"
	"github.com/nerdyworm/rsq"
)

var (
	queue = rsq.NewMemoryAdapter()
)

func TestMain(m *testing.M) {
	config.Env = "test"
	config.Postgres = "dbname=erraroo_test sslmode=disable"
	Setup()

	SetupForTesting()
	defer Shutdown()

	jobs.Use(queue)

	ret := m.Run()
	os.Exit(ret)
}
