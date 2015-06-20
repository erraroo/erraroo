package mailers

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/erraroo/erraroo/config"
	"github.com/erraroo/erraroo/emailer"
	"github.com/erraroo/erraroo/logger"
	"github.com/erraroo/erraroo/models"
)

func DeliverNewErrorNotification(user *models.User, group *models.Error) error {
	return notifyUser(user, newErrorEmailView(group))
}

type groupEmailView struct {
	URL     string
	Message string
	Subject string
}

func newErrorEmailView(g *models.Error) groupEmailView {
	return groupEmailView{
		Subject: fmt.Sprintf("[erraroo] %s", g.Message),
		Message: g.Message,
		URL:     fmt.Sprintf("%s/projects/%d/groups/%d/errors/latest", config.MailerBaseURL, g.ProjectID, g.ID),
	}
}

func notifyUser(user *models.User, group groupEmailView) error {
	body := bytes.NewBufferString("")

	tmpl, err := template.New("T").Parse(groupNotifcationEmailTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(body, group)
	if err != nil {
		return err
	}

	err = emailer.Send(user.Email, group.Subject, body.String())
	if err != nil {
		logger.Error("could not send", "err", err, "email", user.Email)
		return err
	}

	logger.Info("delivered", "name", "groups.new", "email", user.Email, "url", group.URL)
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
