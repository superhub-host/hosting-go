package superhub

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/guregu/null"
)

const CurrentUserReference = "@self"

// User - пользователь, зарегистрированный на сайте хостинга. Поля отражают основные параметры, используемые в ЛК.
type User struct {
	// Идентификатор пользователя в системе.
	// Не имеет ничего общего с идентификатором в панели Pterodactyl.
	ID int64 `json:"id"`

	// Адрес электронной почты, используемый пользователем для авторизации.
	// Не может повторяться у разных пользователей.
	Email string `json:"email"`

	// Никнейм пользователя.
	// Не может повторяться у разных пользователей.
	Name string `json:"name"`

	// Текущий баланс пользователя в рублях.
	Balance float64 `json:"balance"`

	// Discord содержит информацию о привязанном аккаунте Discord.
	// API возвращает это поле даже если привязанного аккаунта нет.
	Discord LinkedDiscord `json:"discord"`

	// VK содержит информацию о привязанном аккаунте ВК.
	// API возвращает это поле даже если привязанного аккаунта нет.
	VK LinkedVK `json:"vk"`

	// Referral содержит информацию, связанную с реферальной системой.
	// Включает как информацию о реферале текущего пользователя,
	// так и о параметрах для реферальной системы самого пользователя.
	Referral Referral `json:"referral"`

	// HasMfaEnabled имеет значение true, если на аккаунте пользователя
	// включена двухфакторная аутентификация.
	HasMfaEnabled bool `json:"hasMfaEnabled"`

	// HadTestServer имеет значение true, если пользователь когда-либо имел тестовый сервер.
	// Соответственно, значение false отражает то, что в данный момент пользователь может
	// заказать новый тестовый сервер.
	HadTestServer bool `json:"hadTestServer"`

	// Дата регистрации пользователя.
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt null.Time `json:"updatedAt"`
}

// LinkedDiscord содержит информацию о привязанном аккаунте пользователя в Discord.
type LinkedDiscord struct {
	// Идентификатор пользователя в Discord.
	ID null.String `json:"id"`

	// Имеет значение true, если пользователь получил бонус за привязку аккаунта.
	AcquiredLinkBonus bool `json:"linkBonus"`
}

// LinkedVK содержит информацию о привязанном аккаунте пользователя в ВК.
type LinkedVK struct {
	// Идентификатор пользователя в ВК.
	ID null.Int `json:"id"`

	// Имеет значение true, если пользователь получил бонус за привязку аккаунта.
	AcquiredLinkBonus bool `json:"linkBonus"`

	// Имеет значение true, если пользователь получил бонус за отзыв в обсуждении.
	AcquiredFeedbackBonus bool `json:"feedbackBonus"`
}

// Referral содержит информацию, связанную с реферальной системой.
// Включает как информацию о реферале текущего пользователя,
// так и о параметрах для реферальной системы самого пользователя.
type Referral struct {
	// Имеет значение true, если пользователь получил бонус за первое пополнение баланса.
	AcquiredBonus bool `json:"bonus"`

	// Реферальный код пользователя.
	Code string `json:"code"`

	// Идентификатор реферала - владельца реферального кода, который текущий пользователь использовал при регистрации.
	UserID null.Int `json:"userId"`
}

// HasLinkedDiscord возвращает true, если пользователь привязал аккаунт в Discord.
func (u *User) HasLinkedDiscord() bool {
	return u.Discord.ID.Valid
}

// HasLinkedVK возвращает true, если пользователь привязал аккаунт в ВК.
func (u *User) HasLinkedVK() bool {
	return u.VK.ID.Valid
}

// HasReferral возвращает true, если пользователь зарегистрировался по приглашению от другого пользователя.
func (u *User) HasReferral() bool {
	return u.Referral.UserID.Valid
}

// GetOwnedServers получает список серверов, владельцем которых является данный пользователь.
// Если передан параметр external = true, в структуре полученных серверов будет доступно поле ExternalServer, если для
// конкретного сервера доступен внешний сервер.
func (u *User) GetOwnedServers(client *Client, external bool) (*[]Server, error) {
	return client.GetOwnedServers(u.ID, external)
}

func (c *Client) getUser(id string) (*User, error) {
	return InvokeEndpoint[User](c, http.MethodGet, fmt.Sprintf("/users/%s", id), nil)
}

// GetUser получает пользователя по указанному числовому идентификатору.
// Чтобы получить информацию о текущем пользователе, используйте GetCurrentUser.
func (c *Client) GetUser(id int64) (*User, error) {
	return c.getUser(strconv.FormatInt(id, 10))
}

// GetCurrentUser получает информацию о владельце учётных данных, с помощью которых производится авторизация.
func (c *Client) GetCurrentUser() (*User, error) {
	return c.getUser(CurrentUserReference)
}

// GetOwnedServers получает список серверов, владельцем которых является пользователь с заданным идентификатором ownerID.
// Если передан параметр external = true, в структуре полученных серверов будет доступно поле ExternalServer, если для
// конкретного сервера доступен внешний сервер.
func (c *Client) GetOwnedServers(ownerID int64, external bool) (*[]Server, error) {
	return InvokeEndpoint[[]Server](c, http.MethodGet, fmt.Sprintf("/users/%d/servers?external=%t", ownerID, external), nil)
}
