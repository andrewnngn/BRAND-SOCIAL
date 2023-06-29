package services

import (
	"errors"
	"strconv"

	"github.com/Cedar-81/swype/dto"
	"github.com/Cedar-81/swype/initializers"
	"github.com/Cedar-81/swype/models"
)

func GetLimitedPost(page, limit string) ([]models.Post, error) {
	var pageval int64 = 1
	var limitval int64 = 10
	var err error
	var posts []models.Post

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

	initializers.DB.
		Limit(int(limitval)).Offset(int(offset)).
		Preload("LikedBy").
		Preload("Comments.Replies").
		Preload("Comments.LikedBy").
		Find(&posts)

	return posts, nil
}

func GetLimitedPostByUserID(page, limit, userid string) ([]models.Post, error) {
	var pageval int64 = 1
	var limitval int64 = 10
	var err error
	var posts []models.Post

	_, err = GetUserByID(userid)
	if err != nil {
		return posts, err
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

	initializers.DB.
		Where("user_id = ?", userid).
		Limit(int(limitval)).Offset(int(offset)).
		Preload("LikedBy").
		Preload("Comments.Replies").
		Preload("Comments.LikedBy").
		Find(&posts)

	return posts, nil
}

func GetAllPost() ([]models.Post, error) {
	var posts []models.Post
	result := initializers.DB.
		Preload("LikedBy").
		Preload("Comments.Replies").
		Preload("Comments.LikedBy").
		Find(&posts)

	if result.Error != nil {
		return posts, errors.New("unable to get posts")
	}

	return posts, nil
}

func GetPostsByUserId(id string) ([]models.Post, error) {
	var posts []models.Post = []models.Post{}

	_, err := GetUserByID(id)
	if err != nil {
		return posts, err
	}

	result := initializers.DB.Where("user_id = ?", id).Find(&posts)
	if result.Error != nil {
		return posts, errors.New("unable to get posts by: " + id)
	}

	return posts, nil
}

func GetPostById(id string) (models.Post, error) {
	var post models.Post

	if err := initializers.DB.Preload("Comments").Where("id = ?", id).First(&post).Error; err != nil {
		return post, err
	}

	return post, nil
}

func CreatePost(input *dto.PostDto, userid string) (models.Post, error) {
	var post models.Post
	useridnum, _ := strconv.ParseUint(userid, 10, 64)

	//make sure post has content
	if input.TextContent == "" && input.ImageURL == "" && input.VideoURL == "" {
		return post, errors.New("make sure at least one post content is inputed")
	}

	//get user
	user, err := GetUserByID(userid)
	if err != nil {
		return post, err
	}

	// Create book
	post = models.Post{
		TextContent: input.TextContent,
		ImageURL:    input.ImageURL,
		VideoURL:    input.VideoURL,
		UserID:      uint(useridnum),
		User:        user,
	}
	if err := initializers.DB.Create(&post).Error; err != nil {
		return post, err
	}

	return post, nil
}

func GetSinglePostByUser(userid string, postid string) (models.Post, error) {
	var post models.Post

	if err := initializers.DB.Preload("Comments").Where("id = ? AND user_id = ?", postid, userid).First(&post).Error; err != nil {
		return post, err
	}

	return post, nil
}

func GetPostByID(postid string) (models.Post, error) {
	var post models.Post

	if err := initializers.DB.Preload("Comments").
		Preload("LikedBy").
		Where("id = ?", postid).
		First(&post).Error; err != nil {
		return post, errors.New("unabele to get post with id " + postid)
	}

	return post, nil
}

func UpdatePost(input dto.PostDto, userid string, postid string) (models.Post, error) {
	var newinput models.Post
	useridnum, _ := strconv.ParseUint(userid, 10, 64)

	//check if post with id belongs to logged in user
	post, err := GetSinglePostByUser(userid, postid)
	if err != nil {
		return post, errors.New("User does not have post with id '" + userid + "'")
	}

	//get user
	user, err := GetUserByID(userid)
	if err != nil {
		return post, errors.New("Unable to get user with id '" + postid + "'")
	}

	//modify input
	newinput = models.Post{
		TextContent: input.TextContent,
		ImageURL:    input.ImageURL,
		VideoURL:    input.VideoURL,
		UserID:      uint(useridnum),
		User:        user,
	}

	result := initializers.DB.Model(&post).Updates(newinput)
	if result.Error != nil {
		return post, errors.New("unable to update post")
	}

	return newinput, nil
}

func GetPostsFromFollowedUsers(userid string) ([]models.Post, error) {
	followingList, err := GetFollowersOrFollowing(userid, Following) // list of users follwed by user
	userids := []string{userid}                                      //id's of users in following list including user
	var posts []models.Post

	if err != nil {
		return nil, errors.New("unable to get following")
	}

	for _, user := range followingList {
		id := strconv.Itoa(int(user.ID))
		userids = append(userids, id)
	}

	initializers.DB.
		Where("user_id IN ?", userids).
		Order("updated_at DESC").
		Find(&posts)

	return posts, nil
}

func LikePost(postid string, userid string) error {
	user, err := GetUserByID(userid)
	if err != nil {
		return errors.New("unable to get user")
	}

	post, err := GetPostByID(postid)
	if err != nil {
		return errors.New("unable to get post")
	}

	//confirm user havent liked post before
	for _, user := range post.LikedBy {
		if userid == strconv.FormatUint(uint64(user.ID), 10) {
			return errors.New("unable to like post more than once")
		}
	}

	post.Likes++

	//handle adding user to post likedby list
	post.LikedBy = append(post.LikedBy, user)

	result := initializers.DB.Save(&post)

	if result.Error != nil {
		return errors.New("unable to like post")
	}

	postuserid := strconv.FormatUint(uint64(post.UserID), 10)
	err = PushNotification(postuserid, "postlike")

	if err != nil {
		return err
	}

	return nil
}
