package transport

import (
	"encoding/json"
	"fmt"
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
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("No slugs found for user: %s", handle)})
		return
	}

	json.NewEncoder(w).Encode(slugs)
}
