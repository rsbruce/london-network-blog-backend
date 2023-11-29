package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (handler *HttpHandler) GetTextContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	slug := params["slug"]

	textContent, err := handler.DB_conn.GetTextContent(slug)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(textContent)
}
