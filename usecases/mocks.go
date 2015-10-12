package usecases

import (
	"fmt"
	"testing"
	"time"

	"github.com/erraroo/erraroo/models"
	"github.com/stretchr/testify/assert"
	"github.com/tuvistavie/securerandom"
)

type mockSender struct {
	sends []map[string]string
}

func (m *mockSender) Send(to, subject, body string) error {
	payload := map[string]string{
		"to":      to,
		"subject": subject,
		"body":    body,
	}
	m.sends = append(m.sends, payload)
	return nil
}

func (m *mockSender) Clear() {
	m.sends = []map[string]string{}
}

func uniqEmail() string {
	return fmt.Sprintf("%d@example.com", time.Now().Nanosecond())
}

func uuid() string {
	uuid, _ := securerandom.Uuid()
	return uuid
}

func aup(t *testing.T) (*models.Account, *models.User, *models.Project) {
	account := account(t)
	user := user(t, account)
	project := project(t, account)
	return account, user, project
}

func account(t *testing.T) *models.Account {
	account, err := models.CreateAccount()
	assert.Nil(t, err)
	assert.NotNil(t, account)
	return account
}

func project(t *testing.T, account *models.Account) *models.Project {
	project, err := models.Projects.Create("test", account.ID)
	assert.Nil(t, err)
	assert.NotNil(t, project)
	return project
}

func user(t *testing.T, account *models.Account) *models.User {
	email := uniqEmail()
	user, err := models.Users.Create(email, "password", account)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	return user
}

func makeEvent(t *testing.T, project *models.Project, payload string) *models.Event {
	event := models.NewEvent(project, "js.error", payload)
	err := models.Events.Insert(event)
	assert.Nil(t, err)
	assert.NotNil(t, event)
	return event
}
