package authroutes

import (
	"encoding/json"
	"log"
	"net/http"

	"rsbruce/blogsite-api/internal/authdata"
)

type Service struct {
	AuthDataService *authdata.Service
}

func (svc *Service) UserHandle(w http.ResponseWriter, r *http.Request) {
	var authResponse struct {
		Message string `json:"message"`
		Handle  string `json:"handle"`
	}

	accessToken := r.Header.Get("Authorization")
	handle, err := svc.AuthDataService.GetHandleFromAccessToken(accessToken)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		authResponse.Message = "Failed"
	} else {
		authResponse.Message = "Success"
		authResponse.Handle = handle
	}

	json.NewEncoder(w).Encode(authResponse)
}

func (svc *Service) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var credentials authdata.LoginCredentials
	err := decoder.Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid JSON payload for this route.")
		return
	}

	id, err := svc.AuthDataService.CheckPassword(credentials)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Could not authenticate")
		return
	}

	accessToken, refreshToken, err := svc.AuthDataService.GenerateTokensFromId(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error creating refresh token")
		return
	}

	var auth_response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	auth_response.AccessToken = accessToken
	auth_response.RefreshToken = refreshToken

	json.NewEncoder(w).Encode(auth_response)
}