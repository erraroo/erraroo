package models

import (
	"github.com/erraroo/erraroo/emailer"
	_ "github.com/lib/pq"
)

func Setup(config string) error {
	var err error

	MainStore, err = NewStore(config)
	if err != nil {
		return err
	}

	Accounts = &accountsStore{MainStore}
	Errors = &errorsStore{MainStore}
	Groups = &groupsStore{MainStore}
	Projects = &projectsStore{MainStore}
	Timings = &timingsStore{MainStore}
	Users = &usersStore{MainStore}

	return nil
}

func Shutdown() {
	MainStore.Close()
}

// SetupForTesting setups up the database and such for testing
func SetupForTesting() {
	encrypter = &dummyPasswordEncrypter{}
	emailer.Use(&emailer.DummySender{})
	Migrate()
}

func Migrate() {
	MainStore.Migrate()
}
