package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID    uint
	NotifType string `json:"notifType" gorm:"non null"`
	Message   string `json:"message" gorm:"non null"`
}

func (Notification) TableName() {
}
