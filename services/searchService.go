package services

import (
	"errors"
	"sort"

	"github.com/Cedar-81/swype/initializers"
	"github.com/Cedar-81/swype/models"
)

func SearchUsers(searchparam string) ([]models.User, error) {
	var users []models.User

	result := initializers.DB.Where("users.username LIKE ?", "%"+searchparam+"%").Find(&users)

	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func SearchPosts(searchparam string) ([]models.Post, error) {
	var posts []models.Post

	result := initializers.DB.
		Where("posts.text_content LIKE ?", "%"+searchparam+"%").Find(&posts)

	if result.Error != nil {
		return posts, errors.New("unable to successfully search Posts")
	}

	return posts, nil
}

func SearchComments(searchparam string) ([]models.Comment, error) {
	var comments []models.Comment

	result := initializers.DB.
		Where("comments.text LIKE ?", "%"+searchparam+"%").Find(&comments)

	if result.Error != nil {
		return comments, errors.New("unable to successfully search Comments")
	}

	return comments, nil
}

func SearchVideo(searchparam string) ([]models.Model, error) {
	var data []models.Model
	var posts []models.Post
	var comments []models.Comment
	var err error

	posts, err = SearchPosts(searchparam)
	if err != nil {
		return data, err
	}

	comments, err = SearchComments(searchparam)
	if err != nil {
		return data, err
	}

	for _, post := range posts {
		if post.VideoURL != "" {
			data = append(data, post)
		}
	}

	for _, comment := range comments {
		if comment.VideoURL != "" {
			data = append(data, comment)
		}
	}

	return data, nil

}

func SearchImage(searchparam string) ([]models.Model, error) {
	var data []models.Model
	var posts []models.Post
	var comments []models.Comment
	var err error

	posts, err = SearchPosts(searchparam)
	if err != nil {
		return data, err
	}

	comments, err = SearchComments(searchparam)
	if err != nil {
		return data, err
	}

	for _, post := range posts {
		if post.ImageURL != "" {
			data = append(data, post)
		}
	}

	for _, comment := range comments {
		if comment.ImageURL != "" {
			data = append(data, comment)
		}
	}

	return data, nil

}

func SearchTextContent(searchparam string) ([]models.Model, error) {
	var data []models.Model
	var posts []models.Post
	var comments []models.Comment
	var err error

	posts, err = SearchPosts(searchparam)
	if err != nil {
		return data, err
	}

	comments, err = SearchComments(searchparam)
	if err != nil {
		return data, err
	}

	for _, post := range posts {
		data = append(data, post)
	}

	for _, comment := range comments {
		data = append(data, comment)
	}

	SortModelsByUpdatedAtDesc(data)

	return data, nil

}

func SortModelsByUpdatedAtDesc(models []models.Model) {
	sort.Slice(models, func(i, j int) bool {
		return models[i].GetUpdatedAt().After(models[j].GetUpdatedAt())
	})
}
