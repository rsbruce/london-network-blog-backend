package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/models"
)

type PostFeedItemHandler struct {
	GetLatestAllAuthors func(w http.ResponseWriter, r *http.Request)
	GetLatestForAuthor  func(w http.ResponseWriter, r *http.Request)
}

func NewPostFeedItemHandler(service *models.FeedItemService) PostFeedItemHandler {
	var handler PostFeedItemHandler

	handler.GetLatestAllAuthors = func(w http.ResponseWriter, r *http.Request) {

		post_feed, err := service.Store.GetLatestPostFeed()
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(post_feed)
	}

	handler.GetLatestForAuthor = func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		handle := params["handle"]

		author_post_feed, err := service.Store.GetFeedItemPostsForAuthor(handle)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(author_post_feed)
	}

	return handler
}
