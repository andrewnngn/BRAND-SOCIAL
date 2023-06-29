package models

import "time"

type Model interface {
	TableName()
	GetUpdatedAt() time.Time
	GetLikes() uint
}
