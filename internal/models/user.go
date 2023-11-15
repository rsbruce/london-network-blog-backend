package models

import (
	"errors"
	"log"
)

var (
	ErrFetchingUser = errors.New("Could not retrieve user")
)

type User struct {
	ID              []byte `json:"id"`
	Handle          string `json:"handle"`
	Blurb           string `json:"blurb"`
	Display_name    string `json:"display_name"`
	Display_picture string `json:"display_picture"`
	User_role       string `json:"user_role"`
	Created_at      string `json:"created_at"`
}
type UserProfile struct {
	User        User           `json:"user"`
	LatestPosts []FeedItemPost `json:"posts"`
}

type UserStore interface {
	GetUser(string) (User, error)
}

type UserService struct {
	Store UserStore
}

func NewUserService(store UserStore) *UserService {
	return &UserService{
		Store: store,
	}
}

func (s *UserService) GetUser(handle string) (User, error) {
	user, err := s.Store.GetUser(handle)
	if err != nil {
		log.Fatalf("An error occured fetching the latest post feed: %s", err.Error())
		return User{}, ErrFetchingUser
	}
	return user, nil
}
