package api

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type CreateListRequest struct {
	Title string `json:"title"`
}

type CreateTaskRequest struct {
	Text string `json:"text"`
}

type Error struct {
	Code    string  `json:"code"`
	Details *string `json:"details,omitempty"`
	Message string  `json:"message"`
}

type List struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`

	Id openapi_types.UUID `json:"id"`

	Title string `json:"title"`
}

type Task struct {
	Completed *bool     `json:"completed,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	ListId    string    `json:"list_id"`
	Text      string    `json:"text"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateListRequest struct {
	Title string `json:"title"`
}

type UpdateTaskRequest struct {
	Completed *bool   `json:"completed,omitempty"`
	Text      *string `json:"text,omitempty"`
}

type Id = openapi_types.UUID

type Limit = int

type Offset = int

type NotFound = Error

type ServerError = Error

type ValidationError = Error

type ListListsParams struct {
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	Offset *Offset `form:"offset,omitempty" json:"offset,omitempty"`
}

type GetTasksParams struct {
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

type CreateListJSONRequestBody = CreateListRequest

type UpdateListJSONRequestBody = UpdateListRequest

type CreateTaskJSONRequestBody = CreateTaskRequest

type UpdateTaskJSONRequestBody = UpdateTaskRequest
