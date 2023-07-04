package superhub

import (
	"fmt"

	"gopkg.in/guregu/null.v4"
	"net/http"
)

type NodeComponent string

const (
	NodeComponentCPU NodeComponent = "cpu"
)

// Node - узел хостинга - физический сервер, на котором размещаются сервера, покупаемые пользователями. Поля включают
// в себя общую информацию о данном узле, которая используется при оформлении сервера. Эта информация не дублирует
// аналогичную ей, доступную в панели Pterodactyl, кроме идентификатора, который всегда совпадает с панелью.
type Node struct {
	// Идентификатор ноды. Совпадает с идентификатором ноды в панели Pterodactyl.
	ID int64 `json:"id"`

	// Название ноды.
	Name string `json:"name"`

	// Имя хоста, которое разрешается на IP адрес ноды через DNS.
	Hostname string `json:"hostname"`

	// Множитель стоимости сервера на этой ноде.
	Multiplier float64 `json:"multiplier"`

	// Название набора цен, используемого для серверов на этой ноде.
	PriceSetName string `json:"priceSetName"`

	// Название линейки тарифов, используемой для серверов на этой ноде.
	TariffSetName string `json:"tariffSetName"`

	// Компоненты ноды - информация об установленных комплектующих.
	// На текущий момент предоставляется только информация о модели процессора (NodeComponentCPU).
	Components map[NodeComponent]string `json:"components"`

	// Лимиты по ресурсам, доступным пользователям для покупки.
	Limits Resources `json:"limits"`

	// Нагрузка - число от 0 до 1, показывающее загруженность ноды.
	// 0 - нет нагрузки, 1 - максимальная нагрузка.
	Load float64 `json:"load"`

	// Информация о физическом расположении ноды.
	Location NodeLocation `json:"location"`

	// Внешний адрес, по которому доступна нода.
	PublicAddress AddressPair `json:"publicAddress"`

	// Скрыта ли нода от пользователей? Если true, то нода не будет отображена при покупке сервера.
	Hidden bool `json:"hidden"`
}

// NodeLocation - информация о физическом расположении ноды.
type NodeLocation struct {
	// Страна, в которой располагается дата-центр.
	Country string `json:"country"`

	// Город, в котором располагается дата-центр.
	City string `json:"city"`

	// Код локации. Например, MSK-1.
	Code string `json:"code"`
}

// AddressPair содержит IPv4 и IPv6 адреса, указывающие на один узел.
type AddressPair struct {
	V4 null.String `json:"v4"`
	V6 null.String `json:"v6"`
}

// NodeLoad является обёрткой для значения нагрузки ноды.
// Используется только при сериализации и десериализации запросов и ответов.
type NodeLoad struct {
	// Нагрузка - число от 0 до 1, показывающее загруженность ноды.
	// 0 - нет нагрузки, 1 - максимальная нагрузка.
	Load float64 `json:"load"`
}

// GetLimits получает лимиты по ресурсам, доступным пользователям при покупке сервера на данной ноде.
func (n *Node) GetLimits(client *Client) (*Resources, error) {
	return client.GetNodeLimits(n.ID)
}

// UpdateLimits изменяет лимиты по ресурсам, доступным пользователям при покупке сервера на данной ноде.
func (n *Node) UpdateLimits(client *Client, limits *Resources) (*Resources, error) {
	return client.UpdateNodeLimits(n.ID, limits)
}

// GetLoad получает текущую нагрузку на ноду.
func (n *Node) GetLoad(client *Client) (*NodeLoad, error) {
	return client.GetNodeLoad(n.ID)
}

// UpdateLoad обновляет информацию о загруженности ноды.
func (n *Node) UpdateLoad(client *Client, load *NodeLoad) (*NodeLoad, error) {
	return client.UpdateNodeLoad(n.ID, load)
}

// GetNode получает информацию о ноде с заданным идентификатором.
func (c *Client) GetNode(id int64) (*Node, error) {
	return InvokeEndpoint[Node](c, http.MethodGet, fmt.Sprintf("/nodes/%d", id), nil)
}

// GetNodes получает список всех доступных нод.
func (c *Client) GetNodes() (*[]Node, error) {
	return InvokeEndpoint[[]Node](c, http.MethodGet, "/nodes", nil)
}

// GetNodeLimits получает лимиты по ресурсам, доступным пользователям при покупке сервера на данной ноде.
func (c *Client) GetNodeLimits(id int64) (*Resources, error) {
	return InvokeEndpoint[Resources](c, http.MethodGet, fmt.Sprintf("/nodes/%d/limits", id), nil)
}

// UpdateNodeLimits изменяет лимиты по ресурсам, доступным пользователям при покупке сервера на данной ноде.
func (c *Client) UpdateNodeLimits(id int64, limits *Resources) (*Resources, error) {
	return InvokeEndpoint[Resources](c, http.MethodPut, fmt.Sprintf("/nodes/%d/limits", id), limits)
}

// GetNodeLoad получает текущую нагрузку на ноду.
func (c *Client) GetNodeLoad(id int64) (*NodeLoad, error) {
	return InvokeEndpoint[NodeLoad](c, http.MethodGet, fmt.Sprintf("/nodes/%d/load", id), nil)
}

// UpdateNodeLoad обновляет информацию о загруженности ноды.
func (c *Client) UpdateNodeLoad(id int64, load *NodeLoad) (*NodeLoad, error) {
	return InvokeEndpoint[NodeLoad](c, http.MethodPut, fmt.Sprintf("/nodes/%d/load", id), load)
}
