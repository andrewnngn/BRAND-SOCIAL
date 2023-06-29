package dto

import "github.com/Cedar-81/swype/models"

type SignupDto struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

//TODO implement struct update email and password as both needs email and password

type UpdateUserDto struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	ProfileImg string `json:"profileimg"`
	Bio        string `json:"bio"`
}

type UpdatePasswordDto struct {
	Username string `json:"username"`
	Resetter string `json:"resetter"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponseDto struct {
	Name          string                `json:"name"`
	Username      string                `json:"username"`
	Email         string                `json:"email"`
	ProfileImg    string                `json:"profileimg"`
	Role          string                `json:"role"`
	Bio           string                `json:"bio"`
	Comments      []models.Comment      `json:"comments"`
	Notifications []models.Notification `json:"notifications"`
	Posts         []models.Post         `json:"posts"`
	Following     []*models.User        `json:"following"`
	Followers     []*models.User        `json:"followers"`
}
