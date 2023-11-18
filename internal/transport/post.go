package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/models"
)

type PostHandler struct {
	GetPostPage func(w http.ResponseWriter, r *http.Request)
	NewPost     func(w http.ResponseWriter, r *http.Request)
	postService *models.PostService
}

func (handler *PostHandler) getPostPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	slug := params["slug"]

	post_with_user, err := handler.postService.Store.GetPostWithUser(slug)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(post_with_user)
}

func (handler *PostHandler) newPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var new_post models.Post
	err := decoder.Decode(&new_post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "500 - parse error"})
		log.Fatal(err)
	}

	post, err := handler.postService.NewPost(new_post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "500 - db error"})
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(post)

}

func NewPostHandler(post_svc *models.PostService) PostHandler {
	var handler PostHandler
	handler.postService = post_svc
	handler.GetPostPage = handler.getPostPage
	handler.NewPost = handler.newPost

	return handler
}
