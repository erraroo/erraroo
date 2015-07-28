package models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
)

type jsEvent struct {
	Language  string                 `json:"language"`
	Libaries  []jsEventLibrary       `json:"libaries"`
	Plugins   []plugin               `json:"plugins"`
	Trace     trace                  `json:"trace"`
	URL       string                 `json:"url"`
	UserAgent string                 `json:"userAgent"`
	Version   string                 `json:"version"`
	Processed bool                   `json:"processed"`
	Logs      []jsLog                `json:"logs"`
	Userdata  map[string]interface{} `json:"userdata"`
}

type jsLog struct {
	Level     string  `json:"level"`
	Timestamp float64 `json:"timestamp"`
	Message   string  `json:"message"`
}

type jsEventLibrary struct {
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

type jsErrorEvent struct{ *Event }

func (e *jsErrorEvent) PreCreate() error {
	e.Event.Checksum = e.Checksum()

	resources := NewResourceStore()

	jse, err := e.unmarshal()
	if err != nil {
		return err
	}

	for i := range jse.Trace.Stack {
		err = populateFrameContext(&jse.Trace.Stack[i], resources)
		if err != nil {
			return err
		}
	}

	jse.Processed = true
	payload, err := json.Marshal(jse)
	if err != nil {
		return err
	}

	e.Payload = string(payload)
	return nil
}

func (e *jsErrorEvent) Checksum() string {
	js, _ := e.unmarshal()

	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%d", e.ProjectID))
	io.WriteString(h, js.Trace.Name)
	io.WriteString(h, js.Trace.Message)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (e *jsErrorEvent) unmarshal() (jsEvent, error) {
	var js jsEvent
	err := json.Unmarshal([]byte(e.Payload), &js)
	return js, err
}

func (e *jsErrorEvent) Name() string {
	jse, _ := e.unmarshal()
	return jse.Trace.Name
}

func (e *jsErrorEvent) Message() string {
	jse, _ := e.unmarshal()
	return jse.Trace.Message
}

func (e *jsErrorEvent) Libaries() []Library {
	js, _ := e.unmarshal()

	libs := []Library{}

	for _, l := range js.Libaries {
		libs = append(libs, Library{ProjectID: e.ProjectID, Name: l.Name, Version: l.Version})
	}

	return libs
}

func populateFrameContext(f *frame, resources *resourcesStore) error {
	resource, err := resources.FindByURL(f.URL)
	if err != nil {
		return err
	}

	f.SourceContext = resource.Context(f.Line, f.Column)
	return nil
}
