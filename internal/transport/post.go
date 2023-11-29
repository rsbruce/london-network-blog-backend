package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/database"
)

func (handler *HttpHandler) GetPostPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	slug := params["slug"]

	post_with_user, err := handler.DB_conn.GetPostWithUser(slug)
	if err != nil {
		log.Fatal(err)
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "500 - parse error"})
		log.Fatal(err)
	}

	post, err := handler.DB_conn.NewPost(new_post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "500 - db error"})
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(post)

}
