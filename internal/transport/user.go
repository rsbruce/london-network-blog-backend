package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/database"

	"github.com/gorilla/mux"
)

func (handler *HttpHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	handle := params["handle"]

	author, err := handler.DB_conn.GetUser(handle)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("Could not get user with handle: %s", handle)})
		return
	}

	feed_item_posts, err := handler.DB_conn.GetFeedItemPostsForAuthor(handle)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("Failed to retrieve post items for author with handle: %s.", handle)})
		return
	}

	user_profile := database.UserProfile{
		User:        author,
		LatestPosts: feed_item_posts,
	}

	json.NewEncoder(w).Encode(user_profile)
}

func (handler *HttpHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint hit")
	params := mux.Vars(r)
	handle := params["handle"]

	decoder := json.NewDecoder(r.Body)

	var user_changes database.User
	err := decoder.Decode(&user_changes)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Invalid JSON payload for this route."})
		return
	}

	user, err := handler.DB_conn.UpdateUser(handle, user_changes)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not update user."})
		return
	}

	json.NewEncoder(w).Encode(user)
}
