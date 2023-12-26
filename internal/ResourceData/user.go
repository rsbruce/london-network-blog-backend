package database

import (
	"database/sql"

	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserProfileRow struct {
	User        UserRow
	LatestPosts []FeedItemPostRow
}
type UserRow struct {
	ID              sql.NullInt64
	Handle          sql.NullString
	Blurb           sql.NullString
	Display_name    sql.NullString
	Display_picture sql.NullString
	User_role       sql.NullString
	Created_at      sql.NullString
}

type UserAuth struct {
	ID       int64  `json:"id"`
	Handle   string `json:"handle"`
	Password string `json:"password"`
}

type UserAuthRow struct {
	ID       sql.NullInt64
	Handle   sql.NullString
	Password sql.NullString
}

type User struct {
	ID              int64  `json:"id"`
	Handle          string `json:"handle"`
	Blurb           string `json:"blurb"`
	Display_name    string `json:"display_name"`
	Display_picture string `json:"display_picture"`
	User_role       string `json:"user_role"`
	Created_at      string `json:"created_at"`
}
type UserProfile struct {
	User        User           `json:"user"`
	LatestPosts []FeedItemPost `json:"posts"`
}

func userFromRow(row UserRow) User {
	return User{
		ID:              row.ID.Int64,
		Handle:          row.Handle.String,
		Blurb:           row.Blurb.String,
		Display_name:    row.Display_name.String,
		Display_picture: row.Display_picture.String,
		User_role:       row.User_role.String,
		Created_at:      row.Created_at.String,
	}
}

func (db *Database) GetUser(handle string) (User, error) {

	row := db.Client.QueryRow("SELECT id, handle, blurb, display_name, display_picture, user_role, created_at FROM user WHERE handle = ?", handle)
	var user_row UserRow

	if err := row.Scan(
		&user_row.ID,
		&user_row.Handle,
		&user_row.Blurb,
		&user_row.Display_name,
		&user_row.Display_picture,
		&user_row.User_role,
		&user_row.Created_at,
	); err != nil {
		return User{}, err
	}

	return userFromRow(user_row), nil
}

func (db *Database) UpdateUser(user User) (User, error) {

	var query_fields []string

	if user.Handle != "" {
		query_fields = append(query_fields, "handle = :handle")
	}
	if user.Display_name != "" {
		query_fields = append(query_fields, "display_name = :display_name")
	}
	if user.Blurb != "" {
		query_fields = append(query_fields, "blurb = :blurb")
	}

	query := "UPDATE user SET " + strings.Join(query_fields, ", ") + " WHERE id = :id"

	rows, err := db.Client.NamedQuery(query, user)

	if err != nil {
		return User{}, err
	}
	if err := rows.Close(); err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) UpdatePassword(userAuth UserAuth) error {

	password := []byte(userAuth.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		panic(err)
	}

	_, err = db.Client.Exec(`UPDATE user SET password = ? WHERE id = ?`, hashedPassword, userAuth.ID)

	return err
}

func (db *Database) UserHandleFromId(id int64) (string, error) {
	var handle string
	row := db.Client.QueryRow(`SELECT handle FROM user WHERE id = ?`, id)
	err := row.Scan(&handle)
	if err != nil {
		return "", err
	}

	return handle, nil
}

func (db *Database) UpdateDisplayPicture(id int64, imagePath string) error {
	_, err := db.Client.Exec(`UPDATE user SET display_picture = ? WHERE id = ?`, os.Getenv("HOST_NAME")+"/"+imagePath, id)

	return err
}
