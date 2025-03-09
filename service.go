package superhub

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type InstanceKind struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
}

type ServiceGroup struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
}

func (c *Client) GetServiceGroups() (*ServiceGroup, error) {
	return InvokeEndpoint[ServiceGroup](c, http.MethodGet, "/service-groups", nil, nil)
}

func (c *Client) GetServiceGroup(serviceGroupId uuid.UUID) (*[]ServiceGroup, error) {
	return InvokeEndpoint[[]ServiceGroup](c, http.MethodGet, fmt.Sprintf("/service-groups/%s", serviceGroupId), nil, nil)
}
