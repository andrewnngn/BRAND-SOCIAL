package controllers

import (
	"net/http"
	"strconv"

	"github.com/Cedar-81/swype/dto"
	"github.com/Cedar-81/swype/models"
	"github.com/Cedar-81/swype/services"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	user := c.Value("user").(models.User)
	userid := strconv.FormatUint(uint64(user.ID), 10)
	postid := c.Param("postid")
	parentid := c.Query("parentid")

	// Validate input
	var input dto.CreateCommentDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := services.CreateComment(input, postid, parentid, userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comment})
}

func GetComments(c *gin.Context) {
	comments, err := services.GetAllComments()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})

}

func GetCommentsByPostId(c *gin.Context) {
	var comments []models.Comment
	postID := c.Param("postid")
	page := c.Query("page")
	limit := c.Query("limit")
	var err error

	switch {

	case page != "" && postID != "" || limit != "" && postID != "":
		comments, err = services.GetLimitedCommentsByPostId(postID, page, limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	default:
		comments, err = services.GetCommentsByPostId(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

func GetComment(c *gin.Context) {
	var comment models.Comment

	comment, err := services.GetCommentById(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comment})
}

func UpdateComment(c *gin.Context) {
	commentid := c.Param("id")
	user := c.Value("user").(models.User)
	userid := strconv.FormatUint(uint64(user.ID), 10)

	// Validate input
	var input dto.UpdateCommentDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.UpdateComment(input, userid, commentid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "comment updated successfully"})

}

func DeleteComment(c *gin.Context) {
	user := c.Value("user").(models.User)
	userid := strconv.FormatUint(uint64(user.ID), 10)

	err := services.DeleteComment(c.Param("id"), userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "comment deleted successfully"})
}

func LikeComment(c *gin.Context) {
	user := c.Value("user").(models.User)
	userid := strconv.FormatUint(uint64(user.ID), 10)
	postid := c.Param("commentid")

	err := services.LikeComment(postid, userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully liked comment"})
}
