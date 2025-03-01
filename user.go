package superhub

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

const CurrentUserReference = "@self"

// VerificationLevel — уровень подтверждения пользователя. Чем выше уровень, тем выше доверие к пользователю со стороны
// систем хостинга.
type VerificationLevel string

const (
	VerificationLevelNone   VerificationLevel = "NONE"
	VerificationLevelLow    VerificationLevel = "LOW"
	VerificationLevelMedium VerificationLevel = "MEDIUM"
	VerificationLevelHigh   VerificationLevel = "HIGH"
)

// User — пользователь, зарегистрированный на сайте хостинга. Поля отражают основные параметры, используемые в ЛК.
type User struct {
	// Идентификатор пользователя в системе.
	ID uuid.UUID `json:"id"`

	// Адрес электронной почты, используемый пользователем для авторизации.
	// Не может повторяться у разных пользователей.
	Email string `json:"email"`

	// Имя пользователя.
	// Не может повторяться у разных пользователей.
	Name string `json:"name"`

	// Текущий баланс пользователя.
	Balance Amount `json:"balance"`

	// Уровень подтверждения пользователя.
	VerificationLevel VerificationLevel `json:"verificationLevel"`

	// Referral содержит информацию, связанную с реферальной системой.
	// Включает как информацию о реферале текущего пользователя,
	// так и о параметрах для реферальной системы самого пользователя.
	Referral Referral `json:"referral"`

	// HasMfaEnabled имеет значение true, если на аккаунте пользователя
	// включена двухфакторная аутентификация.
	HasMfaEnabled bool `json:"hasMfaEnabled"`

	// Дата регистрации пользователя.
	CreatedAt time.Time `json:"createdAt"`
}

// Referral содержит информацию, связанную с реферальной системой.
// Включает как информацию о реферале текущего пользователя,
// так и о параметрах для реферальной системы самого пользователя.
type Referral struct {
	// Реферальный код пользователя.
	Code string `json:"code"`

	// Идентификатор реферала — владельца реферального кода, который текущий пользователь использовал при регистрации.
	UserId *uuid.UUID `json:"userId"`
}

// GetOwnedInstances получает список услуг, владельцем которых является данный пользователь.
func (u *User) GetOwnedInstances(client *Client) (*[]Instance, error) {
	return client.GetOwnedInstances(u.ID)
}

// GetPayments получает список платежей, связанных с данным пользователем.
func (u *User) GetPayments(client *Client, params *PaginationParams) (*[]Payment, error) {
	return client.GetUserPayments(u.ID, params)
}

// CreatePayment создаёт платёж для данного пользователя.
func (u *User) CreatePayment(client *Client, form PaymentCreationForm) (*Payment, error) {
	return client.CreatePayment(u.ID, form)
}

func (c *Client) getUser(id string) (*User, error) {
	return InvokeEndpoint[User](c, http.MethodGet, fmt.Sprintf("/users/%s", id), nil, nil)
}

// GetUser получает пользователя по указанному числовому идентификатору.
// Чтобы получить информацию о текущем пользователе, используйте GetCurrentUser.
func (c *Client) GetUser(userId uuid.UUID) (*User, error) {
	return c.getUser(userId.String())
}

func (c *Client) GetUsers(params *PaginationSearchSortParams) (*[]User, error) {
	query := &url.Values{}
	params.Encode(query)

	return InvokeEndpoint[[]User](c, http.MethodGet, fmt.Sprintf("/users"), query, nil)
}

// GetCurrentUser получает информацию о владельце учётных данных, с помощью которых производится авторизация.
func (c *Client) GetCurrentUser() (*User, error) {
	return c.getUser(CurrentUserReference)
}

// GetOwnedInstances получает список услуг, владельцем которых является пользователь с заданным идентификатором ownerId.
func (c *Client) GetOwnedInstances(ownerId uuid.UUID) (*[]Instance, error) {
	return InvokeEndpoint[[]Instance](c, http.MethodGet, fmt.Sprintf("/users/%s/instances", ownerId), nil, nil)
}
