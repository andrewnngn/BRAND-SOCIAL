package models

import (
	"time"

	"gorm.io/gorm"
)

type Posts []Post

type Post struct {
	gorm.Model
	TextContent string    `json:"text"`
	ImageURL    string    `json:"image"`
	VideoURL    string    `json:"video"`
	Comments    []Comment `json:"comments" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID      uint      `json:"user_id"`
	LikedBy     []User    `json:"likedby" gorm:"many2many:posts_liked_by;"`
	User        User      `json:"author" gorm:"foreignKey:UserID"`
	Likes       uint      `json:"likes"`
}

func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Posts) Less(i, j int) bool { return p[i].UpdatedAt.After(p[j].UpdatedAt) }

func (p Post) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

func (p Post) GetLikes() uint {
	return p.Likes
}

func (Post) TableName() {
}
