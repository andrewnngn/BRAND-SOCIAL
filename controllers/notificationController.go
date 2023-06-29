package controllers

import (
	"net/http"
	"strconv"

	"github.com/Cedar-81/swype/models"
	"github.com/Cedar-81/swype/services"
	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	//get user from req to make sure user is validated
	userdata := c.Value("user").(models.User)
	userID := strconv.FormatUint(uint64(userdata.ID), 10)

	notifications, err := services.GetNotificationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}
