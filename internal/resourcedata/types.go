package resourcedata

import (
	"database/sql"
)

type Post struct {
	ID         int64  `json:"id,omitempty"`
	Author_id  int64  `json:"author_id,omitempty"`
	Title      string `json:"title,omitempty"`
	Subtitle   string `json:"subtitle,omitempty"`
	Slug       string `json:"slug,omitempty"`
	Content    string `json:"content,omitempty"`
	Main_image string `json:"main_image,omitempty"`
	Created_at string `json:"created_at,omitempty"`
	Updated_at string `json:"updated_at,omitempty"`
	Deleted_at string `json:"deleted_at,omitempty"`
}

type PostRow struct {
	ID         int64
	Author_id  int64
	Title      string
	Subtitle   sql.NullString
	Slug       string
	Content    string
	Main_image sql.NullString
	Created_at string
	Updated_at sql.NullString
	Deleted_at sql.NullString
}

func (pr *PostRow) GetPost() *Post {
	return &Post{
		ID:         pr.ID,
		Title:      pr.Title,
		Subtitle:   pr.Subtitle.String,
		Slug:       pr.Slug,
		Content:    pr.Content,
		Main_image: pr.Main_image.String,
		Created_at: pr.Created_at,
		Updated_at: pr.Updated_at.String,
		Deleted_at: pr.Deleted_at.String,
	}
}

func (p *Post) GetRow() PostRow {
	postRow := PostRow{
		ID:         p.ID,
		Author_id:  p.Author_id,
		Title:      p.Title,
		Slug:       p.Slug,
		Created_at: p.Created_at,
	}
	if p.Subtitle != "" {
		postRow.Subtitle = sql.NullString{String: p.Subtitle, Valid: true}
	}
	if p.Main_image != "" {
		postRow.Main_image = sql.NullString{String: p.Main_image, Valid: true}
	}
	if p.Updated_at != "" {
		postRow.Updated_at = sql.NullString{String: p.Updated_at, Valid: true}
	}
	if p.Deleted_at != "" {
		postRow.Deleted_at = sql.NullString{String: p.Deleted_at, Valid: true}
	}
	return postRow
}

type User struct {
	ID              int64  `json:"id,omitempty"`
	Display_name    string `json:"display_name,omitempty"`
	Display_picture string `json:"display_picture,omitempty"`
	Email           string `json:"email,omitempty"`
	Handle          string `json:"handle,omitempty"`
	Blurb           string `json:"blurb,omitempty"`
	User_role       string `json:"user_role,omitempty"`
	Created_at      string `json:"created_at,omitempty"`
	Updated_at      string `json:"updated_at,omitempty"`
	Deleted_at      string `json:"deleted_at,omitempty"`
}

type UserRow struct {
	ID              int64
	Display_name    sql.NullString
	Display_picture sql.NullString
	Handle          string
	Blurb           sql.NullString
	User_role       sql.NullString
	Created_at      string
	Updated_at      sql.NullString
	Deleted_at      sql.NullString
}

func (ur *UserRow) GetUser() *User {
	return &User{
		ID:              ur.ID,
		Display_name:    ur.Display_name.String,
		Display_picture: ur.Display_picture.String,
		Handle:          ur.Handle,
		Blurb:           ur.Blurb.String,
		User_role:       ur.User_role.String,
		Created_at:      ur.Created_at,
		Updated_at:      ur.Updated_at.String,
		Deleted_at:      ur.Deleted_at.String,
	}
}

func (u *User) GetRow() UserRow {
	userRow := UserRow{
		ID:         u.ID,
		Handle:     u.Handle,
		Created_at: u.Created_at,
	}
	if u.Display_name != "" {
		userRow.Display_name = sql.NullString{String: u.Display_name, Valid: true}
	}
	if u.Display_picture != "" {
		userRow.Display_picture = sql.NullString{String: u.Display_picture, Valid: true}
	}
	if u.Blurb != "" {
		userRow.Blurb = sql.NullString{String: u.Blurb, Valid: true}
	}
	if u.Updated_at != "" {
		userRow.Updated_at = sql.NullString{String: u.Updated_at, Valid: true}
	}
	if u.Deleted_at != "" {
		userRow.Deleted_at = sql.NullString{String: u.Deleted_at, Valid: true}
	}
	return userRow
}

type FeedItem struct {
	Post `json:"post"`
	User `json:"user"`
}
