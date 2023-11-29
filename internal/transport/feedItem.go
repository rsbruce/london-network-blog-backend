package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (handler *HttpHandler) GetLatestAllAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	post_feed, err := handler.DB_conn.GetLatestPostFeed()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(post_feed)
}

func (handler *HttpHandler) GetLatestForAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	handle := params["handle"]

	author_post_feed, err := handler.DB_conn.GetFeedItemPostsForAuthor(handle)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(author_post_feed)
}
