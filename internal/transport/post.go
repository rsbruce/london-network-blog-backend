package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/database"

	"github.com/gorilla/mux"
)

func (handler *HttpHandler) GetPostPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	slug := params["slug"]

	post_with_user, err := handler.DB_conn.GetPostWithUser(slug)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("Post not found with slug: %s", slug)})
		return
	}

	json.NewEncoder(w).Encode(post_with_user)
}

func (handler *HttpHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var new_post database.Post
	err := decoder.Decode(&new_post)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Invalid JSON payload for this route."})
		return
	}

	post, err := handler.DB_conn.NewPost(new_post)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not add post."})
		return
	}

	json.NewEncoder(w).Encode(post)

}
