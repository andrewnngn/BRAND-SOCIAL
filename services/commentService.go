package services

import (
	"errors"
	"strconv"

	"github.com/Cedar-81/swype/dto"
	"github.com/Cedar-81/swype/initializers"
	"github.com/Cedar-81/swype/models"
)

func GetCommentsByPostId(id string) ([]models.Comment, error) {
	var comments []models.Comment = []models.Comment{}

	_, err := GetPostById(id)

	if err != nil {
		return comments, errors.New("Unable to find post with id " + id)
	}

	if err == nil {
		initializers.DB.Where("post_id = ?", id).Find(&comments)
	}

	return comments, nil
}

func GetLimitedCommentsByPostId(id, page, limit string) ([]models.Comment, error) {
	var pageval int64 = 1
	var limitval int64 = 10
	var err error
	var comments []models.Comment

	_, err = GetPostById(id)

	if err != nil {
		return comments, errors.New("Unable to find post with id " + id)
	}

	if page != "" {
		pageval, err = strconv.ParseInt(page, 10, 64)
	}

	if limit != "" {
		limitval, err = strconv.ParseInt(limit, 10, 64)
	}

	if err != nil {
		return nil, errors.New("please enter a valid integer")
	}

	offset := (pageval - 1) * limitval

	if err == nil {
		initializers.DB.Where("post_id = ?", id).
			Limit(int(limitval)).Offset(int(offset)).
			Find(&comments)
	}

	return comments, nil
}

func GetAllComments() ([]models.Comment, error) {
	var comments []models.Comment
	result := initializers.DB.Find(&comments)

	if result.Error != nil {
		return comments, errors.New("unable to get comments")
	}

	return comments, nil
}

func GetCommentById(id string) (models.Comment, error) {
	var comment models.Comment

	if err := initializers.DB.
		Preload("Replies").
		Preload("LikedBy").
		Where("id = ?", id).
		First(&comment).Error; err != nil {
		return comment, err
	}

	return comment, nil
}

func GetCommentByUserAndCommentId(userid string, commentid string) (models.Comment, error) {
	var comment models.Comment

	if err := initializers.DB.Where("id = ? AND user_id = ?", commentid, userid).First(&comment).Error; err != nil {
		return comment, err
	}

	return comment, nil
}

func CreateComment(input dto.CreateCommentDto, postid string, parentid string, userid string) (models.Comment, error) {
	var comment models.Comment
	postidint, _ := strconv.ParseUint(postid, 36, 10)
	useridint, _ := strconv.ParseUint(userid, 36, 10)
	parentidint, _ := strconv.ParseUint(parentid, 36, 10)
	var parentiduint uint = uint(parentidint)

	//make sure post has content
	if input.Text == "" {
		return comment, errors.New("make sure content is not empty")
	}

	post, err := GetPostById(postid)

	if err != nil {
		return comment, errors.New("unable to get post with id " + postid)
	}

	if parentid != "" {
		comment, err = GetCommentById(parentid)
		if err != nil {
			return comment, errors.New("unable to get comment with id " + parentid)
		}
	}

	if parentid == "" {
		// Create comment
		comment = models.Comment{
			Text:   input.Text,
			PostID: uint(postidint),
			UserID: uint(useridint),
		}
	} else {
		// Create comment
		comment = models.Comment{
			Text:     input.Text,
			PostID:   uint(postidint),
			UserID:   uint(useridint),
			ParentID: &parentiduint,
		}
	}

	result := initializers.DB.Create(&comment)

	if result.Error != nil {
		return comment, errors.New("unable to create comment")
	}

	//check if author is the creator of the comment
	if post.UserID != uint(useridint) {
		commentuserid := strconv.FormatUint(uint64(comment.UserID), 10)
		err = PushNotification(commentuserid, "newcomment")
		if err != nil {
			return comment, err
		}
	}

	return comment, nil
}

func UpdateComment(input dto.UpdateCommentDto, userid string, commentid string) error {

	comment, err := GetCommentByUserAndCommentId(userid, commentid)

	if err != nil {
		return errors.New("unable to get comment with userid: " + userid + " and commentid: " + commentid)
	}

	result := initializers.DB.Model(&comment).Updates(input)

	if result.Error != nil {
		return errors.New("unable to update comment")
	}

	return nil
}

func DeleteComment(commentid string, userid string) error {

	/*
		Note: The only users with ability to delet a post
		are the author of the comment and the author of
		the post that contains the comment
	*/

	comment, err := GetCommentById(commentid)

	if err != nil {
		return errors.New("comment with id " + commentid + " not found")
	}

	postid := strconv.FormatUint(uint64(comment.PostID), 10)

	//find post by id and check its author
	post, err := GetPostById(postid)
	if err != nil {
		return errors.New("failed to get post with comment")
	}

	// delete post if user or post author is trying to delete post
	postauthor := strconv.FormatUint(uint64(post.UserID), 10)
	commentauthor := strconv.FormatUint(uint64(comment.UserID), 10)

	if postauthor == userid || commentauthor == userid {
		result := initializers.DB.Delete(&comment)

		if result.Error != nil {
			return errors.New("unable to delete comment")
		}
	} else {
		return errors.New("sorry you don't have the ability to delete this comment")
	}

	return nil
}

func LikeComment(commentid string, userid string) error {
	user, err := GetUserByID(userid)
	if err != nil {
		return errors.New("unable to get user")
	}

	comment, err := GetCommentById(commentid)
	if err != nil {
		return errors.New("unable to get comment")
	}

	//comfirm user haven't liked comment before
	for _, user := range comment.LikedBy {
		if userid == strconv.FormatUint(uint64(user.ID), 10) {
			return errors.New("unable to like comment more than once")
		}
	}

	comment.Likes++

	//handle adding user to post likedby list
	comment.LikedBy = append(comment.LikedBy, user)

	result := initializers.DB.Save(&comment)

	if result.Error != nil {
		return errors.New("unable to like comment")
	}

	commentuserid := strconv.FormatUint(uint64(comment.UserID), 10)
	err = PushNotification(commentuserid, "commentlike")

	if err != nil {
		return err
	}

	return nil
}
