package dto

import "github.com/Cedar-81/swype/models"

type CreateCommentDto struct {
	Text string `json:"content"`
}

type UpdateCommentDto struct {
	Text string `json:"content"`
}

type CommentResponseDto struct {
	Text     string `json:"content"`
	PostID   uint
	ParentID *uint
	Parent   *models.Comment
	Replies  []models.Comment
}
