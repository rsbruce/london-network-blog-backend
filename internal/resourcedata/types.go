package resourcedata

import (
	"database/sql"
)

type Post struct {
	ID        int64  `json:"id,omitempty"`
	AuthorID        int64  `json:"author_id,omitempty"`
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
	ID        int64
	AuthorID        int64
	Title     string
	Subtitle  sql.NullString
	Slug      string
	Content   string
	MainImage sql.NullString
	CreatedAt string
	UpdatedAt sql.NullString
	DeletedAt sql.NullString
}

func (pr *PostRow) GetPost() *Post {
	return &Post{
		ID:        pr.ID,
		Title:     pr.Title,
		Subtitle:  pr.Subtitle.String,
		Slug:      pr.Slug,
		Content:   pr.Content,
		MainImage: pr.MainImage.String,
		CreatedAt: pr.CreatedAt,
		UpdatedAt: pr.UpdatedAt.String,
		DeletedAt: pr.DeletedAt.String,
	}
}

func (p *Post) GetRow() PostRow {
	postRow := PostRow{
		ID: p.ID,
		AuthorID: p.AuthorID,
		Title: p.Title,
		Slug: p.Slug,
		CreatedAt: p.CreatedAt,
	}
	if p.Subtitle != "" {
		postRow.Subtitle = sql.NullString{String: p.Subtitle, Valid: true}
	}
	if p.MainImage != "" {
		postRow.MainImage = sql.NullString{String: p.MainImage, Valid: true}
	}
	if p.UpdatedAt != "" {
		postRow.UpdatedAt = sql.NullString{String: p.UpdatedAt, Valid: true}
	}
	if p.DeletedAt != "" {
		postRow.DeletedAt = sql.NullString{String: p.DeletedAt, Valid: true}
	}
	return postRow
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
