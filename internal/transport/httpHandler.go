package transport

import (
	"net/http"
	"rsbruce/blogsite-api/internal/database"
)

type HttpHandler struct {
	DB_conn *database.Database
}

type ResponseMessage struct {
	Message string
}

func NewHttpHandler(db *database.Database) *HttpHandler {
	return &HttpHandler{DB_conn: db}
}

func (handler *HttpHandler) HandleCors(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}
	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.WriteHeader(http.StatusOK)
}
