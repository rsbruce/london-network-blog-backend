package database

import (
	"database/sql"
	"fmt"
	"rsbruce/blogsite-api/internal/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

type PostRow struct {
	ID         []byte
	Author_id  []byte
	Slug       sql.NullString
	Title      sql.NullString
	Subtitle   sql.NullString
	Content    sql.NullString
	Main_image sql.NullString
	Created_at sql.NullString
}

func postFromRow(row PostRow) models.Post {
	return models.Post{
		ID:         row.ID,
		Author_id:  row.Author_id,
		Slug:       row.Slug.String,
		Title:      row.Title.String,
		Subtitle:   row.Subtitle.String,
		Content:    row.Content.String,
		Main_image: row.Main_image.String,
		Created_at: row.Created_at.String,
	}
}

func partialUserFromRow(row UserRow) models.User {
	return models.User{
		Handle:          row.Handle.String,
		Display_name:    row.Display_name.String,
		Display_picture: row.Display_picture.String,
	}
}

func (db *Database) NewPost(post models.Post) (models.Post, error) {
	post.ID = uuid.NewV4().Bytes()
	postRow := PostRow{
		ID:         post.ID,
		Author_id:  post.Author_id,
		Slug:       sql.NullString{String: post.Slug, Valid: true},
		Title:      sql.NullString{String: post.Title, Valid: true},
		Subtitle:   sql.NullString{String: post.Subtitle, Valid: true},
		Content:    sql.NullString{String: post.Content, Valid: true},
		Main_image: sql.NullString{String: post.Main_image, Valid: true},
		Created_at: sql.NullString{String: time.Now().Format(time.DateOnly), Valid: true},
	}

	rows, err := db.Client.NamedQuery(
		`INSERT INTO post 
		(id, author_id, slug, title, subtitle, content, main_image, created_at) VALUES
		(:id, :author_id, :slug, :title, :subtitle, :content, :main_image, :created_at) `,
		postRow,
	)
	if err != nil {
		return models.Post{}, fmt.Errorf("failed to insert post: %w", err)
	}
	if err = rows.Close(); err != nil {
		return models.Post{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return post, nil
}

func postPageFromRow(post_row PostRow, user_row UserRow) models.PostPage {
	post := postFromRow(post_row)
	user := userFromRow(user_row)

	return models.PostPage{
		Post: post,
		User: user,
	}
}

func (db *Database) GetPostWithUser(slug string) (models.PostPage, error) {
	row := db.Client.QueryRow(
		`SELECT post.id, post.author_id, post.title, post.subtitle, post.content, post.slug, post.main_image, post.created_at, user.display_name, user.display_picture, user.handle 
        FROM post 
        INNER JOIN user on post.author_id = user.id 
        WHERE post.slug = ?`, slug)
	var post_row PostRow
	var user_row UserRow

	if err := row.Scan(
		&post_row.ID,
		&post_row.Author_id,
		&post_row.Title,
		&post_row.Subtitle,
		&post_row.Content,
		&post_row.Slug,
		&post_row.Main_image,
		&post_row.Created_at,
		&user_row.Display_name,
		&user_row.Display_picture,
		&user_row.Handle,
	); err != nil {
		return models.PostPage{}, fmt.Errorf("getPost %v", err)
	}

	return postPageFromRow(post_row, user_row), nil
}
