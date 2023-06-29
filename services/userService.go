package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/Cedar-81/swype/dto"
	"github.com/Cedar-81/swype/initializers"
	"github.com/Cedar-81/swype/models"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

type FollowerOrFollowing int

const (
	Follower FollowerOrFollowing = iota
	Following
)

func CreateUser(input dto.SignupDto, hash string) error {
	var username string = input.Username

	// Extract username from email if username isn't provided
	if username == "" {
		emailarr := strings.Split(input.Email, "@")
		username = emailarr[0]
	}

	//create user
	user := models.User{
		Username: username,
		Role:     "USER",
		Email:    input.Email,
		Password: hash,
	}

	result := initializers.DB.Create(&user)

	return result.Error
}

func HashPass(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return string(hashedPass), errors.New("unable to hash password")
	}

	return string(hashedPass), nil

}

func CreateAdminUser(input dto.SignupDto, hash string) error {
	var username string = input.Username

	// Extract username from email if username isn't provided
	if username == "" {
		emailarr := strings.Split(input.Email, "@")
		username = emailarr[0]
	}

	//create user
	user := models.User{
		Username: username,
		Role:     "ADMIN",
		Email:    input.Email,
		Password: hash,
	}

	result := initializers.DB.Create(&user)

	return result.Error
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	if err := initializers.DB.
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User

	if err := initializers.DB.
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByID(id string) (models.User, error) {
	var user models.User

	if err := initializers.DB.
		Preload("Posts").
		Preload("Followers").
		Preload("Following").
		Preload("Notifications").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetFollowersOrFollowing(userID string, valueType FollowerOrFollowing) ([]*models.User, error) {
	user, err := GetUserByID(userID)

	if err != nil {
		return nil, errors.New("unable to get user with ID " + userID)
	}

	switch valueType {
	case Follower:
		// followers := user.Followers
		// for _, user := range followers {
		// 	user.Password = ""
		// 	user.Resetter = nil
		// }
		return user.Followers, nil
	case Following:
		// following := user.Following
		// for _, user := range following {
		// 	user.Password = ""
		// 	user.Resetter = nil
		// }
		return user.Following, nil
	default:
		return nil, errors.New("enter a valid value type")

	}
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := initializers.DB.Preload("Following").Preload("Followers").Preload("Posts").Find(&users)

	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func AddFollower(userID string, followerID string) error {
	/*
		Note in this instance the client making the request is the follower
		and the person to be followed is the user
	*/

	//Make sure user isn't trying to follow themselves
	if userID == followerID {
		return errors.New("sorry you can't follow yourself")
	}

	// Get the user and follower
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	follower, err := GetUserByID(followerID)
	if err != nil {
		return err
	}

	// Check if user is already followed
	for _, existingFollower := range user.Followers {
		if existingFollower.ID == follower.ID {
			return errors.New("you already follow this User")
		}
	}

	// Add the follower to the user's followers list
	user.Followers = append(user.Followers, &follower)

	// Add the user to the follower's following list
	follower.Following = append(follower.Following, &user)

	// Save the user and follower
	err = initializers.DB.Save(&user).Error
	if err != nil {
		return err
	}

	err = initializers.DB.Save(&follower).Error
	if err != nil {
		return err
	}

	return nil
}

func RemoveFollower(userID string, followerID string) error {
	//Make sure user isn't trying to unfollow themselves
	if userID == followerID {
		return errors.New("sorry you can't unfollow yourself")
	}

	// Get the user and follower
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	follower, err := GetUserByID(followerID)
	if err != nil {
		return err
	}

	// Remove follower from user's followers list
	if err := initializers.DB.Model(&user).Association("Followers").Delete(&follower); err != nil {
		return err
	}

	// Remove user from follower's following list
	if err := initializers.DB.Model(&follower).Association("Following").Delete(&user); err != nil {
		return err
	}

	return nil
}

func DeleteUser(userid string) error {
	user, err := GetUserByID(userid)

	if err != nil {
		return errors.New("user with id '" + userid + "' does not exist")
	}

	if err := initializers.DB.Delete(&user).Error; err != nil {
		return errors.New("couldn't complete delete")
	}

	return nil

}

func GenerateResetter(input dto.UpdatePasswordDto) error {
	var user models.User
	var err error
	length := 10

	if input.Username == "" {
		_, err = GetUserByEmail(input.Email)
		if err != nil {
			return errors.New("user with email does not exist")
		}
	}

	if input.Email == "" {
		_, err = GetUserByUsername(input.Email)
		if err != nil {
			return errors.New("user with username does not exist")
		}
	}

	//generate resetter
	randomBytes := make([]byte, length)

	_, err = rand.Read(randomBytes)
	if err != nil {
		return err
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	randomString = randomString[:length]

	//save to db
	user = models.User{
		Resetter: &randomString,
	}

	if input.Username == "" {
		err = initializers.DB.Where("email = ?", input.Email).Updates(&user).Error
	}

	if input.Email == "" {
		err = initializers.DB.Where("username = ?", input.Username).Updates(&user).Error
	}

	if err != nil {
		return errors.New("unable to set resetter")
	}

	return nil
}

func ChangePass(input dto.UpdatePasswordDto, hash string) error {
	var user models.User
	var err error

	//check if email or username is entered and user with credentials exist
	if input.Email == "" {
		user, err = GetUserByUsername(input.Username)
	}

	if input.Username == "" {
		user, err = GetUserByEmail(input.Email)
	}

	if err != nil {
		return errors.New("please make sure you entered a valid email or username ")

	}

	//confirm that resseter is correcet
	if input.Resetter == *user.Resetter && len(*user.Resetter) > 0 {
		user = models.User{
			Password: hash,
			Resetter: nil,
		}

		if input.Email == "" {
			err = initializers.DB.Model(&user).
				Where("email = ?", input.Username).
				Updates(&user).Error
			if err != nil {
				return errors.New("unable to update password")
			}
			return nil
		}

		if input.Username == "" {
			err = initializers.DB.Model(&user).
				Where("email = ?", input.Email).
				Updates(&user).Error
			if err != nil {
				return errors.New("unable to update password")
			}
			return nil
		}
	}

	return errors.New("invalid resetter value")
}

func HandleResponseForUser(users models.User) (dto.UserResponseDto, error) {
	var userresponse dto.UserResponseDto
	config := &mapstructure.DecoderConfig{
		Result:           &userresponse,
		TagName:          "json",
		WeaklyTypedInput: true, // Allow weak type conversions
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return userresponse, errors.New("Unable to create decoder for response data")
	}

	if err := decoder.Decode(&users); err != nil {
		return userresponse, errors.New("Unable to decode response data")
	}

	return userresponse, nil
}

func HandleResponseForUsers(users []models.User) ([]dto.UserResponseDto, error) {
	var userresponse []dto.UserResponseDto
	config := &mapstructure.DecoderConfig{
		Result:           &userresponse,
		TagName:          "json",
		WeaklyTypedInput: true, // Allow weak type conversions
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return userresponse, errors.New("Unable to create decoder for response data")
	}

	if err := decoder.Decode(&users); err != nil {
		return userresponse, errors.New("Unable to decode response data")
	}

	return userresponse, nil
}
