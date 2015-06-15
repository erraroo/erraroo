package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/erraroo/erraroo/api"
	"github.com/erraroo/erraroo/app"
	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
	"github.com/nerdyworm/rsq"
)

func main() {
	err := models.Setup(config.Postgres)
	if err != nil {
		logger.Fatal("could not connect to database", "err", err)
	}

	jobs.Use(rsq.NewSqsAdapter(rsq.SqsOptions{
		LongPollTimeout:   config.SqsLongPollTimeout,
		MessagesPerWorker: config.SqsMessagesPerWorker,
		QueueURL:          config.SqsQueueURL,
	}))

	erraroo := app.New()
	defer erraroo.Shutdown()

	a := cli.NewApp()
	a.Name = "erraroo"
	a.Author = "Benjamin Silas Rhodes"
	a.Email = "ben@nerdyworm.com"
	a.Version = "0.0.3"
	a.Commands = []cli.Command{
		cli.Command{
			Name:        "server",
			Description: "start the http server",
			Action: func(c *cli.Context) {
				startServer(erraroo)
			},
		},

		cli.Command{
			Name:        "workers",
			Description: "start workers",
			Action: func(c *cli.Context) {
				startWorkers(erraroo)
				ch := make(chan os.Signal)
				signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
				<-ch
			},
		},

		cli.Command{
			Name:        "migrate",
			Description: "run the migrations on the database",
			Action: func(c *cli.Context) {
				config.MigrationsPath = c.Args().First()
				models.Migrate()
			},
		},

		cli.Command{
			Name:        "development",
			Description: "start both http server and workers",
			Action: func(c *cli.Context) {
				startWorkers(erraroo)
				startServer(erraroo)
			},
		},

		cli.Command{
			Name:        "process",
			Description: "process an error at the cli",
			Action: func(c *cli.Context) {
				id, err := api.StrToID(c.Args().First())
				if err != nil {
					logger.Fatal("could not parse id", "err", err)
				}

				err = jobs.AfterCreateErrorFn(id)
				if err != nil {
					logger.Fatal("could not complete job", "err", err)
				}
			},
		},
	}

	a.Run(os.Args)
}

func startServer(a *app.App) {
	logger.Info("server listening", "port", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), a)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func startWorkers(a *app.App) {
	for i := 0; i < config.QueueWorkers; i++ {
		logger.Info("starting worker", "number", i)
		go jobs.Work(a)
	}
}
