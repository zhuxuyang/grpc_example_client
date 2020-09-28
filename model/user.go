package model

import (
	"github.com/jinzhu/gorm"
)

type UserA struct {
	gorm.Model
	Payload string
}

func (m *UserA) TableName() string {
	return "user_a"
}
