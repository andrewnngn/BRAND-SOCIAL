package controllers

import (
	"net/http"

	"github.com/Cedar-81/swype/models"
	"github.com/Cedar-81/swype/services"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	var data []models.Model // empty array that can hold all models as they all implement the Model interface
	var users []models.User
	var err error

	searchtype := c.Query("type")
	searchstring := c.Param("searchparam")

	if searchstring == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a search string"})
		return
	}

	switch searchtype {
	case "video":
		data, err = services.SearchVideo(searchstring)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	case "image":
		data, err = services.SearchImage(searchstring)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	case "user":
		users, err = services.SearchUsers(searchstring)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	case "content":
		data, err = services.SearchTextContent(searchstring)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	default:
		users, err = services.SearchUsers(searchstring)
		data, err = services.SearchTextContent(searchstring)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "user": users})
}
