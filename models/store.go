package models

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"go/build"

	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/logger"
	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	_ "github.com/lib/pq"
	"github.com/tanel/dbmigrate"
)

var (
	builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	store    *Store
	Accounts AccountsStore
	Errors   ErrorsStore
	Groups   GroupsStore
	Prefs    PrefsStore
	Projects ProjectsStore
	Timings  TimingsStore
	Users    UsersStore
)

// Store is the abstraction used to interact with the
// database.
type Store struct {
	*sqlx.DB
}

// NewStore initializes a new Store
func NewStore(config string) (*Store, error) {
	db, err := sqlx.Connect("postgres", config)
	if err != nil {
		logger.Error("could not connect to postgres", "err", err)
		return nil, err
	}

	return &Store{db}, nil
}

// Close closes all the connections
func (s *Store) Close() {
	s.DB.Close()
}

// Migrate the database to the latest version
func (s *Store) Migrate() {
	path := config.MigrationsPath
	if path == "" {
		pkg, err := build.Default.Import("github.com/erraroo/erraroo", "", 0x0)
		if err != nil {
			panic(err)
		}

		path = filepath.Join(pkg.Dir, "db", "migrations")
	}

	err := dbmigrate.Run(s.DB.DB, path)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func (s *Store) logQuery(start time.Time, query string, args ...interface{}) {
	end := time.Since(start)
	query = strings.Replace(query, "\n", "", -1)
	query = strings.Replace(query, "\t", "", -1)

	if int(end.Nanoseconds()/1000000) > 100 {
		logger.Info("slow query", "query", query, "args", fmt.Sprintf("%v", args), "runtime", time.Since(start))
	} else {
		logger.Debug("store", "query", query, "args", fmt.Sprintf("%v", args), "runtime", time.Since(start))
	}
}

func (s *Store) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	defer s.logQuery(start, query, args...)
	return s.DB.Query(query, args...)
}

func (s *Store) QueryRow(query string, args ...interface{}) *sql.Row {
	start := time.Now()
	defer s.logQuery(start, query, args...)
	return s.DB.QueryRow(query, args...)
}

func (s *Store) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	defer s.logQuery(start, query, args...)
	return s.DB.Exec(query, args...)
}

func (s *Store) Select(dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	defer s.logQuery(start, query, args...)
	return s.DB.Select(dest, query, args...)
}

func (s *Store) Get(dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	defer s.logQuery(start, query, args...)
	return s.DB.Get(dest, query, args...)
}
