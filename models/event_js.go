package models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/erraroo/erraroo/logger"
)

type jsEvent struct {
	Language  string           `json:"language"`
	Libaries  []jsEventLibrary `json:"libaries"`
	Plugins   []plugin         `json:"plugins"`
	Trace     trace            `json:"trace"`
	URL       string           `json:"url"`
	UserAgent string           `json:"userAgent"`
	Version   string           `json:"version"`
	Processed bool             `json:"processed"`
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

func (e *jsErrorEvent) IsAsync() bool {
	return true
}

func (e *jsErrorEvent) PreProcess() error {
	e.Event.Checksum = e.Checksum()
	return nil
}

func (e *jsErrorEvent) PostProcess() error {
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

	key := fmt.Sprintf("%d", e.ID)
	err = put(key, payload)
	if err != nil {
		return err
	}

	//e.Payload = string(payload)
	err = Events.Update(e.Event)
	if err != nil {
		logger.Error("updating event", "err", err, "event", e.ID)
		return err
	}

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

	payload, err := get(fmt.Sprintf("%d", e.ID))
	log.Println("FETCH")
	if err != nil {
		return js, err
	}

	return js, json.Unmarshal(payload, &js)

	//err := json.Unmarshal([]byte(e.Payload), &js)
	//return js, err
}

func (e *jsErrorEvent) Name() string {
	jse, _ := e.unmarshal()
	return jse.Trace.Name
}

func (e *jsErrorEvent) Message() string {
	jse, _ := e.unmarshal()
	return jse.Trace.Message
}

func (e *jsErrorEvent) Tags() []Tag {
	js, _ := e.unmarshal()

	tags := []Tag{}

	//for _, l := range js.Libaries {
	//tags = append(tags, l.Tag())
	//}

	if js.UserAgent != "" {
		tags = append(tags, Tag{
			Key:   "js.useragent",
			Value: js.UserAgent,
			Label: "UserAgent",
		})
	}

	if js.URL != "" {
		tags = append(tags, Tag{
			Key:   "js.url",
			Value: js.URL,
			Label: "URL",
		})
	}

	return tags
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

type jsTimingEvent struct{ *Event }

func (e *jsTimingEvent) IsAsync() bool {
	return false
}

func (e *jsTimingEvent) PreProcess() error {
	return nil
}

func (e *jsTimingEvent) PostProcess() error {
	return nil
}

func (e *jsTimingEvent) Checksum() string {
	return ""
}

func (e *jsTimingEvent) Name() string {
	return "timing event"
}

func (e *jsTimingEvent) Message() string {
	return "timing event recorded"
}

func (e *jsTimingEvent) Tags() []Tag {
	return []Tag{}
}

func (e *jsTimingEvent) Libaries() []Library {
	return []Library{}
}

type jsLogEvent struct{ *Event }

func (e *jsLogEvent) IsAsync() bool {
	return false
}

func (e *jsLogEvent) PreProcess() error {
	logger.Info("js.log", "payload", e.Payload)
	return nil
}

func (e *jsLogEvent) PostProcess() error {
	return nil
}

func (e *jsLogEvent) Checksum() string {
	return ""
}

func (e *jsLogEvent) Name() string {
	return "log event"
}

func (e *jsLogEvent) Message() string {
	return "log event"
}

func (e *jsLogEvent) Tags() []Tag {
	return []Tag{}
}

func (e *jsLogEvent) Libaries() []Library {
	return []Library{}
}
