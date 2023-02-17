package superhub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/guregu/null"
)

// PaymentAmount - размер платежа. Описывает сумму и валюту, в которой проводится платёж.
type PaymentAmount struct {
	// Сумма платежа в валюте, соответствующей значению Currency.
	Sum float64 `json:"sum"`

	// Валюта платежа.
	Currency string `json:"currency"`
}

// PaymentSourceType - тип "источника" платежа - действия, вызвавшего создание данного платежа.
type PaymentSourceType string

const (
	// PaymentSourceTopUp - пополнение баланса.
	PaymentSourceTopUp PaymentSourceType = "TOP_UP"

	// PaymentSourceServerService - списание за игровой сервер.
	PaymentSourceServerService PaymentSourceType = "SERVER_SERVICE"

	// PaymentSourceReferral - проценты от пополнений приглашённых пользователей.
	PaymentSourceReferral PaymentSourceType = "REFERRAL"

	// PaymentSourceReferralWelcomeBonus - приветственный бонус для пользователей, зарегистрированных по приглашению.
	PaymentSourceReferralWelcomeBonus PaymentSourceType = "REFERRAL_WELCOME_BONUS"

	// PaymentSourceOther используется как стандартное значение для типа источника платежа.
	PaymentSourceOther PaymentSourceType = "OTHER"
)

// PaymentSource - "источник" платежа. Описывает то, почему был создан данный платёж. Значения, связанные с источником
// платежа, могут быть использованы для фильтрации и группировки платежей.
type PaymentSource struct {
	// Тип источника (см. PaymentSourceType)
	Type PaymentSourceType `json:"type,omitempty"`

	// Идентификатор источника. Для некоторых типов всегда имеет пустое значение ("TOP_UP", "OTHER"),
	// для других - всегда непустое. Например, для типа "REFERRAL" будет содержать значение пользователя, от которого
	// получен бонус по реферальной системе.
	ID null.String `json:"id,omitempty"`
}

// PaymentMode отражает режим, в котором система обрабатывает платёж.
type PaymentMode string

const (
	// PaymentModeProduction - стандартный режим, при котором обработка платежа производится полностью.
	PaymentModeProduction = "PRODUCTION"

	// PaymentModeTest - тестовый режим, производится полная обработка платежа без изменения баланса пользователя.
	PaymentModeTest = "TEST"
)

// Payment описывает платёж - сущность, используемую для хранения истории изменения баланса пользователя на хостинге.
// Платежи могут иметь как положительную, так и отрицательную сумму. Платежи с положительной суммой отражают пополнения
// баланса, будь то пополнение пользователем или администрацией хостинга. Платежи с отрицательной суммой отражают
// списания средств с баланса пользователя, например, для оплаты услуг хостинга.
type Payment struct {
	// Идентификатор платежа. В текущей реализации представляет собой последовательность из 16 байт, представленную
	// в шестнадцатеричном виде. При этом не гарантируется, что все идентификаторы будут в таком формате в будущем.
	ID string `json:"id"`

	// Идентификатор пользователя, для которого был проведён платёж. Т.е. изменение баланса, описанное данным платежом,
	// производилось с балансом пользователя, имеющего идентификатор, равный UserID.
	UserID int64 `json:"userId"`

	// Сумма платежа.
	Amount PaymentAmount `json:"amount"`

	// Описание платежа. Может быть произвольной строкой или отсутствовать вообще.
	Description null.String `json:"description"`

	// "Источник" платежа - действие, которое вызвало создание данного платежа.
	Source PaymentSource `json:"source"`

	// Режим проведения платежа. В большинстве случаев имеет значение "PRODUCTION", т.е. платёж обрабатывается
	// полностью. В зависимости от данного значения платёж в системе может обрабатываться по-разному.
	Mode PaymentMode `json:"mode"`

	// Завершён ли платёж? Значение false показывает, что изменение баланса, описываемое платежом, пока не было
	// произведено. Например, при пополнении баланса создаётся платёж, у которого Completed = false, но после того,
	// как пользователь производит оплату, сумма платежа зачисляется на баланс, а Completed изменяется на true.
	Completed bool `json:"completed"`

	// Дата создания платежа.
	CreatedAt time.Time `json:"createdAt"`

	// Дата последнего обновления информации о платеже.
	UpdatedAt null.Time `json:"updatedAt"`
}

// GetPayments получает список всех платежей в системе.
func (c *Client) GetPayments() (*[]Payment, error) {
	return InvokeEndpoint[[]Payment](c, http.MethodGet, "/payments", nil)
}

// GetUserPayments получает список платежей пользователя.
func (c *Client) GetUserPayments(userID int64) (*[]Payment, error) {
	return InvokeEndpoint[[]Payment](c, http.MethodGet, fmt.Sprintf("/users/%d/payments", userID), nil)
}

type PaymentCreationForm struct {
	// Сумма платежа в рублях.
	Amount float64 `json:"amount"`

	// Описание платежа.
	Description null.String `json:"description"`

	// Источник платежа.
	Source *PaymentSource `json:"source"`
}

// CreatePayment создаёт платёж для данного пользователя.
func (c *Client) CreatePayment(userID int64, form PaymentCreationForm) (*Payment, error) {
	return InvokeEndpoint[Payment](c, http.MethodPost, fmt.Sprintf("/users/%d/payments", userID), form)
}
