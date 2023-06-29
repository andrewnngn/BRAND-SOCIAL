package models

import (
	"time"

	"gorm.io/gorm"
)

type Comments []Comment

type Comment struct {
	gorm.Model
	Text     string `json:"content"`
	ImageURL string `json:"image"`
	VideoURL string `json:"video"`
	PostID   uint
	UserID   uint
	Likes    uint   `json:"likes"`
	LikedBy  []User `json:"likedby" gorm:"many2many:comments_liked_by;"`
	ParentID *uint
	Parent   *Comment  `gorm:"foreignKey:ParentID"`
	Replies  []Comment `gorm:"foreignKey:ParentID"`
}

func (c Comments) Len() int           { return len(c) }
func (c Comments) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Comments) Less(i, j int) bool { return c[i].UpdatedAt.After(c[j].UpdatedAt) }

func (c Comment) GetUpdatedAt() time.Time {
	return c.UpdatedAt
}

func (c Comment) GetLikes() uint {
	return c.Likes
}

func (Comment) TableName() {
}
