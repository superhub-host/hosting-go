package superhub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Price struct {
	FullAmount   Amount `json:"fullAmount"`
	ActualAmount Amount `json:"actualAmount"`
}

type Tariff struct {
	ID                uuid.UUID         `json:"id"`
	DisplayName       string            `json:"displayName"`
	Description       string            `json:"description"`
	Price             Price             `json:"price"`
	Options           OptionContainer   `json:"options"`
	MaxPerUser        *int              `json:"maxPerUser"`
	VerificationLevel VerificationLevel `json:"verificationLevel"`
	UserCanSuspend    bool              `json:"userCanSuspend"`
	DomainLimit       int               `json:"domainLimit"`
	Public            bool              `json:"public"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

type OrderOption struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
	Price       Price     `json:"price"`
	OptionKey   string    `json:"optionKey"`
	MinQuantity float64   `json:"minQuantity"`
	MaxQuantity float64   `json:"maxQuantity"`
}

type BillingPeriod struct {
	ID            uuid.UUID `json:"id"`
	DisplayName   string    `json:"displayName"`
	Description   string    `json:"description"`
	PeriodSeconds *int      `json:"periodSeconds"`
}

type TariffSet struct {
	ID            uuid.UUID         `json:"id"`
	ServiceGroup  ServiceGroup      `json:"serviceGroup"`
	PricingPolicy PricingPolicyType `json:"pricingPolicy"`
	DisplayName   string            `json:"displayName"`
	Description   string            `json:"description"`
}

type TariffGroup struct {
	ID            uuid.UUID       `json:"id"`
	InstanceKind  InstanceKind    `json:"instanceKind"`
	TariffSet     TariffSet       `json:"tariffSet"`
	BillingPeriod BillingPeriod   `json:"billingPeriod"`
	DisplayName   string          `json:"displayName"`
	Description   string          `json:"description"`
	Options       OptionContainer `json:"options"`
}

func (g *TariffGroup) GetTariffs(client *Client) (*[]Tariff, error) {
	return client.GetTariffs(g.ID)
}

func (g *TariffGroup) GetOrderOptions(client *Client) (*[]OrderOption, error) {
	return client.GetOrderOptions(g.ID)
}

func (c *Client) GetTariffGroups() (*[]TariffGroup, error) {
	return InvokeEndpoint[[]TariffGroup](c, http.MethodGet, "/tariff-groups", nil, nil)
}

func (c *Client) GetTariffGroup(tariffGroupId uuid.UUID) (*TariffGroup, error) {
	return InvokeEndpoint[TariffGroup](c, http.MethodGet, fmt.Sprintf("/tariff-groups/%s", tariffGroupId), nil, nil)
}

func (c *Client) GetTariffs(tariffGroupId uuid.UUID) (*[]Tariff, error) {
	return InvokeEndpoint[[]Tariff](c, http.MethodGet, fmt.Sprintf("/tariff-groups/%s/tariffs", tariffGroupId), nil, nil)
}

func (c *Client) GetOrderOptions(tariffGroupId uuid.UUID) (*[]OrderOption, error) {
	return InvokeEndpoint[[]OrderOption](c, http.MethodGet, fmt.Sprintf("/tariff-groups/%s/order-options", tariffGroupId), nil, nil)
}
