package superhub

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// Node — узел хостинга — физический сервер, на котором размещаются услуги, покупаемые пользователями. Поля включают
// в себя общую информацию о данном узле, которая используется при оформлении услуги.
type Node struct {
	// Идентификатор ноды.
	ID uuid.UUID `json:"id"`

	// Название ноды.
	Name string `json:"name"`

	// Информация о физическом расположении ноды.
	Location NodeLocation `json:"location"`

	// Информация о подключении к ноде.
	Connectivity NodeConnectivity `json:"connectivity"`

	// Информация о защите ноды от DDoS атак.
	Protection NodeProtection `json:"protection"`

	// Конфигурация сервера.
	Build NodeBuild `json:"build"`

	// Идентификатор линейки тарифов, используемой для серверов на этой ноде.
	TariffSetId string `json:"tariffSetId"`

	// Уровень доступности в месяц в процентах.
	SLA string `json:"sla"`

	// Нагрузка — число от 0 до 1, показывающее загруженность ноды.
	// 0 — нет нагрузки, 1 — максимальная нагрузка.
	Load float64 `json:"load"`

	// Доступна ли нода для самостоятельного размещения серверов пользователями?
	IsPublic bool `json:"isPublic"`

	// Доступна ли нода для размещения серверов в принципе?
	IsAvailable bool `json:"isAvailable"`

	tariffGroupId uuid.UUID
}

// NodeLocation — информация о физическом расположении ноды.
type NodeLocation struct {
	// Название локации. Например, MSK-1.
	Name string `json:"name"`

	// Город, в котором располагается дата-центр.
	City string `json:"city"`

	// Страна, в которой располагается дата-центр.
	CountryCode string `json:"countryCode"`
}

type NodeConnectivity struct {
	Inet4Address *string `json:"inet4Address"`
	Inet6Address *string `json:"inet6Address"`
}

// ProtectionLevel показывает уровень защиты расположения от DDoS атак.
type ProtectionLevel string

const (
	// ProtectionLevelNone — отсутствие защиты от DDoS атак.
	ProtectionLevelNone ProtectionLevel = "None"

	// ProtectionLevelBasic — базовый уровень защиты от DDoS атак.
	ProtectionLevelBasic ProtectionLevel = "Basic"

	// ProtectionLevelFull — максимальный уровень защиты от DDoS атак.
	ProtectionLevelFull ProtectionLevel = "Full"
)

// ProtectionKind показывает, каким образом обеспечивается защита расположения от DDoS атак.
type ProtectionKind string

const (
	// ProtectionKindISP используется на расположениях, где защита от DDoS атак обеспечивается исключительно средствами
	// интернет-провайдера, либо не обеспечивается вовсе.
	ProtectionKindISP = "ISP"

	// ProtectionKindExternal используется на расположениях, где клиенты обязуются использовать внешнюю защиту от DDoS
	// атак, предоставляемую либо хостингом, либо сторонним провайдером.
	ProtectionKindExternal = "External"

	// ProtectionKindInternal используется на расположениях, где применяется встроенная защита от DDoS атак.
	ProtectionKindInternal = "Internal"
)

// NodeProtection содержит информацию об уровне защищённости расположения от DDoS атак.
type NodeProtection struct {
	// Вид защиты от DDoS атак на расположении.
	Kind ProtectionKind

	// Уровень защищённости расположения от DDoS атак.
	Level ProtectionLevel
}

// AddressPair содержит IPv4 и IPv6 адреса, указывающие на один узел.
type AddressPair struct {
	V4 null.String `json:"v4"`
	V6 null.String `json:"v6"`
}

// NodeBuild содержит информацию о комплектующих сервера, на котором запускаются пользовательские услуги.
type NodeBuild struct {
	// Информация о центральном процессоре сервера.
	CPU CPU
}

// CPU описывает модель центрального процессора.
type CPU struct {
	// Производитель процессора.
	Vendor string

	// Модель процессора.
	Model string

	// Количество потоков.
	Threads int

	// Тактовая частота в ГГц.
	Frequency float64

	// Список результатов бенчмарков для процессора.
	Benchmarks []Benchmark
}

// Benchmark содержит информацию о результатах выполнения теста производительности ЦПУ.
type Benchmark struct {
	// Название теста.
	Name string

	// Итоговый счёт в тесте.
	Value float64
}

// NodeLoad является обёрткой для значения нагрузки ноды.
// Используется только при сериализации и десериализации запросов и ответов.
type NodeLoad struct {
	// Нагрузка — число от 0 до 1, показывающее загруженность ноды.
	// 0 — нет нагрузки, 1 — максимальная нагрузка.
	Load float64 `json:"load"`
}

// UpdateLoad обновляет информацию о загруженности ноды.
func (n *Node) UpdateLoad(client *Client, load *NodeLoad) (*NodeLoad, error) {
	return client.UpdateNodeLoad(n.tariffGroupId, n.ID, load)
}

// GetNode получает информацию о ноде с заданным идентификатором.
func (c *Client) GetNode(tariffGroupId uuid.UUID, nodeId uuid.UUID) (*Node, error) {
	n, err := InvokeEndpoint[Node](c, http.MethodGet, fmt.Sprintf("/tariff-groups/%s/nodes/%s", tariffGroupId, nodeId), nil, nil)
	if err != nil {
		return nil, err
	}

	n.tariffGroupId = tariffGroupId
	return n, nil
}

// GetNodes получает список всех доступных нод.
func (c *Client) GetNodes(tariffGroupId uuid.UUID) (*[]Node, error) {
	nodes, err := InvokeEndpoint[[]Node](c, http.MethodGet, fmt.Sprintf("/tariff-groups/%s/nodes", tariffGroupId), nil, nil)
	if err != nil {
		return nil, err
	}

	for _, n := range *nodes {
		n.tariffGroupId = tariffGroupId
	}

	return nodes, nil
}

// UpdateNodeLoad обновляет информацию о загруженности ноды.
func (c *Client) UpdateNodeLoad(tariffGroupId uuid.UUID, nodeId uuid.UUID, load *NodeLoad) (*NodeLoad, error) {
	return InvokeEndpoint[NodeLoad](c, http.MethodPut, fmt.Sprintf("/tariff-groups/%s/nodes/%s/load", tariffGroupId, nodeId), nil, load)
}
