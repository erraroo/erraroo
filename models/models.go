package models

import (
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
	Errors = &errorsStore{store}
	Groups = &groupsStore{store}
	Plans = &plansStore{store}
	Prefs = &prefsStore{store}
	Projects = &projectsStore{store}
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
