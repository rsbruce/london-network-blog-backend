package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (handler *HttpHandler) GetTextContent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	textContent, err := handler.DB_conn.GetTextContent(slug)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "No Text Content found with slug: " + slug})
		return
	}
	json.NewEncoder(w).Encode(textContent)
}
