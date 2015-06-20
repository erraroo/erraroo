package models

import (
	"database/sql"
	"time"

	"github.com/erraroo/erraroo/logger"
)

// EventsStore is the interface to error data
type GroupsStore interface {
	FindOrCreate(*Project, *Event) (*Group, error)
	FindQuery(GroupQuery) (GroupResults, error)
	FindByID(int64) (*Group, error)
	Update(*Group) error
	Touch(*Group) error
}

type GroupQuery struct {
	ProjectID int64
	Resolved  bool
	QueryOptions
}

type GroupResults struct {
	Groups []*Group
	Total  int64
	Query  GroupQuery
}
type groupsStore struct{ *Store }

func (s *groupsStore) FindOrCreate(p *Project, e *Event) (*Group, error) {
	group, err := s.findByProjectIDAndChecksum(p.ID, e.Checksum)
	if err == ErrNotFound {
		group = newGroup(p, e)
		return group, s.insert(group)
	} else if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupsStore) findByProjectIDAndChecksum(id int64, checksum string) (*Group, error) {
	group := &Group{}
	query := "select * from groups where project_id=$1 and checksum=$2 limit 1"
	err := s.Get(group, query, id, checksum)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupsStore) insert(group *Group) error {
	query := "insert into groups (message, checksum, project_id, occurrences, last_seen_at) values($1,$2,$3,$4,$5) returning id"
	err := s.QueryRow(query,
		group.Message,
		group.Checksum,
		group.ProjectID,
		group.Occurrences,
		group.LastSeenAt,
	).Scan(&group.ID)

	if err != nil {
		logger.Error("error inserting group", "err", err)
		return err
	}

	group.WasInserted = true
	return nil
}

func (s *groupsStore) Touch(g *Group) error {
	query := "update groups set occurrences=(select count(*) from errors where errors.checksum = $1), last_seen_at=now_utc(), resolved='f', updated_at=now_utc() where id=$2"
	_, err := s.Exec(query, g.Checksum, g.ID)
	return err
}

func (s *groupsStore) Update(group *Group) error {
	group.UpdatedAt = time.Now().UTC()

	query := "update groups set occurrences=$1, last_seen_at=$2, resolved=$3, updated_at=$4, muted=$5 where id=$6"
	_, err := s.Exec(query,
		group.Occurrences,
		group.LastSeenAt,
		group.Resolved,
		group.UpdatedAt,
		group.Muted,
		group.ID,
	)

	return err
}

func (s *groupsStore) FindQuery(q GroupQuery) (GroupResults, error) {
	groups := GroupResults{}
	groups.Query = q
	groups.Groups = []*Group{}

	countQuery := builder.Select("count(*)").From("groups")
	findQuery := builder.Select("*").From("groups")

	countQuery = countQuery.Where("project_id=?", q.ProjectID)
	findQuery = findQuery.Where("project_id=?", q.ProjectID)

	findQuery = findQuery.Limit(uint64(q.PerPageOrDefault())).Offset(uint64(q.Offset()))

	findQuery = findQuery.OrderBy("last_seen_at desc, created_at desc")

	query, args, _ := findQuery.ToSql()
	err := s.Select(&groups.Groups, query, args...)
	if err != nil {
		return groups, err
	}

	query, args, _ = countQuery.ToSql()
	err = s.Get(&groups.Total, query, args...)
	if err != nil {
		return groups, err
	}

	return groups, err
}

func (s *groupsStore) FindByID(id int64) (*Group, error) {
	group := &Group{}
	query := "select * from groups where id=$1 limit 1"
	err := s.Get(group, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return group, nil
}
