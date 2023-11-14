package postFeedItem

import (
	"errors"
	"log"
)

var (
	ErrFetchingLatestPostFeed = errors.New("Could not retrieve latest post feed")
	ErrFetchingAuthorPostFeed = errors.New("Could not retrieve latest post feed for author")
)

type FeedItem struct {
	Post FeedItemPost `json:"post"`
	User FeedItemUser `json:"user"`
}

type FeedItemPost struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Slug       string `json:"slug"`
	Created_at string `json:"created_at"`
}

type FeedItemUser struct {
	Display_name    string `json:"display_name"`
	Display_picture string `json:"display_picture"`
	Handle          string `json:"handle"`
}

type PostFeedItemStore interface {
	GetLatestPostFeed() ([]FeedItem, error)
	GetFeedItemPostsForAuthor(string) ([]FeedItemPost, error)
}

type PostFeedItemService struct {
	Store PostFeedItemStore
}

func NewService(store PostFeedItemStore) *PostFeedItemService {
	return &PostFeedItemService{
		Store: store,
	}
}

func (s *PostFeedItemService) GetLatestPostFeed() ([]FeedItem, error) {
	feedItems, err := s.Store.GetLatestPostFeed()
	if err != nil {
		log.Fatalf("An error occured fetching the latest post feed: %s", err.Error())
		return nil, ErrFetchingLatestPostFeed
	}
	return feedItems, nil
}

func (s *PostFeedItemService) GetFeedItemPostsForAuthor(handle string) ([]FeedItemPost, error) {
	feedItemPosts, err := s.Store.GetFeedItemPostsForAuthor(handle)
	if err != nil {
		log.Fatalf("An error occured fetching the latest post feed for %s: %s", handle, err.Error())
		return nil, ErrFetchingAuthorPostFeed
	}
	return feedItemPosts, nil
}
