package emailer

import (
	"log"

	"github.com/erraroo/erraroo/config"
	"github.com/sendgrid/sendgrid-go"
)

var sender Sender

func init() {
	Use(&sendGridSender{
		config.SendGridKey,
		config.SendGridUser,
	})
}

func Use(s Sender) {
	sender = s
}

type Sender interface {
	Send(to, subject, body string) error
}

func Send(to string, subject string, body string) error {
	return sender.Send(to, subject, body)
}

type sendGridSender struct {
	key  string
	user string
}

func (s *sendGridSender) Send(to, subject, body string) error {
	sg := sendgrid.NewSendGridClient(config.SendGridUser, config.SendGridKey)
	message := sendgrid.NewMail()
	message.AddTo(to)
	message.SetFrom("ben@nerdyworm.com")
	message.SetSubject(subject)
	message.SetHTML(body)
	if err := sg.Send(message); err != nil {
		log.Printf("[sendGridSender] [error] Send to=%s subject=%s error=%2\n", to, subject, err)
		return err
	}

	return nil
}

type DummySender struct{}

func (n *DummySender) Send(to, subject, body string) error {
	return nil
}
