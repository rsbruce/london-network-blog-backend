package resourceroutes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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