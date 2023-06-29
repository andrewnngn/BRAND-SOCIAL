package middleware

import (
	"net/http"

	"github.com/Cedar-81/swype/models"
	"github.com/gin-gonic/gin"
)

func CheckSuperAdmin(c *gin.Context) {
	user := c.Value("user").(models.User)

	if user.Role == "SUPER" {
		//continue
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
