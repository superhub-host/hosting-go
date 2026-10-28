package superhub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// InstanceComment — комментарий к услуге. Представляет собой произвольную справочную или аудиторскую
// заметку, оставленную сотрудником хостинга или внутренней подсистемой. 
type InstanceComment struct {
	// Идентификатор комментария в системе.
	ID uuid.UUID `json:"id"`

	// Идентификатор услуги, к которой относится комментарий.
	InstanceID uuid.UUID `json:"instanceId"`

	// Происхождение комментария.
	Source string `json:"source"`

	// Текст комментария.
	Text string `json:"text"`

	// Дата создания комментария.
	CreatedAt time.Time `json:"createdAt"`

	// Дата последнего изменения комментария. Будет nil, если комментарий не изменялся.
	UpdatedAt *time.Time `json:"updatedAt"`
}

// InstanceCommentForm — параметры создания и изменения комментария к услуге. При изменении комментария значения
// обоих полей заменяются полностью.
type InstanceCommentForm struct {
	// Происхождение комментария.
	Source string `json:"source"`

	// Текст комментария.
	Text string `json:"text"`
}

// Update изменяет комментарий к услуге.
func (c *InstanceComment) Update(client *Client, form InstanceCommentForm) (*InstanceComment, error) {
	return client.UpdateInstanceComment(c.InstanceID.String(), c.ID.String(), form)
}

// Delete удаляет комментарий к услуге.
func (c *InstanceComment) Delete(client *Client) error {
	return client.DeleteInstanceComment(c.InstanceID.String(), c.ID.String())
}

// GetComments получает комментарии к услуге, отсортированные по дате создания в порядке убывания.
func (i *Instance) GetComments(client *Client) (*[]InstanceComment, error) {
	return client.GetInstanceComments(i.ID.String())
}

// GetComment получает комментарий к услуге по его идентификатору.
func (i *Instance) GetComment(client *Client, commentID string) (*InstanceComment, error) {
	return client.GetInstanceComment(i.ID.String(), commentID)
}

// CreateComment создаёт комментарий к услуге.
func (i *Instance) CreateComment(client *Client, form InstanceCommentForm) (*InstanceComment, error) {
	return client.CreateInstanceComment(i.ID.String(), form)
}

// UpdateComment изменяет комментарий к услуге.
func (i *Instance) UpdateComment(client *Client, commentID string, form InstanceCommentForm) (*InstanceComment, error) {
	return client.UpdateInstanceComment(i.ID.String(), commentID, form)
}

// DeleteComment удаляет комментарий к услуге.
func (i *Instance) DeleteComment(client *Client, commentID string) error {
	return client.DeleteInstanceComment(i.ID.String(), commentID)
}

// GetInstanceComments получает комментарии к услуге, отсортированные по дате создания в порядке убывания.
func (c *Client) GetInstanceComments(instanceID string) (*[]InstanceComment, error) {
	return InvokeEndpoint[[]InstanceComment](c, http.MethodGet, fmt.Sprintf("/instances/%s/comments", instanceID), nil, nil)
}

// GetInstanceComment получает комментарий к услуге по его идентификатору. Вернёт ошибку 404, если комментарий
// не существует или относится к другой услуге.
func (c *Client) GetInstanceComment(instanceID, commentID string) (*InstanceComment, error) {
	return InvokeEndpoint[InstanceComment](c, http.MethodGet, fmt.Sprintf("/instances/%s/comments/%s", instanceID, commentID), nil, nil)
}

// CreateInstanceComment создаёт комментарий к услуге.
func (c *Client) CreateInstanceComment(instanceID string, form InstanceCommentForm) (*InstanceComment, error) {
	return InvokeEndpoint[InstanceComment](c, http.MethodPost, fmt.Sprintf("/instances/%s/comments", instanceID), nil, form)
}

// UpdateInstanceComment изменяет комментарий к услуге, полностью заменяя значения полей. Вернёт ошибку 404, если
// комментарий не существует или относится к другой услуге.
func (c *Client) UpdateInstanceComment(instanceID, commentID string, form InstanceCommentForm) (*InstanceComment, error) {
	return InvokeEndpoint[InstanceComment](c, http.MethodPut, fmt.Sprintf("/instances/%s/comments/%s", instanceID, commentID), nil, form)
}

// DeleteInstanceComment удаляет комментарий к услуге. Вернёт ошибку 404, если комментарий не существует или относится
// к другой услуге.
func (c *Client) DeleteInstanceComment(instanceID, commentID string) error {
	return InvokeVoidEndpoint(c, http.MethodDelete, fmt.Sprintf("/instances/%s/comments/%s", instanceID, commentID), nil, nil)
}
