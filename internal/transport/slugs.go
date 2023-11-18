package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/models"
)

type SlugsHandler struct {
	Get func(w http.ResponseWriter, r *http.Request)
}

func NewSlugsHandler(slugs_svc *models.SlugsService) SlugsHandler {
	var handler SlugsHandler
	handler.Get = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		params := mux.Vars(r)
		handle := params["handle"]

		slugs, err := slugs_svc.Store.GetSlugsForUser(handle)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(slugs)
	}

	return handler
}
