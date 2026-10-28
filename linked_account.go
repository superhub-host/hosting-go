package superhub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// LinkedAccount — аккаунт пользователя на внешней платформе, привязанный к его учётной записи на хостинге.
// Привязанные аккаунты используются для доступа к возможностям платформ (см. PlatformFeature).
type LinkedAccount struct {
	// Идентификатор привязанного аккаунта в системе.
	ID string `json:"id"`

	// Идентификатор платформы, на которой зарегистрирован привязанный аккаунт.
	PlatformID string `json:"platformId"`

	// Идентификатор аккаунта на внешней платформе.
	ExternalID string `json:"externalId"`

	// Имя аккаунта на внешней платформе, например, логин.
	Name string `json:"name"`

	// Отображаемое имя аккаунта на внешней платформе.
	DisplayName string `json:"displayName"`

	// Дата привязки аккаунта.
	LinkedAt time.Time `json:"linkedAt"`
}

// LinkActionType — тип действия, которое пользователь должен выполнить, чтобы привязать аккаунт.
type LinkActionType string

const (
	// LinkActionTypeFollowLink — пользователю необходимо перейти по ссылке, указанной в LinkAction.LinkURL.
	LinkActionTypeFollowLink LinkActionType = "FOLLOW_LINK"

	// LinkActionTypeSendCode — пользователю необходимо отправить код, указанный в LinkAction.Code,
	// на внешнюю платформу.
	LinkActionTypeSendCode LinkActionType = "SEND_CODE"

	// LinkActionTypeProvideCode — пользователю необходимо ввести код, полученный на внешней платформе.
	// Дополнительные параметры для такого действия не передаются.
	LinkActionTypeProvideCode LinkActionType = "PROVIDE_CODE"
)

// LinkAction — параметры действия, которое пользователь должен выполнить для привязки аккаунта. Набор заполненных
// полей зависит от значения LinkInstructions.ActionType.
type LinkAction struct {
	// Ссылка, по которой необходимо перейти пользователю.
	// Заполняется только при LinkActionTypeFollowLink.
	LinkURL string `json:"linkUrl,omitempty"`

	// Код, который пользователь должен отправить на внешнюю платформу.
	// Заполняется только при LinkActionTypeSendCode.
	Code string `json:"code,omitempty"`
}

// LinkInstructions — инструкция по привязке аккаунта внешней платформы, которую необходимо показать пользователю.
type LinkInstructions struct {
	// Заголовок инструкции.
	Title string `json:"title"`

	// Описание инструкции.
	Description string `json:"description"`

	// Тип действия, которое необходимо выполнить пользователю. Возможные значения смотри в LinkActionType.
	ActionType LinkActionType `json:"actionType"`

	// Параметры действия, которое необходимо выполнить пользователю.
	Action LinkAction `json:"action"`
}

// LinkAccountForm — параметры запроса инструкции по привязке аккаунта.
type LinkAccountForm struct {
	// Идентификатор платформы, аккаунт на которой требуется привязать.
	PlatformID string `json:"platformId"`
}

// PlatformFeature — возможность внешней платформы, доступная пользователю благодаря привязанному аккаунту.
type PlatformFeature struct {
	// Идентификатор возможности.
	ID string `json:"id"`

	// Отображаемое название возможности.
	DisplayName string `json:"displayName"`

	// Описание возможности.
	Description string `json:"description"`

	// Текст, который следует показать на кнопке, запускающей использование возможности.
	ActionText string `json:"actionText"`

	// OnlyOnce имеет значение true, если возможность может быть использована только один раз.
	OnlyOnce bool `json:"onlyOnce"`
}

// PlatformFeatureStatus — возможность внешней платформы вместе с информацией о её доступности
// для конкретного привязанного аккаунта.
type PlatformFeatureStatus struct {
	// Описание возможности.
	Feature PlatformFeature `json:"feature"`

	// CanUse имеет значение true, если возможность может быть использована прямо сейчас.
	CanUse bool `json:"canUse"`

	// IsUsed имеет значение true, если возможность уже была использована.
	IsUsed bool `json:"isUsed"`
}

// UseFeatureResult — результат использования возможности внешней платформы.
type UseFeatureResult struct {
	// IsApproved имеет значение true, если возможность была успешно использована.
	IsApproved bool `json:"isApproved"`

	// Сообщение, описывающее результат использования возможности.
	Message string `json:"message"`
}

// GetLinkedAccounts получает список аккаунтов внешних платформ, привязанных к данному пользователю.
func (u *User) GetLinkedAccounts(client *Client) (*[]LinkedAccount, error) {
	return client.GetLinkedAccounts(u.ID)
}

// GetLinkedAccount получает привязанный аккаунт данного пользователя по его идентификатору.
func (u *User) GetLinkedAccount(client *Client, linkedAccountID string) (*LinkedAccount, error) {
	return client.GetLinkedAccount(u.ID, linkedAccountID)
}

// LinkAccount получает инструкцию по привязке аккаунта внешней платформы к данному пользователю.
func (u *User) LinkAccount(client *Client, form LinkAccountForm) (*LinkInstructions, error) {
	return client.LinkAccount(u.ID, form)
}

// GetPlatformFeatures получает список возможностей внешней платформы, доступных данному пользователю
// благодаря привязанному аккаунту.
func (u *User) GetPlatformFeatures(client *Client, linkedAccountID string) (*[]PlatformFeatureStatus, error) {
	return client.GetPlatformFeatures(u.ID, linkedAccountID)
}

// UsePlatformFeature использует возможность внешней платформы от имени привязанного аккаунта данного пользователя.
func (u *User) UsePlatformFeature(client *Client, linkedAccountID, featureID string) (*UseFeatureResult, error) {
	return client.UsePlatformFeature(u.ID, linkedAccountID, featureID)
}

// GetPlatformFeatures получает список возможностей внешней платформы, доступных благодаря данному
// привязанному аккаунту. Идентификатор пользователя userID должен принадлежать владельцу аккаунта.
func (a *LinkedAccount) GetPlatformFeatures(client *Client, userID uuid.UUID) (*[]PlatformFeatureStatus, error) {
	return client.GetPlatformFeatures(userID, a.ID)
}

// UsePlatformFeature использует возможность внешней платформы от имени данного привязанного аккаунта.
// Идентификатор пользователя userID должен принадлежать владельцу аккаунта.
func (a *LinkedAccount) UsePlatformFeature(client *Client, userID uuid.UUID, featureID string) (*UseFeatureResult, error) {
	return client.UsePlatformFeature(userID, a.ID, featureID)
}

// GetLinkedAccounts получает список аккаунтов внешних платформ, привязанных к пользователю
// с заданным идентификатором userID.
func (c *Client) GetLinkedAccounts(userID uuid.UUID) (*[]LinkedAccount, error) {
	return InvokeEndpoint[[]LinkedAccount](c, http.MethodGet, fmt.Sprintf("/users/%s/linked-accounts", userID), nil, nil)
}

// GetLinkedAccount получает привязанный аккаунт пользователя по его идентификатору. Вернёт ошибку 404, если
// аккаунт не существует или привязан к другому пользователю.
func (c *Client) GetLinkedAccount(userID uuid.UUID, linkedAccountID string) (*LinkedAccount, error) {
	return InvokeEndpoint[LinkedAccount](c, http.MethodGet, fmt.Sprintf("/users/%s/linked-accounts/%s", userID, linkedAccountID), nil, nil)
}

// LinkAccount получает инструкцию по привязке аккаунта внешней платформы к пользователю с заданным
// идентификатором userID. Сама привязка производится внешней платформой после того, как пользователь
// выполнит описанное в инструкции действие.
func (c *Client) LinkAccount(userID uuid.UUID, form LinkAccountForm) (*LinkInstructions, error) {
	return InvokeEndpoint[LinkInstructions](c, http.MethodPost, fmt.Sprintf("/users/%s/linked-accounts", userID), nil, form)
}

// GetPlatformFeatures получает список возможностей внешней платформы, доступных благодаря привязанному аккаунту.
// Вернёт ошибку 404, если аккаунт не существует или привязан к другому пользователю.
func (c *Client) GetPlatformFeatures(userID uuid.UUID, linkedAccountID string) (*[]PlatformFeatureStatus, error) {
	return InvokeEndpoint[[]PlatformFeatureStatus](c, http.MethodGet, fmt.Sprintf("/users/%s/linked-accounts/%s/platform-features", userID, linkedAccountID), nil, nil)
}

// UsePlatformFeature использует возможность внешней платформы от имени привязанного аккаунта. Вернёт ошибку 404,
// если аккаунт не существует или привязан к другому пользователю, и ошибку 409, если возможность не может быть
// использована — в этом случае причина будет указана в сообщении ошибки.
func (c *Client) UsePlatformFeature(userID uuid.UUID, linkedAccountID, featureID string) (*UseFeatureResult, error) {
	return InvokeEndpoint[UseFeatureResult](c, http.MethodPost, fmt.Sprintf("/users/%s/linked-accounts/%s/platform-features/%s", userID, linkedAccountID, featureID), nil, nil)
}
