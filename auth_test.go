package superhub

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func TestPersistentToken_AuthorizeRequest(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
		return
	}

	bytes := make([]byte, 128)
	_, err = rand.Read(bytes)
	if err != nil {
		t.Error(err)
		return
	}

	token := base64.StdEncoding.EncodeToString(bytes)
	wrappedToken := NewPersistentToken(id, token)

	request, err := http.NewRequest("GET", "https://google.com/", nil)
	if err != nil {
		t.Error(err)
		return
	}

	wrappedToken.AuthorizeRequest(request)

	assert.Equal(t, request.Header.Get("Authorization"), fmt.Sprintf("Persistent %s", id.String()))
	assert.Equal(t, request.Header.Get("X-Access-Token"), token)
}
