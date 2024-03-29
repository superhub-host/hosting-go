package superhub

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v4"
)

// ServerState - состояние сервера. Значение отражает то, готов ли сервер к работе после установки. При этом, значение
// не зависит от того, доступен ли сервер в данный момент, если установка уже прошла.
type ServerState string

const (
	// ServerStateError используется, когда системе не удалось полностью создать сервер.
	ServerStateError ServerState = "ERROR"

	// ServerStateInstalling - стандартное значение состояния, которое означает, что сервер находится в процессе создания.
	ServerStateInstalling ServerState = "INSTALLING"

	// ServerStateReady используется, когда сервер готов к работе.
	ServerStateReady ServerState = "READY"
)

// Server - сервер, принадлежащий определённому пользователю. Содержит только общую информацию, используемую в ЛК,
// такую как стоимость, скидка, дата создания и т.п.
// Более подробная информация, например, о лимитах ресурсов, доступна в структуре внешнего сервера (см. ExternalServer)
type Server struct {
	// ID - идентификатор сервера в системе.
	// Не имеет ничего общего с идентификатором в панели Pterodactyl.
	ID int64 `json:"id"`

	// Идентификатор пользователя, являющегося владельцем данного сервера.
	OwnerID int64 `json:"ownerId"`

	// Текущее состояние сервера (см. ServerState)
	State ServerState `json:"state"`

	// Стоимость сервера в рублях.
	Cost ServerCost `json:"cost"`

	// Параметры, связанные с биллингом сервера.
	Billing ServerBillingConfig `json:"billing"`

	// Параметры, связанные с доменом сервера.
	Domain ServerDomainConfig `json:"domain"`

	// ExternalServer - структура, содержащая информацию о внешнем сервере. Имеет значение nil, если конкретная
	// конечная точка, от которой была получена информация о сервере, не подразумевает получение информации о
	// внешнем сервере.
	ExternalServer *ExternalServer `json:"externalServer"`

	// Дата окончания срока действия сервера. Не имеет ничего общего с датой блокировки сервера за неуплату.
	// Обычно имеет значение только на бесплатных серверах, срок действия которых ограничен двумя днями.
	ExpiresAt null.Time `json:"expiresAt"`

	// Дата заморозки сервера пользователем. Имеет значение только когда пользователь заморозил сервер сам
	// и имеет пустое значение, если сервер заблокирован за неуплату.
	FrozenAt null.Time `json:"frozenAt"`

	// Дата создания сервера.
	CreatedAt time.Time `json:"createdAt"`

	// Дата последнего обновления информации о сервере.
	UpdatedAt time.Time `json:"updatedAt"`
}

// ServerCost показывает текущую стоимость сервера. Включает в себя базовую стоимость и стоимость заморозки.
type ServerCost struct {
	// Базовая стоимость сервера.
	Base float64 `json:"base"`

	// Стоимость заморозки сервера.
	Freeze null.Float `json:"freeze"`
}

// ServerTariffMode - тарифный режим сервера. Отражает одновременно период списаний и способ тарификации.
type ServerTariffMode string

const (
	// TariffModeDailyResources - режим, при котором пользователь выбирает произвольную конфигурацию сервера на своё
	// усмотрение, а списание производится ежедневно.
	TariffModeDailyResources ServerTariffMode = "DAILY_RESOURCES"

	// TariffModeMonthlyTariff - режим, при котором пользователь выбирает готовый тариф, а списание производится
	// ежемесячно.
	TariffModeMonthlyTariff ServerTariffMode = "MONTHLY_TARIFF"
)

// TariffModeBase - основа тарифного режима. Отражает способ тарификации услуги.
type TariffModeBase string

const (
	// TariffModeBaseTariff используется на серверах, где пользователь выбирает тариф из заранее заготовленного списка.
	TariffModeBaseTariff = "TARIFF"

	// TariffModeBaseResourceLimits используется на серверах, где пользователь выбирает произвольную конфигурацию.
	TariffModeBaseResourceLimits = "RESOURCE_LIMITS"
)

// BillingPeriod - период выставления счетов на оплату сервера.
type BillingPeriod string

const (
	// BillingPeriodOnce - единоразовая оплата.
	BillingPeriodOnce = "ONCE"

	// BillingPeriodDaily - ежедневная оплата.
	BillingPeriodDaily = "DAILY"

	// BillingPeriodMonthly - ежемесячная оплата.
	BillingPeriodMonthly = "MONTHLY"
)

// PricingPolicyType - тип политики ценообразования. Показывает, какой механизм использует система для подсчёта
// стоимости конкретной услуги, приобретённой пользователем.
type PricingPolicyType string

const (
	// FixedPricingPolicy (фиксированная политика ценообразования).
	// Стоимость услуги фиксируется при покупке и не зависит от каких-либо других параметров.
	FixedPricingPolicy = "FIXED"
)

// ServerBillingConfig - параметры, связанные с выставлением счетов владельцу сервера.
type ServerBillingConfig struct {
	// Ценовая политика, используемая на сервере.
	PricingPolicy PricingPolicyType `json:"pricingPolicy"`

	// Текущий тарифный режим сервера.
	TariffMode ServerTariffMode `json:"tariffMode"`

	// Идентификатор тарифа, который используется на сервере в данный момент. Имеет значение только если пользователь
	// выбрал готовый тариф, а не свою конфигурацию.
	TariffID null.String `json:"tariffId"`

	// Основание текущего тарифного режима сервера.
	Base TariffModeBase `json:"base"`

	// Период выставления счетов владельцу сервера.
	Period BillingPeriod `json:"period"`
}

// ServerDomainConfig содержит информацию о домене, связанном с сервером.
type ServerDomainConfig struct {
	// Домен сервера. Может иметь в качестве значения также и IP адрес сервера с портом, если домен не настроен.
	Summary string `json:"summary"`

	// Содержимое CNAME записи для TCPShield. Имеет значение только, если пользователь настроил домен от хостинга
	// и включил поддержку TCPShield.
	TcpShieldRecord string `json:"TCPShieldRecord"`
}

// IsFrozenByUser возвращает true, если в данный момент сервер заморожен пользователем.
func (s *Server) IsFrozenByUser() bool {
	return s.FrozenAt.Valid && s.Cost.Freeze.Valid
}

// IsTemporary возвращает true, если срок действия данного сервера ограничен.
func (s *Server) IsTemporary() bool {
	return s.ExpiresAt.Valid
}

// GetExternalServer получает данные о внешнем сервере, соответствующем данному (внутреннему).
func (s *Server) GetExternalServer(client *Client) (*ExternalServer, error) {
	return client.GetExternalServer(s.ID)
}

// GetPricing получает актуальную информацию о стоимости сервера.
func (s *Server) GetPricing(client *Client) (*ServicePricing, error) {
	return client.GetServerPricing(s.ID)
}

// Block блокирует сервер. Если сервер заморожен пользователем, заморозка снимается, и только после этого сервер
// блокируется. Вернёт ошибку 409, если сервер уже заблокирован.
func (s *Server) Block(client *Client) error {
	return client.BlockServer(s.ID)
}

// Unblock разблокирует сервер. Вернёт ошибку 409, если сервер не заблокирован.
func (s *Server) Unblock(client *Client) error {
	return client.UnblockServer(s.ID)
}

// GetServers получает список всех серверов, доступных в системе.
func (c *Client) GetServers() (*[]Server, error) {
	return InvokeEndpoint[[]Server](c, http.MethodGet, "/servers", nil)
}

// GetServer получает информацию о сервере с данным идентификатором.
func (c *Client) GetServer(id int64) (*Server, error) {
	return InvokeEndpoint[Server](c, http.MethodGet, fmt.Sprintf("/servers/%d", id), nil)
}

// BlockServer блокирует сервер с заданным идентификатором. Если сервер заморожен пользователем, заморозка снимается,
// и только после этого сервер блокируется. Вернёт ошибку 409, если сервер уже заблокирован.
func (c *Client) BlockServer(id int64) error {
	return InvokeVoidEndpoint(c, http.MethodPost, fmt.Sprintf("/servers/%d/blocking", id), nil)
}

// UnblockServer разблокирует сервер с заданным идентификатором. Вернёт ошибку 409, если сервер не заблокирован.
func (c *Client) UnblockServer(id int64) error {
	return InvokeVoidEndpoint(c, http.MethodDelete, fmt.Sprintf("/servers/%d/blocking", id), nil)
}

// ExternalServer - информация о сервере во внешней системе. Сейчас берётся только из панели Pterodactyl.
type ExternalServer struct {
	// ControlURL - адрес страницы с панелью управления сервером.
	// Например, https://panel.superhub.host/servers/1a2b3c4d
	ControlURL string `json:"controlUrl"`

	// Name - название сервера.
	Name string `json:"name"`

	// Identifier - идентификатор сервера.
	Identifier string `json:"identifier"`

	// ResourceLimits содержит информацию об ограничениях по ресурсам, доступным серверу.
	ResourceLimits Resources `json:"resourceLimits"`

	// FeatureLimits содержит информацию об ограничениях дополнительных возможностей сервера.
	FeatureLimits FeatureLimits `json:"featureLimits"`

	// IsSuspended имеет значение true, если сервер заморожен во внешней системе.
	IsSuspended bool `json:"suspended"`

	// NodeID - идентификатор узла во внешней системе, на котором располагается сервер.
	NodeID int64 `json:"nodeId"`

	// NestID - идентификатор nest в Pterodactyl, используемого на сервере.
	NestID int64 `json:"nestId"`

	// EggID - идентификатор egg в Pterodactyl, используемого на сервере.
	EggID int64 `json:"eggId"`
}

// Resources содержит информацию об основных ресурсах сервера - ЦПУ, ОЗУ, диске.
type Resources struct {
	// CPU - количество ядер ЦПУ.
	CPU float64 `json:"cpu"`

	// Memory - оперативная память в ГБ.
	Memory float64 `json:"memory"`

	// Disk - размер дискового пространства в ГБ.
	Disk float64 `json:"disk"`
}

// FeatureLimits содержит информацию о доступных дополнительных возможностях сервера.
type FeatureLimits struct {
	// Databases - количество баз данных, которые может создать пользователь.
	Databases int64 `json:"databases"`

	// Backups - количество резервных копий, которые может создать пользователь.
	Backups int64 `json:"backups"`
}

// GetExternalServer получает данные о внешнем сервере, соответствующем внутреннему с заданным идентификатором internalID.
func (c *Client) GetExternalServer(internalID int64) (*ExternalServer, error) {
	return InvokeEndpoint[ExternalServer](c, http.MethodGet, fmt.Sprintf("/servers/%d/external", internalID), nil)
}

// ServicePricing - структура, содержащая информацию о текущей стоимости конкретной услуги.
type ServicePricing struct {
	// Текущая стоимость услуги в рублях.
	ActualCost float64 `json:"actualCost"`

	// Тип политики ценообразования, которая используется в данный момент для данной услуги.
	PricingPolicyType PricingPolicyType `json:"pricingPolicyType"`
}

// GetServerPricing получает актуальную информацию о стоимости сервера.
func (c *Client) GetServerPricing(serverID int64) (*ServicePricing, error) {
	return InvokeEndpoint[ServicePricing](c, http.MethodGet, fmt.Sprintf("/servers/%d/pricing", serverID), nil)
}
