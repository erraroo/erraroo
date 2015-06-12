package models

import (
	"database/sql"
	"log"
	"time"
)

// ErrorsStore is the interface to error data
type GroupsStore interface {
	FindOrCreate(*Project, *Error) (*Group, error)
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

func (s *groupsStore) FindOrCreate(p *Project, e *Error) (*Group, error) {
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
		log.Println(err)
		return err
	}

	return group.AfterInsert()
}

func (s *groupsStore) Touch(g *Group) error {
	query := "update groups set occurrences=occurrences+1, last_seen_at=(select now() at time zone 'utc'), resolved='f', updated_at=now() where id=$1"
	_, err := s.Exec(query, g.ID)
	return err
}

func (s *groupsStore) Update(group *Group) error {
	group.UpdatedAt = time.Now().UTC()

	query := "update groups set occurrences=$1, last_seen_at=$2, resolved=$3, updated_at=$4 where id=$5"
	_, err := s.Exec(query,
		group.Occurrences,
		group.LastSeenAt,
		group.Resolved,
		group.UpdatedAt,
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
