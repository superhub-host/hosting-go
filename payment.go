package superhub

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

// Amount — размер платежа. Описывает сумму и валюту, в которой проводится платёж.
type Amount struct {
	// Сумма платежа в валюте, соответствующей значению Currency.
	Amount float64 `json:"amount"`

	// Валюта платежа.
	Currency string `json:"currency"`
}

// PaymentSourceType — тип "источника" платежа — действия, вызвавшего создание данного платежа.
type PaymentSourceType string

const (
	// PaymentSourceDeposit — пополнение баланса.
	PaymentSourceDeposit PaymentSourceType = "DEPOSIT"

	// PaymentSourceInstance — списание за услугу.
	PaymentSourceInstance PaymentSourceType = "INSTANCE"

	// PaymentSourceReferral — проценты от пополнений приглашённых пользователей.
	PaymentSourceReferral PaymentSourceType = "REFERRAL"

	// PaymentSourceDepositBonus — бонус за пополнение баланса.
	PaymentSourceDepositBonus PaymentSourceType = "DEPOSIT_BONUS"

	// PaymentSourcePersonnelPayout — внутренние выплаты персоналу.
	PaymentSourcePersonnelPayout PaymentSourceType = "PERSONNEL_PAYOUT"

	// PaymentSourceRefund — возвраты за услуги на баланс и вывод средств с баланса.
	PaymentSourceRefund PaymentSourceType = "REFUND"

	// PaymentSourcePrizePayout — выплаты призов за участие в событиях хостинга.
	PaymentSourcePrizePayout PaymentSourceType = "PRIZE_PAYOUT"

	// PaymentSourceReferralWelcomeBonus — приветственный бонус для пользователей, зарегистрированных по приглашению.
	PaymentSourceReferralWelcomeBonus PaymentSourceType = "REFERRAL_WELCOME_BONUS"

	// PaymentSourceAccountLinkBonus — бонус за привязку аккаунта.
	PaymentSourceAccountLinkBonus PaymentSourceType = "ACCOUNT_LINK_BONUS"

	// PaymentSourceOther используется как стандартное значение для типа источника платежа.
	PaymentSourceOther PaymentSourceType = "OTHER"
)

// PaymentSource — «источник» платежа. Описывает то, почему был создан данный платёж. Значения, связанные с источником
// платежа, могут быть использованы для фильтрации и группировки платежей.
type PaymentSource struct {
	// Тип источника (см. PaymentSourceType)
	Type PaymentSourceType `json:"type,omitempty"`

	// Идентификатор источника. Для некоторых типов всегда имеет пустое значение ("TOP_UP", "OTHER"),
	// для других — всегда непустое. Например, для типа "REFERRAL" будет содержать значение пользователя, от которого
	// получен бонус по реферальной системе.
	ID *uuid.UUID `json:"id,omitempty"`
}

// PaymentMode отражает режим, в котором система обрабатывает платёж.
type PaymentMode string

const (
	// PaymentModeProduction — стандартный режим, при котором обработка платежа производится полностью.
	PaymentModeProduction = "PRODUCTION"

	// PaymentModeTest — тестовый режим, производится полная обработка платежа без изменения баланса пользователя.
	PaymentModeTest = "TEST"
)

// Payment описывает платёж — сущность, используемую для хранения истории изменения баланса пользователя на хостинге.
// Платежи могут иметь как положительную, так и отрицательную сумму. Платежи с положительной суммой отражают пополнения
// баланса, будь то пополнение пользователем или администрацией хостинга. Платежи с отрицательной суммой отражают
// списания средств с баланса пользователя, например, для оплаты услуг хостинга.
type Payment struct {
	// Идентификатор платежа. В текущей реализации представляет собой последовательность из 16 байт, представленную
	// в шестнадцатеричном виде. Не гарантируется, что все идентификаторы будут в таком формате в будущем.
	ID string `json:"id"`

	// Сумма платежа.
	Amount Amount `json:"amount"`

	// Описание платежа. Может быть произвольной строкой или отсутствовать вообще.
	Description *string `json:"description"`

	// «Источник» платежа — действие, которое вызвало создание данного платежа.
	Source PaymentSource `json:"source"`

	// Режим проведения платежа. В большинстве случаев имеет значение "PRODUCTION", т.е. платёж обрабатывается
	// полностью. В зависимости от данного значения платёж в системе может обрабатываться по-разному.
	Mode PaymentMode `json:"mode"`

	// Дата создания платежа.
	CreatedAt time.Time `json:"createdAt"`
}

// GetPayments получает список всех платежей в системе.
func (c *Client) GetPayments(params *PaginationParams) (*[]Payment, error) {
	query := &url.Values{}
	params.Encode(query)

	return InvokeEndpoint[[]Payment](c, http.MethodGet, "/payments", query, nil)
}

// GetUserPayments получает список платежей пользователя.
func (c *Client) GetUserPayments(userID uuid.UUID, params *PaginationParams) (*[]Payment, error) {
	query := &url.Values{}
	params.Encode(query)

	return InvokeEndpoint[[]Payment](c, http.MethodGet, fmt.Sprintf("/users/%s/payments", userID), query, nil)
}

type PaymentCreationForm struct {
	// Сумма платежа.
	Amount Amount `json:"amount"`

	// Описание платежа.
	Description null.String `json:"description"`

	// Источник платежа.
	Source *PaymentSource `json:"source"`
}

// CreatePayment создаёт платёж для данного пользователя.
func (c *Client) CreatePayment(userID uuid.UUID, form PaymentCreationForm) (*Payment, error) {
	return InvokeEndpoint[Payment](c, http.MethodPost, fmt.Sprintf("/users/%s/payments", userID), nil, form)
}
