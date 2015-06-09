package cx

import (
	"github.com/erraroo/erraroo/models"
	"github.com/nerdyworm/rsq"
)

// Ctx is the context that is passed into each request
type Context struct {
	User  *models.User
	Queue rsq.Queue
}
