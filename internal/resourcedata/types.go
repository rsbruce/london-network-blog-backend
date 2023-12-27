package resourcedata

import (
	"database/sql"
)

type Post struct {
	ID        int64  `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Subtitle  string `json:"subtitle,omitempty"`
	Slug      string `json:"slug,omitempty"`
	Content   string `json:"content,omitempty"`
	MainImage string `json:"main_image,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
}

type PostRow struct {
	ID        sql.NullInt64
	Title     sql.NullString
	Subtitle  sql.NullString
	Slug      sql.NullString
	Content   sql.NullString
	MainImage sql.NullString
	CreatedAt sql.NullString
	UpdatedAt sql.NullString
	DeletedAt sql.NullString
}

func (pr *PostRow) GetPost() *Post {
	return &Post{
		ID:        pr.ID.Int64,
		Title:     pr.Title.String,
		Subtitle:  pr.Subtitle.String,
		Slug:      pr.Slug.String,
		Content:   pr.Content.String,
		MainImage: pr.MainImage.String,
		CreatedAt: pr.CreatedAt.String,
		UpdatedAt: pr.UpdatedAt.String,
		DeletedAt: pr.DeletedAt.String,
	}
}

type User struct {
	ID             int64  `json:"id,omitempty"`
	DisplayName    string `json:"display_name,omitempty"`
	DisplayPicture string `json:"display_picture,omitempty"`
	Handle         string `json:"handle,omitempty"`
	Blurb          string `json:"blurb,omitempty"`
	UserRole       string `json:"user_role,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	DeletedAt      string `json:"deleted_at,omitempty"`
}

type UserRow struct {
	ID             sql.NullInt64
	DisplayName    sql.NullString
	DisplayPicture sql.NullString
	Handle         sql.NullString
	Blurb          sql.NullString
	UserRole       sql.NullString
	CreatedAt      sql.NullString
	UpdatedAt      sql.NullString
	DeletedAt      sql.NullString
}

func (ur *UserRow) GetUser() *User {
	return &User{
		ID:             ur.ID.Int64,
		DisplayName:    ur.DisplayName.String,
		DisplayPicture: ur.DisplayPicture.String,
		Handle:         ur.Handle.String,
		Blurb:          ur.Blurb.String,
		UserRole:       ur.UserRole.String,
		CreatedAt:      ur.CreatedAt.String,
		UpdatedAt:      ur.UpdatedAt.String,
		DeletedAt:      ur.DeletedAt.String,
	}
}

type FeedItem struct {
	Post `json:"post"`
	User `json:"user"`
}
