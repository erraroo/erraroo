package models

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/emailer"
)

type Group struct {
	ID          int64
	Message     string
	Checksum    string
	Occurrences int
	Resolved    bool
	LastSeenAt  time.Time `db:"last_seen_at"`
	ProjectID   int64     `db:"project_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type GroupQueryResults struct {
	Groups []*Group
	Total  int64
}

func newGroup(p *Project, e *Error) *Group {
	return &Group{ProjectID: p.ID, Message: e.Message(), Checksum: e.Checksum}
}

func (g *Group) Touch() error {
	g.LastSeenAt = time.Now().UTC()
	g.Occurrences += 1
	g.Resolved = false
	return Groups.Update(g)
}

func (g *Group) AfterInsert() error {
	err := g.deliverEmailNotifcations()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (g *Group) deliverEmailNotifcations() error {
	project, err := Projects.FindByID(g.ProjectID)
	if err != nil {
		return err
	}

	users, err := Users.ByAccountID(project.AccountID)
	if err != nil {
		return err
	}

	for _, user := range users {
		err := DeliverNewGroupNotification(user, g)
		if err != nil {
			log.Printf("[error] %s\n", err)
			continue
		}
	}

	return nil
}

func DeliverNewGroupNotification(user *User, group *Group) error {
	data := newGroupView(group)
	err := notifyUser(user, data)
	if err != nil {
		log.Printf("[error] %s\n", err)
	}

	return err
}

type groupView struct {
	URL     string
	Message string
	Subject string
}

func newGroupView(g *Group) groupView {
	return groupView{
		Subject: fmt.Sprintf("[erraroo] %s", g.Message),
		Message: g.Message,
		URL:     fmt.Sprintf("%s/projects/%d/groups/%d/errors/latest", config.MailerBaseURL, g.ProjectID, g.ID),
	}
}

func notifyUser(user *User, group groupView) error {
	body := bytes.NewBufferString("")

	tmpl, err := template.New("T").Parse(groupNotifcationEmailTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(body, group)
	if err != nil {
		log.Println(err)
		return err
	}

	err = emailer.Send(user.Email, group.Subject, body.String())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const groupNotifcationEmailTemplate = `
{{define "T"}}
<html>
	<head>
	</head>
	<body>
		<p>Hello Friend!</p>
		<p>We are writing to let you know there was an error in your awesome ember app.</p>
		<p>
			<a href="{{.URL}}">{{.Message}}</a>
		</p>
	</body>
</html>
{{end}}
`
