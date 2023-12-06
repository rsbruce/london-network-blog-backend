package database

import (
	"database/sql"

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

func (db *Database) UpdateUser(handle string, user User) (User, error) {

	user_row := UserRow{
		Handle:       sql.NullString{String: handle, Valid: true},
		Display_name: sql.NullString{String: user.Display_name, Valid: true},
		Blurb:        sql.NullString{String: user.Blurb, Valid: true},
	}

	rows, err := db.Client.NamedQuery(
		`UPDATE user SET
		blurb = :blurb,
		display_name = :display_name
		WHERE handle = :handle`,
		user_row,
	)

	if err != nil {
		return User{}, err
	}
	if err := rows.Close(); err != nil {
		return User{}, err
	}

	return userFromRow(user_row), nil
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

func (db *Database) CheckPassword(userAuth UserAuth) (int64, error) {

	password := []byte(userAuth.Password)

	var storedHash []byte
	var id int64
	row := db.Client.QueryRow(`SELECT password, id FROM user WHERE handle = ?`, userAuth.Handle)
	err := row.Scan(&storedHash, &id)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(storedHash, password)

	return id, err
}

func (db *Database) UpdateDisplayPicture(id int64, imagePath string) error {
	_, err := db.Client.Exec(`UPDATE user SET display_picture = ? WHERE id = ?`, "http://localhost:8080/"+imagePath, id)

	return err
}
