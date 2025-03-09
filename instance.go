package superhub

import (
	"fmt"
	"net/http"
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

// Instance — конкретный экземпляр предоставленной определённому пользователю услуги. Содержит общее описание услуги,
// используемое в личном кабинете.
type Instance struct {
	// ID — идентификатор услуги в системе.
	ID uuid.UUID `json:"id"`

	// Идентификатор пользователя, являющегося владельцем данной услуги.
	OwnerID uuid.UUID `json:"ownerId"`

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

type ControlPanel struct {
	URL         string `json:"url"`
	DisplayName string `json:"displayName"`
}

// GetPricing получает актуальную информацию о стоимости услуги.
func (i *Instance) GetPricing(client *Client) (*InstancePricing, error) {
	return client.GetInstancePricing(i.ID)
}

// Block блокирует услугу.
func (i *Instance) Block(client *Client) error {
	return client.BlockInstance(i.ID)
}

// Unblock разблокирует услугу. Вернёт ошибку 409, если услуга не заблокирована.
func (i *Instance) Unblock(client *Client) error {
	return client.UnblockInstance(i.ID)
}

// GetInstances получает список всех услуг, доступных в системе.
func (c *Client) GetInstances() (*[]Instance, error) {
	return InvokeEndpoint[[]Instance](c, http.MethodGet, "/instances", nil, nil)
}

// GetInstance получает информацию об услуге с данным идентификатором.
func (c *Client) GetInstance(id uuid.UUID) (*Instance, error) {
	return InvokeEndpoint[Instance](c, http.MethodGet, fmt.Sprintf("/instances/%s", id), nil, nil)
}

// BlockInstance блокирует услугу с заданным идентификатором.
func (c *Client) BlockInstance(id uuid.UUID) error {
	return InvokeVoidEndpoint(c, http.MethodPost, fmt.Sprintf("/instances/%s/blocking", id), nil, nil)
}

// UnblockInstance разблокирует услугу с заданным идентификатором. Вернёт ошибку 409, если услуга не заблокирована.
func (c *Client) UnblockInstance(id uuid.UUID) error {
	return InvokeVoidEndpoint(c, http.MethodDelete, fmt.Sprintf("/instances/%s/blocking", id), nil, nil)
}

// GetInstancePricing получает актуальную информацию о стоимости услуги.
func (c *Client) GetInstancePricing(instanceID uuid.UUID) (*InstancePricing, error) {
	return InvokeEndpoint[InstancePricing](c, http.MethodGet, fmt.Sprintf("/instances/%s/pricing", instanceID), nil, nil)
}

// GetInstanceControlPanel получает информацию о доступе ко внешней панели управления услугой.
func (c *Client) GetInstanceControlPanel(instanceID uuid.UUID) (*ControlPanel, error) {
	return InvokeEndpoint[ControlPanel](c, http.MethodGet, fmt.Sprintf("/instances/%s/control-panel", instanceID), nil, nil)
}
