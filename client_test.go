package superhub

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.Equal(t, client.GetBaseURL(), DefaultBaseURL)
	assert.Equal(t, client.GetCredentials(), &EmptyCredentials{})
}

func TestNewClientWithCredentials(t *testing.T) {
	client := NewClientWithCredentials(&PersistentToken{})
	assert.Equal(t, client.GetBaseURL(), DefaultBaseURL)
	assert.Equal(t, client.GetCredentials(), &PersistentToken{})
}

func TestClient_GetEndpointURL(t *testing.T) {
	baseURLs := []string{
		"https://api.superhub.host/v3", "https://api.superhub.host/v3/",
		"http://127.0.0.1:8080/v3", "http://127.0.0.1:8080/v3/",
		"http://127.0.0.1:8080", "http://127.0.0.1:8080/",
	}

	endpoints := []string{
		"/hello-world", "hello-world", "hello-world/", "/hello-world/",
		"/hello/world", "hello/world", "hello/world/", "/hello/world/",
	}

	expected := []string{
		"https://api.superhub.host/v3/hello-world",
		"https://api.superhub.host/v3/hello/world",
		"http://127.0.0.1:8080/v3/hello-world",
		"http://127.0.0.1:8080/v3/hello/world",
		"http://127.0.0.1:8080/hello-world",
		"http://127.0.0.1:8080/hello/world",
	}

	for i, currentBaseURL := range baseURLs {
		for j, currentEndpoint := range endpoints {
			currentlyExpected := expected[i-i%2+j/4]

			client := &Client{BaseURL: currentBaseURL}
			endpointURL, err := client.GetEndpointURL(currentEndpoint)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, endpointURL, currentlyExpected)
		}
	}
}
