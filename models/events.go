package models

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/logger"
)

// EventsStore is the interface to event data
type EventsStore interface {
	FindByID(int64) (*Event, error)
	FindQuery(EventQuery) (EventResults, error)
	Insert(*Event) error
	PayloadURL(*Event) string
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

type eventsStore struct {
	*Store
	service *s3.S3
}

func (s *eventsStore) FindByID(id int64) (*Event, error) {
	e := &Event{}
	o := s.First(&e, id)
	if o.RecordNotFound() {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *eventsStore) Insert(e *Event) error {
	err := e.BeforeCreate()
	if err != nil {
		logger.Error("error running before create on event", "err", err)
		return err
	}

	query := "insert into events (checksum, kind, project_id) values($1,$2,$3) returning id, created_at"
	err = s.QueryRow(query, e.Checksum, e.Kind, e.ProjectID).Scan(
		&e.ID,
		&e.CreatedAt,
	)

	return s.putPayload(e)
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

func (s *eventsStore) putPayload(e *Event) error {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	io.Copy(w, strings.NewReader(e.Payload))
	w.Close()

	params := &s3.PutObjectInput{
		Bucket:          aws.String(config.Bucket),
		Body:            bytes.NewReader(b.Bytes()),
		ContentEncoding: aws.String("gzip"),
		ContentType:     aws.String("application/json"),
		Key:             aws.String(e.PayloadKey()),
	}

	_, err := s.service.PutObject(params)
	if err != nil {
		logger.Error("error saving payload to s3", "err", err)
		return err
	}

	return err
}

func (s *eventsStore) PayloadURL(e *Event) string {
	req, _ := s.service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String(e.PayloadKey()),
	})

	url, err := req.Presign(300 * time.Second)
	if err != nil {
		panic(err)
	}

	return url
}
