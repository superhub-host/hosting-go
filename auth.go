package superhub

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Credentials interface {
	AuthorizeRequest(request *http.Request)
}

type EmptyCredentials struct{}

func (EmptyCredentials) AuthorizeRequest(*http.Request) {}

type PersistentToken struct {
	ID          uuid.UUID
	AccessToken string
}

func (p *PersistentToken) AuthorizeRequest(request *http.Request) {
	request.Header.Set("Authorization", fmt.Sprintf("%s %s", "Persistent", p.ID))
	request.Header.Set("X-Access-Token", p.AccessToken)
}

func NewPersistentToken(id uuid.UUID, accessToken string) *PersistentToken {
	return &PersistentToken{ID: id, AccessToken: accessToken}
}
