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
		URL:     fmt.Sprintf("%s/projects/%d/errors/%d/events/latest", config.MailerBaseURL, g.ProjectID, g.ID),
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

func DeliverInvitation(invitation *models.Invitation) error {
	body := bytes.NewBufferString("")

	tmpl, err := template.New("T").Parse(invitationTemplate)
	if err != nil {
		return err
	}

	user, err := models.Users.FindByID(invitation.UserID)
	if err != nil {
		return err
	}

	view := invitationView{
		From:    user,
		URL:     fmt.Sprintf("%s/invitation/%s", config.MailerBaseURL, invitation.Token),
		Subject: "[erraroo] Invitation!",
	}

	err = tmpl.Execute(body, view)
	if err != nil {
		return err
	}

	return emailer.Send(invitation.Address, view.Subject, body.String())
}

type invitationView struct {
	From    *models.User
	URL     string
	Subject string
}

const invitationTemplate = `
{{define "T"}}
<html>
	<head>
	</head>
	<body>
		<p>Hello New Friend!</p>
		<p>You were invited by {{.From.Email}}</p>
		<p>
			<a href="{{.URL}}">{{.URL}}</a>
		</p>
	</body>
</html>
{{end}}
`

func DeliverPasswordRecover(pr *models.PasswordRecover) error {
	body := bytes.NewBufferString("")

	tmpl, err := template.New("T").Parse(passwordRecoverTemplate)
	if err != nil {
		return err
	}

	view := passwordRecoverView{
		URL:     fmt.Sprintf("%s/recover-password/%s", config.MailerBaseURL, pr.Token),
		Subject: "[erraroo] Password Recovery",
	}

	err = tmpl.Execute(body, view)
	if err != nil {
		return err
	}

	return emailer.Send(pr.User.Email, view.Subject, body.String())
}

type passwordRecoverView struct {
	URL     string
	Subject string
}

const passwordRecoverTemplate = `
{{define "T"}}
<html>
	<head>
	</head>
	<body>
		<p>Hello New Friend!</p>
		<p>You just requested a password recovery.  Here you go :)</p>
		<p>
			<a href="{{.URL}}">{{.URL}}</a>
		</p>
	</body>
</html>
{{end}}
`

func DeliverRateLimitNotifcation(user *models.User, project *models.Project) error {
	body := bytes.NewBufferString("")

	tmpl, err := template.New("T").Parse(rateLimitNotifcationTemplate)
	if err != nil {
		return err
	}

	view := rateLimitNotifcationView{
		ProjectName: project.Name,
		Subject:     "[erraroo] Rate Limit Exceeded",
	}

	err = tmpl.Execute(body, view)
	if err != nil {
		return err
	}

	return emailer.Send(user.Email, view.Subject, body.String())
}

type rateLimitNotifcationView struct {
	ProjectName string
	Subject     string
}

const rateLimitNotifcationTemplate = `
{{define "T"}}
<html>
	<head>
	</head>
	<body>
		<p>Hello New Friend!</p>
		<p>Your project {{.ProjectName}} has execeeded it's rate limit.</p>
		<p>I hope your app is ok friend :(</p>
	</body>
</html>
{{end}}
`
