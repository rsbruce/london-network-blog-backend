package database

import (
	"database/sql"
	"fmt"

	"rsbruce/blogsite-api/internal/models"
)

type UserProfileRow struct {
	User        UserRow
	LatestPosts []FeedItemPostRow
}
type UserRow struct {
	ID              []byte
	Handle          sql.NullString
	Blurb           sql.NullString
	Display_name    sql.NullString
	Display_picture sql.NullString
	User_role       sql.NullString
	Created_at      sql.NullString
}

func userFromRow(row UserRow) models.User {
	return models.User{
		ID:              row.ID,
		Handle:          row.Handle.String,
		Blurb:           row.Blurb.String,
		Display_name:    row.Display_name.String,
		Display_picture: row.Display_picture.String,
		User_role:       row.User_role.String,
		Created_at:      row.Created_at.String,
	}
}

func (db *Database) GetUser(handle string) (models.User, error) {

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
		return models.User{}, fmt.Errorf("getUser %v", err)
	}

	return userFromRow(user_row), nil
}

func (db *Database) UpdateUser(handle string, user models.User) (models.User, error) {

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
		return models.User{}, fmt.Errorf("UpdateUser %v", err)
	}
	if err := rows.Close(); err != nil {
		return models.User{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return userFromRow(user_row), nil
}
