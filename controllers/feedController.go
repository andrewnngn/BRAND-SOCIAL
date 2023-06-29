package controllers

import (
	"net/http"
	"strconv"

	"github.com/Cedar-81/swype/models"
	"github.com/Cedar-81/swype/services"
	"github.com/gin-gonic/gin"
)

func GetFeedPosts(c *gin.Context) {
	user := c.Value("user").(models.User)
	userid := strconv.FormatUint(uint64(user.ID), 10)

	posts, err := services.GetPostsFromFollowedUsers(userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}
