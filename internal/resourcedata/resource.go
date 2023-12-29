package resourcedata

import (
	"log"
	"os"
	"strings"

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
		&userRow.Display_name,
		&userRow.Display_picture,
		&userRow.Blurb,
		&userRow.User_role,
		&userRow.Created_at,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return userRow.GetUser(), nil
}

func (svc *Service) UpdateDisplayPicture(id int64, imagePath string) error {
	_, err := svc.DbConn.Exec(`UPDATE user SET display_picture = ? WHERE id = ?`, os.Getenv("HOST_NAME")+"/"+imagePath, id)

	return err
}

func (svc *Service) UpdatePostImage(author_id int64, slug string, imagePath string) error {
	_, err := svc.DbConn.Exec(`UPDATE post SET main_image = ? WHERE author_id = ? AND slug = ?`, os.Getenv("HOST_NAME")+"/"+imagePath, author_id, slug)

	return err
}
func (svc *Service) EditUser(user User) error {
	var (
		columns []string
		values  []interface{}
	)
	if user.Blurb != "" {
		columns = append(columns, "blurb")
		values = append(values, user.Blurb)
	}
	if user.Display_name != "" {
		columns = append(columns, "display_name")
		values = append(values, user.Blurb)
	}
	if user.Display_picture != "" {
		columns = append(columns, "display_picture")
		values = append(values, user.Display_picture)
	}
	values = append(values, user.ID)

	if len(columns) > 0 {
		query := `
			UPDATE user SET ` + strings.Join(columns, " = ?, ") + ` = ?  WHERE id = ?`

		_, err := svc.DbConn.Exec(query, values...)

		return err
	}

	return nil

}
