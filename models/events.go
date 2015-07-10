package models

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/erraroo/erraroo/jobs"
	"github.com/erraroo/erraroo/logger"
)

// EventsStore is the interface to error data
type EventsStore interface {
	Create(token, kind, data string) (*Event, error)
	ListForProject(*Project) ([]*Event, error)
	FindByID(int64) (*Event, error)
	FindQuery(EventQuery) (EventResults, error)
	Update(*Event) error
}

type eventsStore struct{ *Store }

func (s *eventsStore) Create(token, kind, data string) (*Event, error) {
	var err error
	project, err := Projects.FindByToken(token)
	if err != nil {
		return nil, err
	}

	e := &Event{}
	e.Kind = kind
	e.Payload = data
	e.ProjectID = project.ID
	e.CreatedAt = time.Now().UTC()
	e.UpdatedAt = e.CreatedAt

	err = e.PreProcess()
	if err != nil {
		logger.Error("error pre processing event", "kind", kind, "payload", data, "token", token, "err", err)
		return nil, err
	}

	switch kind {
	case "js.error":
		err := s.Save(e).Error
		if err != nil {
			logger.Error("inserting event", "err", err)
			return nil, err
		}

		//key := fmt.Sprintf("%d", e.ID)
		//err = put(key, []byte(e.Payload))

		err = jobs.Push("event.process", e.ID)
		if err != nil {
			return nil, err
		}

	case "js.timing":
		_, err := Timings.Create(project, e.Payload)
		if err != nil {
			return nil, err
		}

	}

	return e, err
}

var tableName = "fun"
var svc = dynamodb.New(nil)

func put(key string, payload []byte) error {
	params := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(key),
			},
			"payload": {
				B: payload,
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.PutItem(params)
	return err
}

func get(key string) ([]byte, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(key),
			},
		},
		TableName: aws.String(tableName),
	}

	resp, err := svc.GetItem(params)
	if err != nil {
		return nil, err
	}

	item := resp.Item["payload"]
	return item.B, nil
}

func (s *eventsStore) ListForProject(p *Project) ([]*Event, error) {
	events := []*Event{}
	query := s.Where("project_id=?", p.ID).Order("created_at desc").Limit(100)
	return events, query.Find(&events).Error
}

func (s *eventsStore) FindByID(id int64) (*Event, error) {
	e := &Event{}
	o := s.First(&e, id)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *eventsStore) Update(e *Event) error {
	err := s.Save(e).Error
	if err != nil {
		return err
	}

	return err
}

type EventQuery struct {
	ProjectID int64
	Checksum  string
	Kind      string
	QueryOptions
}

type EventResults struct {
	Events []*Event
	Total  int64
	Query  EventQuery
}

func (s *eventsStore) FindQuery(q EventQuery) (EventResults, error) {
	events := EventResults{
		Query:  q,
		Events: []*Event{},
	}

	scope := s.Table("events")
	scope = scope.Where("events.project_id=?", q.ProjectID)

	if q.Checksum != "" {
		scope = scope.Where("events.checksum=?", q.Checksum)
	}

	if q.Kind != "" {
		scope = scope.Where("events.kind=?", q.Kind)
	}

	o := scope.Count(&events.Total)
	if o.Error != nil {
		return events, o.Error
	}

	scope = scope.Limit(q.PerPageOrDefault()).Offset(q.Offset())
	scope = scope.Order("events.created_at desc")

	o = scope.Find(&events.Events)
	if o.Error != nil {
		return events, o.Error
	}

	return events, nil
}
