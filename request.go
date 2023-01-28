package superhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func InvokeEndpoint[T any](client *Client, method, path string, body io.Reader) (*T, error) {
	url, err := client.GetEndpointURL(path)
	if err != nil {
		return nil, fmt.Errorf("making endpoint URL: %s", err)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP request: %s", err)
	}

	return ProcessRequest[T](client, request)
}

func ProcessRequest[T any](client *Client, request *http.Request) (*T, error) {
	client.Credentials.AuthorizeRequest(request)

	response, err := client.GetHttpClient().Do(request)
	if err != nil {
		return nil, fmt.Errorf("dispatching request: %s", err)
	}

	data, err := handleResponse[T](response)
	if err != nil {
		return nil, fmt.Errorf("handling response: %s", err)
	}

	return data, nil
}

type ErrorResponse struct {
	ErrorName string    `json:"error"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("request error: %d (%s) on path %s: %s", e.Status, e.ErrorName, e.Path, e.Message)
}

func handleResponse[T any](response *http.Response) (*T, error) {
	if response.StatusCode >= http.StatusBadRequest && response.StatusCode < http.StatusInternalServerError {
		errorResponse, err := handleErrorResponse(response)
		if err != nil {
			return nil, fmt.Errorf("unsuccessful HTTP status: %s", response.Status)
		}

		return nil, errorResponse
	}

	if response.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("internal server error: %s", response.Status)
	}

	return parseResponse[T](response)
}

func handleErrorResponse(response *http.Response) (*ErrorResponse, error) {
	return parseResponse[ErrorResponse](response)
}

func parseResponse[T any](response *http.Response) (*T, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %s", err)
	}

	if contentType := response.Header.Get("Content-Type"); contentType != "application/json" {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	var data T
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("parsing response body: %s", err)
	}

	return &data, nil
}
