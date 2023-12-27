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