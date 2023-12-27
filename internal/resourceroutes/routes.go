package resourceroutes

import (
	"encoding/json"
	"net/http"
	"rsbruce/blogsite-api/internal/authdata"
	"rsbruce/blogsite-api/internal/resourcedata"
	"strconv"

	"github.com/gorilla/mux"
)

type Service struct {
	AuthData     *authdata.Service
	ResourceData *resourcedata.Service
}

func (svc *Service) GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	authorHandle, ah := params["handle"]
	slug, s := params["slug"]

	if !(ah && s) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post, err := svc.ResourceData.GetPost(authorHandle, slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)

}

func (svc *Service) GetPersonalFeed(w http.ResponseWriter, r *http.Request) {
	var limit int
	var err error
	
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} 

	queryLimit := r.URL.Query().Get("limit")
	if queryLimit != "" {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	id, err := svc.AuthData.GetIdFromAccessToken(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	feed, err := svc.ResourceData.GetPersonalFeed(id, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(feed)
}

func (svc *Service) GetFeed(w http.ResponseWriter, r *http.Request) {
	var limit int
	var err error

	queryLimit := r.URL.Query().Get("limit")
	if queryLimit != "" {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	feed, err := svc.ResourceData.GetFeed(limit)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(feed)

}

func (svc *Service) GetSingleUserFeed(w http.ResponseWriter, r *http.Request) {
	var handle string
	var limit int
	var err error

	params := mux.Vars(r)
	handle = params["handle"]

	queryLimit := r.URL.Query().Get("limit")
	if queryLimit != "" {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	feed, err := svc.ResourceData.GetSingleUserFeed(handle, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(feed)

}

func (svc *Service) GetTextContent(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) UploadPhoto(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) CreatePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) EditPost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) DeletePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) RestorePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) EditUser(w http.ResponseWriter, r *http.Request) {

}
