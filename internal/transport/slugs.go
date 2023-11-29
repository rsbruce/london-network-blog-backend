package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (handler *HttpHandler) GetSlugsForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	handle := params["handle"]

	slugs, err := handler.DB_conn.GetSlugsForUser(handle)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(slugs)
}
