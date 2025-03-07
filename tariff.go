package superhub

import (
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
	Options           []OptionValue     `json:"options"`
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
