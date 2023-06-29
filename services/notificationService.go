package services

import (
	"errors"

	"github.com/Cedar-81/swype/initializers"
	"github.com/Cedar-81/swype/models"
)

// get notifications with limit
func GetNotificationsByUserID(userid string) ([]models.Notification, error) {
	var user models.User
	user, err := GetUserByID(userid)

	if err != nil {
		return nil, errors.New("unable to find user")
	}

	return user.Notifications, nil
}

// create notification
func CreateNotification(userid string, message string, notiftype string) (models.Notification, error) {
	var notification models.Notification

	user, err := GetUserByID(userid)
	if err != nil {
		return notification, errors.New("unable to find user")
	}

	notification = models.Notification{
		NotifType: notiftype,
		Message:   message,
		UserID:    user.ID,
	}

	result := initializers.DB.Create(&notification)

	if result.Error != nil {
		return notification, errors.New("unable to create new notification")
	}

	return notification, nil
}

// update notification
func PushNotification(userid string, notiftype string) error {
	var err error
	var notification models.Notification
	user, err := GetUserByID(userid)

	switch notiftype {
	case "postlike":
		msg := user.Username + " liked your post"
		notification, err = CreateNotification(userid, msg, "postlike")
	case "newcomment":
		msg := user.Username + " commented on your post"
		notification, err = CreateNotification(userid, msg, "newcomment")
	case "commentlike":
		msg := user.Username + " liked your comment"
		notification, err = CreateNotification(userid, msg, "commentlike")
	}

	if err != nil {
		return err
	}

	user.Notifications = append(user.Notifications, notification)

	err = initializers.DB.Save(&user).Error
	if err != nil {
		return errors.New("unable to save notification update")
	}

	return nil
}
