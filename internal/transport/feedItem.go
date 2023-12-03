package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (handler *HttpHandler) GetLatestAllAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	post_feed, err := handler.DB_conn.GetLatestPostFeed()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Query for all authors' latest posts failed"})
		return
	}
	json.NewEncoder(w).Encode(post_feed)
}

func (handler *HttpHandler) GetLatestForAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	handle := params["handle"]

	author_post_feed, err := handler.DB_conn.GetActiveFeedItemPostsForAuthor(handle)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("Failed to retrieve post items for author with handle: %s. Check author exists and has posts.", handle)})
		return
	}
	json.NewEncoder(w).Encode(author_post_feed)
}
