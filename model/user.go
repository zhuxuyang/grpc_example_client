package model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type UserA struct {
	gorm.Model
	Payload json.RawMessage
	Ok      bool
}

func (m *UserA) TableName() string {
	return "user_a"
}
