package models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

// Error is the entity that stores error data
type Error struct {
	ID        int64
	Payload   string
	Checksum  string
	ProjectID int64     `db:"project_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type jsError struct {
	Language  string    `json:"language"`
	Libaries  []library `json:"libaries"`
	Plugins   []plugin  `json:"plugins"`
	Trace     trace     `json:"trace"`
	URL       string    `json:"url"`
	UserAgent string    `json:"userAgent"`
	Version   string    `json:"version"`
}

type library struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type plugin struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type trace struct {
	Incomplete bool    `json:"incomplete"`
	Message    string  `json:"message"`
	Mode       string  `json:"mode"`
	Name       string  `json:"name"`
	Stack      []frame `json:"stack"`
}

type frame struct {
	Column int    `json:"column"`
	Func   string `json:"func"`
	Line   int    `json:"line"`
	URL    string `json:"url"`
	SourceContext
}

type SourceContext struct {
	PreContext  []string `json:"preContext"`
	ContextLine string   `json:"contextLine"`
	PostContext []string `json:"postContext"`
	OrigLine    int      `json:"origLine"`
	OrigColumn  int      `json:"origColumn"`
	OrigFile    string   `json:"origFile"`
}

func (e Error) unmarshalJSError() (jsError, error) {
	var js jsError
	err := json.Unmarshal([]byte(e.Payload), &js)
	return js, err
}

func (e *Error) generateChecksum() {
	js, _ := e.unmarshalJSError()

	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%d", e.ProjectID))
	io.WriteString(h, js.Trace.Message)

	e.Checksum = fmt.Sprintf("%x", h.Sum(nil))
}

func (e *Error) PopulateStackContext(resources *resourcesStore) error {
	jse, err := e.unmarshalJSError()
	if err != nil {
		return err
	}

	for i := range jse.Trace.Stack {
		populateFrameContext(&jse.Trace.Stack[i], resources)
	}

	payload, err := json.Marshal(jse)
	if err != nil {
		return err
	}

	e.Payload = string(payload)

	return nil
}

func populateFrameContext(f *frame, resources *resourcesStore) {
	resource, err := resources.FindByURL(f.URL)
	if err != nil {
		log.Printf("error fetching %s `%v`\n", f.URL, err)
		return
	}

	f.SourceContext = resource.Context(f.Line, f.Column)
}

func (e *Error) Message() string {
	jse, _ := e.unmarshalJSError()
	return jse.Trace.Message
}