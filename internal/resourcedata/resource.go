package resourcedata

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	DbConn *sqlx.DB
}

func (svc *Service) GetTextContent(slug string) (string, error) {
	var content string
	
	row := svc.DbConn.QueryRow(`SELECT content FROM text_content WHERE slug = ?`, slug)
	err := row.Scan(&content)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return content, nil
}

func (svc *Service) GetUser(handle string) (*User, error) {
	var userRow UserRow

	row := svc.DbConn.QueryRow(`
		SELECT display_name, display_picture, blurb, user_role, created_at
		FROM user
		WHERE handle = ?
	`, handle)

	err := row.Scan(
		&userRow.DisplayName,
		&userRow.DisplayPicture,
		&userRow.Blurb,
		&userRow.UserRole,
		&userRow.CreatedAt,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return userRow.GetUser(), nil
} 