package controllers

import (
	"net/http"
	"strconv"

	"github.com/Cedar-81/swype/dto"
	"github.com/Cedar-81/swype/initializers"
	"github.com/Cedar-81/swype/models"
	"github.com/Cedar-81/swype/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUser(c *gin.Context) {
	var user models.User

	//get user from req to make sure user is validated
	userdata := c.Value("user").(models.User)
	userID := strconv.FormatUint(uint64(userdata.ID), 10)

	userToFollowID := c.Param("tofollow") //person to follow

	// Get user if exist
	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input dto.UpdateUserDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check for and handle req parameters
	switch {
	case userToFollowID != "":
		err := services.AddFollower(userToFollowID, userID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	default:

		result := initializers.DB.Model(&user).Updates(input)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to update user info!"})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"data": "update successful"})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	users, err := services.GetAllUsers()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not fetch users!"})
		return
	}

	userresponse, err := services.HandleResponseForUsers(users)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userresponse})
}

func GetLoggedUser(c *gin.Context) {
	var user models.User

	//get user from req to make sure user is validated
	userdata := c.Value("user").(models.User)
	userID := strconv.FormatUint(uint64(userdata.ID), 10)

	user, err := services.GetUserByID(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User record not found!"})
		return
	}

	userresponse, err := services.HandleResponseForUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userresponse})

}

func GetUser(c *gin.Context) {
	username := c.Param("username")

	user, err := services.GetUserByUsername(username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User record not found!"})
		return
	}

	userresponse, err := services.HandleResponseForUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userresponse})
}

func DeleteUser(c *gin.Context) {
	//get user from req to make sure user is validated
	userdata := c.Value("user").(models.User)
	userID := strconv.FormatUint(uint64(userdata.ID), 10)

	err := services.DeleteUser(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.SetCookie("Auth", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted and logged out successfully"})
}

func UpdatePassword(c *gin.Context) {
	// Validate input
	var input dto.UpdatePasswordDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	err = services.ChangePass(input, string(hashedPass))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//make sure user is logged out
	c.SetCookie("Auth", "", -1, "", "", false, true)
}

func ResetPassword(c *gin.Context) {
	//make sure user is logged out
	c.SetCookie("Auth", "", -1, "", "", false, true)

	// Validate input
	var input dto.UpdatePasswordDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.GenerateResetter(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
