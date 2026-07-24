package superhub

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// InstanceStatus — состояние услуги. Значение отражает то, доступна ли услуга в данный момент, или причину
// недоступности.
type InstanceStatus string

const (
	InstanceStatusActive        InstanceStatus = "ACTIVE"
	InstanceStatusBlocked       InstanceStatus = "BLOCKED"
	InstanceStatusInstallFailed InstanceStatus = "INSTALL_FAILED"
	InstanceStatusInstalling    InstanceStatus = "INSTALLING"
	InstanceStatusUserSuspended InstanceStatus = "USER_SUSPENDED"
)

type BlockingReason string

const (
	BlockingReasonUnknown      BlockingReason = "UNKNOWN"
	BlockingReasonUnpaid       BlockingReason = "UNPAID"
	BlockingReasonTosViolation BlockingReason = "TOS_VIOLATION"
)

// Instance — конкретный экземпляр предоставленной определённому пользователю услуги. Содержит общее описание услуги,
// используемое в личном кабинете.
type Instance struct {
	// ID — идентификатор услуги в системе.
	ID uuid.UUID `json:"id"`

	// Идентификатор пользователя, являющегося владельцем данной услуги.
	OwnerID uuid.UUID `json:"ownerId"`

	// Идентификатор сервера для технической поддержки.
	Identifier string `json:"identifier"`

	// Название услуги, предоставленное пользователем.
	Name string `json:"name"`

	// Описание услуги, предоставленное пользователем.
	Description string `json:"description"`

	// Текущий тариф услуги.
	Tariff Tariff `json:"tariff"`

	// Динамические параметры, применённые пользователем при заказе.
	Options []InstanceOption `json:"options"`

	// Тип политики ценообразования, применяемый при вычислении стоимости.
	PricingPolicy PricingPolicyType `json:"pricingPolicy"`

	// Текущая стоимость услуги без учёта скидок.
	CurrentCost Amount `json:"currentCost"`

	// Тип услуги.
	InstanceKind InstanceKind `json:"instanceKind"`

	// Группа, к которой относится услуга.
	ServiceGroup ServiceGroup `json:"serviceGroup"`

	// Идентификатор пользователя, с баланса которого будет списываться плата за услугу.
	InvoiceReceiverID uuid.UUID `json:"invoiceReceiverId"`

	// Идентификатор группы тарифов, в которой находится услуга.
	TariffGroupID uuid.UUID `json:"tariffGroupId"`

	Status InstanceStatus `json:"status"`

	// Дата создания услуги.
	CreatedAt time.Time `json:"createdAt"`

	// Заблокирована ли услуга принудительно?
	Blocked bool `json:"blocked"`

	// Причина блокировки услуги.
	BlockingReason BlockingReason `json:"blockingReason"`

	// Может ли услуга быть разморожена?
	CanBeUnblocked bool `json:"canBeUnblocked"`

	// Время принудительной блокировки услуги. Будет nil, если услуга не заблокирована принудительно.
	BlockedAt *time.Time `json:"blockedAt"`

	// Заблокирована ли услуга добровольно?
	UserSuspended bool `json:"userSuspended"`

	// Время добровольной блокировки услуги. Будет nil, если услуга не заблокирована добровольно.
	UserSuspendedAt *time.Time `json:"userSuspendedAt"`
}

type InstanceOption struct {
	ID          uuid.UUID   `json:"id"`
	OrderOption OrderOption `json:"orderOption"`
	Quantity    float64     `json:"quantity"`
	Cost        Amount      `json:"cost"`
}

// PricingPolicyType — тип политики ценообразования. Показывает, какой механизм использует система для подсчёта
// стоимости конкретной услуги, приобретённой пользователем.
type PricingPolicyType string

const (
	// FixedPricingPolicy (фиксированная политика ценообразования).
	// Стоимость услуги фиксируется при покупке и не зависит от каких-либо других параметров.
	FixedPricingPolicy = "FIXED"

	// FrozenDiskUsageBasedPricingPolicy (политика, основанная на занятом месте на диске).
	// Стоимость услуги зависит от использованного услугой места на диске.
	FrozenDiskUsageBasedPricingPolicy = "FIXED"
)

// InstancePricing — параметры, связанные с выставлением счетов владельцу услуги.
type InstancePricing struct {
	// Ценовая политика, использованная при вычислении стоимости.
	PricingPolicy PricingPolicyType `json:"pricingPolicy"`

	// Текущая стоимость с учётом скидок
	ActualCost Amount `json:"actualCost"`

	// Период оплаты услуги.
	BillingPeriod BillingPeriod `json:"billingPeriod"`
}

type InstanceState struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	KindName    string          `json:"kindName"`
	Options     OptionContainer `json:"options"`
}

type ControlPanel struct {
	URL         string `json:"url"`
	DisplayName string `json:"displayName"`
}

type UpdateTariffParams struct {
	TariffID uuid.UUID         `json:"tariffId"`
	Options  map[string]string `json:"options"`
}

type TariffParams struct {
	Tariff  Tariff           `json:"tariff"`
	Options []InstanceOption `json:"options"`
}

type BlockInstanceParams struct {
	BlockingReason BlockingReason `json:"blockingReason"`
}

// GetPricing получает актуальную информацию о стоимости услуги.
func (i *Instance) GetPricing(client *Client) (*InstancePricing, error) {
	return client.GetInstancePricing(i.ID.String())
}

// GetState получает актуальную конфигурацию услуги во внешней системе.
func (i *Instance) GetState(client *Client) (*InstanceState, error) {
	return client.GetInstanceState(i.ID.String())
}

// Block блокирует услугу.
func (i *Instance) Block(client *Client, params *BlockInstanceParams) error {
	return client.BlockInstance(i.ID.String(), params)
}

// Unblock разблокирует услугу. Вернёт ошибку 409, если услуга не заблокирована.
func (i *Instance) Unblock(client *Client) error {
	return client.UnblockInstance(i.ID.String())
}

// UpdateTariff изменяет тариф услуги.
func (i *Instance) UpdateTariff(client *Client, params *UpdateTariffParams) (*TariffParams, error) {
	return client.UpdateInstanceTariff(i.ID.String(), params)
}

// GetPayments получает платежи, свяазнные с услугой.
func (i *Instance) GetPayments(client *Client, params *PaginationParams) (*[]Payment, error) {
	return client.GetInstancePayments(i.ID.String(), params)
}

// Delete удаляет услугу без подтверждения.
func (i *Instance) Delete(client *Client) error {
	return client.DeleteInstance(i.ID.String())
}

// GetInstances получает список всех услуг, доступных в системе.
func (c *Client) GetInstances() (*[]Instance, error) {
	return InvokeEndpoint[[]Instance](c, http.MethodGet, "/instances", nil, nil)
}

// GetInstance получает информацию об услуге с данным идентификатором.
func (c *Client) GetInstance(id string) (*Instance, error) {
	return InvokeEndpoint[Instance](c, http.MethodGet, fmt.Sprintf("/instances/%s", id), nil, nil)
}

// BlockInstance блокирует услугу с заданным идентификатором.
func (c *Client) BlockInstance(id string, params *BlockInstanceParams) error {
	return InvokeVoidEndpoint(c, http.MethodPost, fmt.Sprintf("/instances/%s/blocking", id), nil, params)
}

// UnblockInstance разблокирует услугу с заданным идентификатором. Вернёт ошибку 409, если услуга не заблокирована.
func (c *Client) UnblockInstance(id string) error {
	return InvokeVoidEndpoint(c, http.MethodDelete, fmt.Sprintf("/instances/%s/blocking", id), nil, nil)
}

// DeleteInstance удаляет услугу без подтверждения.
func (c *Client) DeleteInstance(id string) error {
	return InvokeVoidEndpoint(c, http.MethodDelete, fmt.Sprintf("/instances/%s", id), nil, nil)
}

// GetInstancePricing получает актуальную информацию о стоимости услуги.
func (c *Client) GetInstancePricing(instanceID string) (*InstancePricing, error) {
	return InvokeEndpoint[InstancePricing](c, http.MethodGet, fmt.Sprintf("/instances/%s/pricing", instanceID), nil, nil)
}

// GetInstancePayments получает платежи, свяазнные с услугой.
func (c *Client) GetInstancePayments(instanceID string, params *PaginationParams) (*[]Payment, error) {
	query := &url.Values{}
	params.Encode(query)

	return InvokeEndpoint[[]Payment](c, http.MethodGet, fmt.Sprintf("/instances/%s/payments", instanceID), query, nil)
}

// GetInstanceControlPanel получает информацию о доступе ко внешней панели управления услугой.
func (c *Client) GetInstanceControlPanel(instanceID string) (*ControlPanel, error) {
	return InvokeEndpoint[ControlPanel](c, http.MethodGet, fmt.Sprintf("/instances/%s/control-panel", instanceID), nil, nil)
}

// GetInstanceState получает актуальную конфигурацию услуги во внешней системе.
func (c *Client) GetInstanceState(instanceID string) (*InstanceState, error) {
	return InvokeEndpoint[InstanceState](c, http.MethodGet, fmt.Sprintf("/instances/%s/state", instanceID), nil, nil)
}

// UpdateInstanceTariff изменяет тариф услуги.
func (c *Client) UpdateInstanceTariff(instanceID string, params *UpdateTariffParams) (*TariffParams, error) {
	return InvokeEndpoint[TariffParams](c, http.MethodPut, fmt.Sprintf("/instances/%s/tariff", instanceID), nil, params)
}
