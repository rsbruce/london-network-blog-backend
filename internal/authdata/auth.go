package authdata

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DbConn *sqlx.DB
}

type UserClaims struct {
	ID int64
	jwt.StandardClaims
}

type LoginCredentials struct {
	Handle   string
	Password string
}

func (svc *Service) CheckPassword(creds LoginCredentials) (int64, error) {

	password := []byte(creds.Password)

	var storedHash []byte
	var id int64
	row := svc.DbConn.QueryRow(`SELECT password, id FROM user WHERE handle = ?`, creds.Handle)
	err := row.Scan(&storedHash, &id)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(storedHash, password)

	return id, err
}

func (svc *Service) GenerateTokensFromId(id int64) (string, string, error) {
	userClaims := UserClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 120).Unix(),
		},
	}
	signedAccessToken, err := NewAccessToken(userClaims)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}

	signedRefreshToken, err := NewRefreshToken(refreshClaims)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func (svc *Service) GetUserId(r *http.Request) (int64, error) {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		return 0, errors.New("Empty access token")
	}

	userClaims, err := ParseAccessToken(accessToken)
	if err != nil {
		return 0, err
	}
	return userClaims.ID, nil
}

func (svc *Service) GetHandleFromAccessToken(accessToken string) (string, error) {
	var handle string

	userClaims, err := ParseAccessToken(accessToken)
	if err != nil {
		return "", err
	}
	id := userClaims.ID

	row := svc.DbConn.QueryRow(`SELECT handle FROM user WHERE id = ?`, id)
	err = row.Scan(&handle)
	if err != nil {
		return "", err
	}

	return handle, nil
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedAccessToken.Claims.(*UserClaims), err
}

func ParseRefreshToken(refreshToken string) (*jwt.StandardClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedRefreshToken.Claims.(*jwt.StandardClaims), err
}

func (svc *Service) CanEditPost(r *http.Request, post_id int64) bool {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		return false
	}
	user_claims, err := ParseAccessToken(accessToken)
	if err != nil {
		return false
	}
	token_author_id := user_claims.ID
	true_author_id, err := svc.GetAuthorIdFromPostId(post_id)
	if err != nil {
		return false
	}

	return token_author_id == true_author_id
}

func (svc *Service) GetAuthorIdFromPostId(id int64) (int64, error) {
	row := svc.DbConn.QueryRow(`SELECT author_id FROM post WHERE id = ?`, id)
	var author_id int64
	err := row.Scan(&author_id)
	if err != nil {
		return 0, err
	}

	return author_id, nil
}
