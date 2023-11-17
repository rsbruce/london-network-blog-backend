package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/models"
)

type UserProfileHandler struct {
	Get func(w http.ResponseWriter, r *http.Request)
}

func NewUserProfileHandler(userService *models.UserService, feedItemService *models.FeedItemService) UserProfileHandler {
	var handler UserProfileHandler
	handler.Get = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		handle := params["handle"]

		author, err := userService.Store.GetUser(handle)
		if err != nil {
			log.Fatal(err)
		}

		feed_item_posts, err := feedItemService.GetFeedItemPostsForAuthor(handle)
		if err != nil {
			log.Fatal(err)
		}

		user_profile := models.UserProfile{
			User:        author,
			LatestPosts: feed_item_posts,
		}

		json.NewEncoder(w).Encode(user_profile)
	}

	return handler
}
