package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name          string         `json:"name"`
	Username      string         `json:"username" gorm:"unique"`
	Email         string         `json:"email" gorm:"unique;non null"`
	Password      string         `json:"password" gorm:"non null"`
	ProfileImg    string         `json:"profileimg"`
	Role          string         `json:"role" gorm:"default:'USER'"`
	Bio           string         `json:"bio"`
	Comments      []Comment      `json:"comments"`
	Notifications []Notification `json:"notifications" gorm:"constraint:OnDelete:CASCADE;"`
	Posts         []Post         `json:"posts" gorm:"constraint:OnDelete:CASCADE;"`
	Following     []*User        `json:"following" gorm:"many2many:user_following;"`
	Followers     []*User        `json:"followers" gorm:"many2many:user_followers"`
	Resetter      *string        `json:"resetter"`
}

func (User) TableName() {
}
