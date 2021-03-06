package models

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/redis.v3"

	"go/build"

	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/tanel/dbmigrate"
)

var (
	store                 *Store
	Events                EventsStore
	Errors                ErrorsStore
	Invitations           InvitationsStore
	PasswordRecovers      PasswordRecoversStore
	Plans                 PlansStore
	Prefs                 PrefsStore
	Projects              ProjectsStore
	RateLimitNotifcations RateLimitNotifcationsStore
	Timings               TimingsStore
	Users                 UsersStore
)

// Store is the abstraction used to interact with the
// database.
type Store struct {
	gorm.DB
	redis *redis.Client
}

// NewStore initializes a new Store
func NewStore() (*Store, error) {
	dbGorm, err := gorm.Open("postgres", config.Postgres)
	if err != nil {
		logger.Error("could not connect to postgres", "err", err)
		return nil, err
	}

	dbGorm.DB().SetMaxOpenConns(10)
	dbGorm.LogMode(config.LogSql)

	client := redis.NewClient(&redis.Options{
		Addr:        config.Redis,
		PoolSize:    10,
		PoolTimeout: 5 * time.Second,
	})

	return &Store{dbGorm, client}, nil
}

// Close closes all the connections
func (s *Store) Close() {
	s.DB.Close()
	s.redis.Close()
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

	err := dbmigrate.Run(s.DB.DB(), path)
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
	return s.DB.DB().Query(query, args...)
}

func (s *Store) QueryRow(query string, args ...interface{}) *sql.Row {
	start := time.Now()
	defer s.logQuery(start, query, args...)
	return s.DB.DB().QueryRow(query, args...)
}

func (s *Store) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	defer s.logQuery(start, query, args...)
	return s.DB.DB().Exec(query, args...)
}
