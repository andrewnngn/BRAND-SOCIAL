package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Cedar-81/swype/dto"
	"github.com/Cedar-81/swype/models"
	"github.com/Cedar-81/swype/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// get data email & pass from req body
	var input dto.SignupDto
	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})

		return
	}

	//check if user already exists
	_, err := services.GetUserByEmail(input.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User with email already exists",
		})

		return
	}

	// hash password
	hashedPass, err := services.HashPass(input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// create user
	err = services.CreateUser(input, hashedPass)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{})
}

func CreateAdmin(c *gin.Context) {
	// get data email & pass from req body
	var input dto.SignupDto
	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})

		return
	}

	//check if user already exists
	_, err := services.GetUserByEmail(input.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User with email already exists",
		})

		return
	}

	// hash password
	hashedPass, err := services.HashPass(input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}
	// create user
	err = services.CreateAdminUser(input, hashedPass)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{})
}

func LogIn(c *gin.Context) {
	//get email and password from req body
	var input dto.SignupDto
	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})

		return
	}

	//search for reqested user
	user, err := services.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})

		return
	}

	//compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})

		return
	}

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days expiration
	})
	// Sign token using a cryptographic algorithm and secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	//send back jwt token as cookie
	c.SetSameSite(http.SameSiteLaxMode)
	// 30 days expiration = token expiration
	c.SetCookie("Auth", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func LogOut(c *gin.Context) {
	c.SetCookie("Auth", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func Validate(c *gin.Context) {
	user := c.Value("user").(models.User)
	fmt.Println(user.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
	})
}
