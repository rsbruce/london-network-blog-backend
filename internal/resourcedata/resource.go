package resourcedata

import (
	"fmt"
	"log"
	"os"
	"reflect"
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
	fields := getFieldsToUpdate(user, []string{"id", "created_at", "password", "user_role"})
	assignments := make([]string, len(fields))
	for i, field := range fields {
		assignments[i] = fmt.Sprintf("%s = :%s", field, field)
	}

	query := `UPDATE user SET ` +
		strings.Join(assignments, ", ") +
		` WHERE id = :id`

	rows, err := svc.DbConn.NamedQuery(query, user)
	if err != nil {
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}

	return nil
}

func getFieldsToUpdate(obj interface{}, exclude []string) []string {
	v := reflect.ValueOf(obj)
	val := reflect.Indirect(v)
	num_fields := v.NumField()

	var names []string
	var name string
	var field reflect.Value

	for i := 0; i < num_fields; i++ {
		field = v.Field(i)
		if (field.Type().String() == "string" && field.String() != "") ||
			(field.Type().String() == "int64" && field.Int() != 0) {
			name = strings.ToLower(val.Type().Field(i).Name)
			names = append(names, name)
		}
	}

	return unqiueStrings(names, exclude)
}

func unqiueStrings(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
