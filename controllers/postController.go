package controllers

import (
	"net/http"
	"strconv"

	"github.com/Cedar-81/swype/dto"
	"github.com/Cedar-81/swype/initializers"
	"github.com/Cedar-81/swype/models"
	"github.com/Cedar-81/swype/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePost(c *gin.Context) {
	user := c.Value("user").(models.User)
	userid := strconv.FormatUint(uint64(user.ID), 10)

	// Validate input
	var input dto.PostDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := services.CreatePost(&input, userid)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to find user"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to create post"})
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": post})
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	userID := c.Query("userid")
	page := c.Query("page")
	limit := c.Query("limit")
	var err error

	switch {
	case userID != "":
		posts, err = services.GetPostsByUserId(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	case page != "" || limit != "":
		posts, err = services.GetLimitedPost(page, limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	case page != "" && userID != "" || limit != "" && userID != "":
		posts, err = services.GetLimitedPostByUserID(page, limit, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	default:
		posts, err = services.GetAllPost()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func GetPost(c *gin.Context) {
	var post models.Post

	post, err := services.GetPostById(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": post})
}

func UpdatePost(c *gin.Context) {

	user := c.Value("user").(models.User)
	userid := strconv.FormatUint(uint64(user.ID), 10)

	// Validate input
	var input dto.PostDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := services.UpdatePost(input, userid, c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update Successful"})
}

func DeletePost(c *gin.Context) {
	var post models.Post

	// Get model if exist
	post, err := services.GetPostById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	initializers.DB.Delete(&post)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func LikePost(c *gin.Context) {
	user := c.Value("user").(models.User)
	userid := strconv.FormatUint(uint64(user.ID), 10)
	postid := c.Param("postid")

	err := services.LikePost(postid, userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully liked post"})
}
