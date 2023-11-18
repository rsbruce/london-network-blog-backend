package models

import (
	"errors"
	"log"
)

var (
	ErrFetchingPostWithUser = errors.New("Could not retrieve post with user")
	ErrSavingNewPost        = errors.New("Could not save new post")
)

type Post struct {
	ID         []byte `json:"id"`
	Author_id  []byte `json:"author_id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Content    string `json:"content"`
	Main_image string `json:"main_image"`
	Created_at string `json:"created_at"`
}

type PostPage struct {
	Post Post `json:"post"`
	User User `json:"user"`
}

type PostStore interface {
	GetPostWithUser(string) (PostPage, error)
	NewPost(Post) (Post, error)
}

type PostService struct {
	Store PostStore
}

func NewPostService(store PostStore) *PostService {
	return &PostService{
		Store: store,
	}
}

func (s *PostService) GetPostWithUser(slug string) (PostPage, error) {
	post, err := s.Store.GetPostWithUser(slug)
	if err != nil {
		log.Fatalf("An error occured fetching the latest post feed: %s", err.Error())
		return PostPage{}, ErrFetchingPostWithUser
	}
	return post, nil
}

func (s *PostService) NewPost(post Post) (Post, error) {
	post, err := s.Store.NewPost(post)
	if err != nil {
		log.Fatalf("An error occured saving the new post: %s", err.Error())
		return Post{}, ErrSavingNewPost
	}
	return post, nil
}
