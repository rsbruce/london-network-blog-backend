package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/models"

	"github.com/gorilla/mux"
)

type UserProfileHandler struct {
	Get             func(w http.ResponseWriter, r *http.Request)
	Update          func(w http.ResponseWriter, r *http.Request)
	userService     *models.UserService
	feedItemService *models.FeedItemService
}
type ResponseMessage struct {
	Message string
}

func (handler *UserProfileHandler) getUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	handle := params["handle"]

	author, err := handler.userService.Store.GetUser(handle)
	if err != nil {
		log.Fatal(err)
	}

	feed_item_posts, err := handler.feedItemService.GetFeedItemPostsForAuthor(handle)
	if err != nil {
		log.Fatal(err)
	}

	user_profile := models.UserProfile{
		User:        author,
		LatestPosts: feed_item_posts,
	}

	json.NewEncoder(w).Encode(user_profile)
}

func (handler *UserProfileHandler) updateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint hit")
	params := mux.Vars(r)
	handle := params["handle"]

	decoder := json.NewDecoder(r.Body)

	var user_changes models.User
	err := decoder.Decode(&user_changes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "500 - parse error"})
		log.Fatal(err)
	}

	user, err := handler.userService.UpdateUser(handle, user_changes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "500 - db error"})
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(user)
}

func NewUserProfileHandler(userService *models.UserService, feedItemService *models.FeedItemService) UserProfileHandler {
	var handler UserProfileHandler
	handler.userService = userService
	handler.feedItemService = feedItemService
	handler.Get = handler.getUserProfile
	handler.Update = handler.updateUserProfile
	return handler
}
