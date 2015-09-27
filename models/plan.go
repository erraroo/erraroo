package models

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/erraroo/erraroo/logger"
)

var plans = make(map[string]Plan)

func byName(name string) Plan {
	if plan, ok := plans[name]; ok {
		return plan
	} else {
		return plans["default"]
	}
}

func init() {
	plans["default"] = Plan{
		Name:                     "default",
		DataRetentionInDays:      7,
		ProjectsLimit:            3,
		PriceInCents:             2900,
		RateLimit:                20,
		RateLimitDurationSeconds: 60,
	}
	plans["small"] = Plan{
		Name:                     "small",
		DataRetentionInDays:      7,
		ProjectsLimit:            3,
		PriceInCents:             2900,
		RateLimit:                20,
		RateLimitDurationSeconds: 60,
	}
	plans["pro"] = Plan{
		Name:                     "pro",
		DataRetentionInDays:      14,
		ProjectsLimit:            10,
		PriceInCents:             8900,
		RateLimit:                40,
		RateLimitDurationSeconds: 60,
	}
	plans["enterprise"] = Plan{
		Name:                     "enterprise",
		DataRetentionInDays:      30,
		ProjectsLimit:            50,
		PriceInCents:             19900,
		RateLimit:                120,
		RateLimitDurationSeconds: 60,
	}
}

type Plan struct {
	AccountID                int64
	DataRetentionInDays      int
	Name                     string
	PriceInCents             int
	RateLimitDurationSeconds int
	RateLimit                int
	ProjectsLimit            int
}

func (p Plan) RateLimitDuration() time.Duration {
	return time.Duration(p.RateLimitDurationSeconds) * time.Second
}

type PlansStore interface {
	FindByToken(string) (*Plan, error)
	Get(*Account) (*Plan, error)
	Update(*Plan) error

	Create(*Account, string) (*Plan, error)
}

type plansStore struct {
	*Store
}

var plansColumns = []string{
	"account_id",
	"data_retention_in_days",
	"name",
	"price_in_cents",
	"rate_limit_duration_seconds",
	"rate_limit",
	"projects_limit",
}

func (s *plansStore) FindByToken(token string) (*Plan, error) {
	p := &Plan{}

	query := `select
			plans.account_id,
			plans.data_retention_in_days,
			plans.name,
			plans.price_in_cents,
			plans.rate_limit_duration_seconds,
			plans.rate_limit,
			plans.projects_limit
		from plans join projects on projects.token = $1 and projects.account_id = plans.account_id limit 1;`

	err := s.QueryRow(query, token).Scan(
		&p.AccountID,
		&p.DataRetentionInDays,
		&p.Name,
		&p.PriceInCents,
		&p.RateLimitDurationSeconds,
		&p.RateLimit,
		&p.ProjectsLimit,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		logger.Error("finding plan by token", "token", token, "err", err)
		return nil, err
	}

	return p, nil
}

func (s *plansStore) Create(account *Account, name string) (*Plan, error) {
	model := byName(name)
	model.AccountID = account.ID

	query := "insert into plans("
	query += strings.Join(plansColumns, ",")
	query += ") values ($1,$2,$3,$4,$5,$6,$7);"

	_, err := s.Exec(query,
		model.AccountID,
		model.DataRetentionInDays,
		model.Name,
		model.PriceInCents,
		model.RateLimitDurationSeconds,
		model.RateLimit,
		model.ProjectsLimit,
	)

	if err != nil {
		logger.Error("inserting plan", "err", err, "account", account.ID)
	}

	return &model, err
}

func (s *plansStore) Get(account *Account) (*Plan, error) {
	p := new(Plan)

	query := `select
			plans.account_id,
			plans.data_retention_in_days,
			plans.name,
			plans.price_in_cents,
			plans.rate_limit_duration_seconds,
			plans.rate_limit,
			plans.projects_limit
		from plans where account_id = $1 limit 1;`
	err := s.QueryRow(query, account.ID).Scan(
		&p.AccountID,
		&p.DataRetentionInDays,
		&p.Name,
		&p.PriceInCents,
		&p.RateLimitDurationSeconds,
		&p.RateLimit,
		&p.ProjectsLimit,
	)

	if err == sql.ErrNoRows {
		defaultPlan := byName("default")
		logger.Error("using default plan", "account", account.ID)
		return &defaultPlan, nil
	}

	if err != nil {
		logger.Error("getting plans", "err", err, "account", account.ID)
	}

	return p, err
}

func (s *plansStore) Update(p *Plan) error {
	query := `update plans set
		data_retention_in_days=$1,
		name=$2,
		price_in_cents=$3,
		rate_limit_duration_seconds=$4,
		rate_limit=$5,
		projects_limit=$6
	where account_id=$7`

	_, err := s.Exec(query,
		p.DataRetentionInDays,
		p.Name,
		p.PriceInCents,
		p.RateLimitDurationSeconds,
		p.RateLimit,
		p.ProjectsLimit,
		p.AccountID,
	)
	return err
}

type planCache struct {
	plans   map[string]*Plan
	expires map[string]time.Time
	lock    *sync.RWMutex
	ttl     time.Duration
}

func NewPlanCache() *planCache {
	return &planCache{
		plans:   map[string]*Plan{},
		expires: map[string]time.Time{},
		lock:    &sync.RWMutex{},
		ttl:     1 * time.Minute,
	}
}

func (cache *planCache) FindByToken(token string) (*Plan, error) {
	cache.lock.RLock()
	plan := cache.plans[token]
	expires, cached := cache.expires[token]
	cache.lock.RUnlock()

	if cached && expires.After(time.Now()) {
		return plan, nil
	}

	return cache.put(token)
}

func (cache *planCache) put(token string) (*Plan, error) {
	plan, err := Plans.FindByToken(token)
	if err != nil {
		return nil, err
	}

	cache.lock.Lock()
	cache.plans[token] = plan
	cache.expires[token] = time.Now().Add(cache.ttl)
	cache.lock.Unlock()

	logger.Info("cached plan", "token", token, "plan", plan)
	return plan, nil
}
