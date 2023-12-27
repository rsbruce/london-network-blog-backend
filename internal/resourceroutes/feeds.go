package resourceroutes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (svc *Service) GetPersonalFeed(w http.ResponseWriter, r *http.Request) {
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

	id, err := svc.AuthData.GetUserId(r)
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