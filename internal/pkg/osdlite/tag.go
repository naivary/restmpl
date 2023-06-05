package osdlite

import "github.com/google/uuid"

type tag struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func NewTag() *tag {
	return &tag{
		ID: uuid.NewString(),
	}
}

func (t tag) TableName() string {
	return "tags"
}
