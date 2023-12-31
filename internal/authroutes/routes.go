package authroutes

import (
	"encoding/json"
	"log"
	"net/http"

	"rsbruce/blogsite-api/internal/authdata"
)

type Service struct {
	AuthData *authdata.Service
}

func (svc *Service) UserHandle(w http.ResponseWriter, r *http.Request) {
	var authResponse struct {
		Handle string `json:"handle"`
	}

	accessToken := r.Header.Get("Authorization")
	handle, err := svc.AuthData.GetHandleFromAccessToken(accessToken)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
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

	id, err := svc.AuthData.CheckPassword(credentials)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Could not authenticate")
		return
	}

	accessToken, refreshToken, err := svc.AuthData.GenerateTokensFromId(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error creating refresh token")
		return
	}

	var authResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	authResponse.AccessToken = accessToken
	authResponse.RefreshToken = refreshToken

	json.NewEncoder(w).Encode(authResponse)
}

func (svc *Service) RefreshAccess(w http.ResponseWriter, r *http.Request) {
	accessToken, refreshToken, err := svc.AuthData.GenerateTokensFromRefresh(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var authResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	authResponse.AccessToken = accessToken
	authResponse.RefreshToken = refreshToken

	json.NewEncoder(w).Encode(authResponse)
}

func (svc *Service) ResetPassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var creds authdata.LoginCredentials
	err := decoder.Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	err = svc.AuthData.UpdatePassword(id, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
