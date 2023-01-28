package superhub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/guregu/null"
)

// ServerState - состояние сервера. Значение отражает то, готов ли сервер к работе после установки. При этом, значение
// не зависит от того, доступен ли сервер в данный момент, если установка уже прошла.
type ServerState int

const (
	// ServerStateFailed используется, когда системе не удалось полностью создать сервер.
	ServerStateFailed = -1

	// ServerStateUnknown - стандартное значение состояния, которое означает, что сервер находится в процессе создания,
	// либо состояние сервера неизвестно.
	ServerStateUnknown = 0

	// ServerStateNormal используется, когда сервер готов к работе.
	ServerStateNormal = 1
)

// Server - сервер, принадлежащий определённому пользователю. Содержит только общую информацию, используемую в ЛК,
// такую как стоимость, скидка, дата создания и т.п.
// Более подробная информация, например, о лимитах ресурсов, доступна в структуре внешнего сервера (см. ExternalServer)
type Server struct {
	// ID - идентификатор сервера в системе.
	// Не имеет ничего общего с идентификатором в панели Pterodactyl.
	ID int64 `json:"id"`

	// PteroID - идентификатор соответствующего сервера в панели Pterodactyl.
	PteroID int64 `json:"pteroId"`

	// OwnerID - идентификатор пользователя, являющегося владельцем данного сервера.
	OwnerID int64 `json:"ownerId"`

	// State - текущее состояние сервера (см. ServerState)
	State ServerState `json:"state"`

	// Cost - базовая стоимость сервера в рублях с учётом скидки.
	Cost float64 `json:"cost"`

	// Sale - размер скидки, принимает значения от 0 до 1, где 0 - отсутствие скидки, а 1 - скидка 100%.
	Sale float64 `json:"sale"`

	// FreezeCost - стоимость сервера при заморозке. Имеет значение только если сервер заморожен.
	FreezeCost null.Float `json:"freezeCost"`

	// TariffID - идентификатор тарифа, который используется на сервере в данный момент. Имеет значение только если
	// пользователь выбрал готовый тариф, а не свою конфигурацию.
	TariffID null.String `json:"tariffId"`

	// Domain - домен сервера. Может иметь в качестве значения также и IP адрес сервера с портом,
	// если домен не настроен.
	Domain string `json:"domain"`

	// TCPShieldRecord - содержимое CNAME записи для TCPShield. Имеет значение только если пользователь настроил
	// домен от хостинга и включил поддержку TCPShield.
	TCPShieldRecord null.String `json:"tcpshieldRecord"`

	// ExternalServer - структура, содержащая информацию о внешнем сервере. Имеет значение nil, если конкретная
	// конечная точка, от которой была получена информация о сервере, не подразумевает получение информации о
	// внешнем сервере.
	ExternalServer *ExternalServer `json:"externalServer"`

	// ExpiresAt - дата окончания срока действия сервера. Не имеет ничего общего с датой блокировки сервера за неуплату.
	// Обычно имеет значение только на бесплатных серверах, срок действия которых ограничен двумя днями.
	ExpiresAt null.Time `json:"expiresAt"`

	// FrozenAt - дата заморозки сервера пользователем. Имеет значение только когда пользователь заморозил сервер сам
	// и имеет пустое значение, если сервер заблокирован за неуплату.
	FrozenAt null.Time `json:"frozenAt"`

	// CreatedAt - дата создания сервера.
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt - дата последнего обновления информации о сервере.
	UpdatedAt time.Time `json:"updatedAt"`
}

// IsFrozenByUser возвращает true, если в данный момент сервер заморожен пользователем.
func (s *Server) IsFrozenByUser() bool {
	return s.FrozenAt.Valid && s.FreezeCost.Valid
}

// IsTemporary возвращает true, если срок действия данного сервера ограничен.
func (s *Server) IsTemporary() bool {
	return s.ExpiresAt.Valid
}

// GetExternalServer получает данные о внешнем сервере, соответствующем данному (внутреннему).
func (s *Server) GetExternalServer(client *Client) (*ExternalServer, error) {
	return client.GetExternalServer(s.ID)
}

// GetServers получает список всех серверов, доступных в системе.
func (c *Client) GetServers() (*[]Server, error) {
	return InvokeEndpoint[[]Server](c, http.MethodGet, "/servers", nil)
}

// GetServer получает информацию о сервере с данным идентификатором.
func (c *Client) GetServer(id int64) (*Server, error) {
	return InvokeEndpoint[Server](c, http.MethodGet, fmt.Sprintf("/servers/%d", id), nil)
}

// GetExternalServer получает данные о внешнем сервере, соответствующем внутреннему с заданным идентификатором internalID.
func (c *Client) GetExternalServer(internalID int64) (*ExternalServer, error) {
	return InvokeEndpoint[ExternalServer](c, http.MethodGet, fmt.Sprintf("/servers/%d/external", internalID), nil)
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
	// CPU - процент нагрузки на ЦПУ. 100% = 1 ядро.
	CPU int64 `json:"cpu"`

	// Memory - оперативная память в МБ.
	Memory int64 `json:"memory"`

	// Disk - размер дискового пространства в МБ.
	Disk int64 `json:"disk"`
}

// FeatureLimits содержит информацию о доступных дополнительных возможностях сервера.
type FeatureLimits struct {
	// Databases - количество баз данных, которые может создать пользователь.
	Databases int64 `json:"databases"`

	// Backups - количество резервных копий, которые может создать пользователь.
	Backups int64 `json:"backups"`
}
