package models

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/erraroo/erraroo/emailer"
	_ "github.com/lib/pq"
)

func Setup(config string) error {
	var err error

	store, err = NewStore(config)
	if err != nil {
		return err
	}

	Accounts = &accountsStore{store}
	Events = &eventsStore{store, s3.New(nil)}
	Errors = &errorsStore{store}
	Invitations = &invitationsStore{store}
	PasswordRecovers = &passwordRecoversStore{store}
	Plans = &plansStore{store}
	Prefs = &prefsStore{store}
	Projects = &projectsStore{store}
	RateLimitNotifcations = &rateLimitNotifcationsStore{store}
	Timings = &timingsStore{store}
	Users = &usersStore{store}

	return nil
}

func Shutdown() {
	store.Close()
}

// SetupForTesting setups up the database and such for testing
func SetupForTesting() {
	encrypter = &dummyPasswordEncrypter{}
	emailer.Use(&emailer.DummySender{})
	Migrate()
}

func Migrate() {
	store.Migrate()
}
