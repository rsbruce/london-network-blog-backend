package database

import (
	"database/sql"
	"time"
)

type PostRow struct {
	ID         sql.NullInt64
	Author_id  sql.NullInt64
	Slug       sql.NullString
	Title      sql.NullString
	Subtitle   sql.NullString
	Content    sql.NullString
	Main_image sql.NullString
	Created_at sql.NullString
	Updated_at sql.NullString
	Deleted_at sql.NullString
}

type Post struct {
	ID         int64  `json:"id"`
	Author_id  int64  `json:"author_id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Content    string `json:"content"`
	Main_image string `json:"main_image"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Deleted_at string `json:"deleted_at"`
}

type PostPage struct {
	Post Post `json:"post"`
	User User `json:"user"`
}

func postFromRow(row PostRow) Post {
	return Post{
		ID:         row.ID.Int64,
		Author_id:  row.Author_id.Int64,
		Slug:       row.Slug.String,
		Title:      row.Title.String,
		Subtitle:   row.Subtitle.String,
		Content:    row.Content.String,
		Main_image: row.Main_image.String,
		Created_at: row.Created_at.String,
	}
}

func partialUserFromRow(row UserRow) User {
	return User{
		Handle:          row.Handle.String,
		Display_name:    row.Display_name.String,
		Display_picture: row.Display_picture.String,
	}
}
func (db *Database) UpdatePost(post Post) (Post, error) {
	postRow := PostRow{
		ID:         sql.NullInt64{Int64: post.ID, Valid: post.ID != 0},
		Slug:       sql.NullString{String: post.Slug, Valid: post.Slug != ""},
		Title:      sql.NullString{String: post.Title, Valid: post.Title != ""},
		Subtitle:   sql.NullString{String: post.Subtitle, Valid: post.Subtitle != ""},
		Content:    sql.NullString{String: post.Content, Valid: post.Content != ""},
		Main_image: sql.NullString{String: post.Main_image, Valid: post.Main_image != ""},
		Updated_at: sql.NullString{String: time.Now().Format(time.DateTime), Valid: true},
	}

	rows, err := db.Client.NamedQuery(
		`UPDATE post SET
		slug = :slug,
		title = :title,
		subtitle = :subtitle,
		content = :content,
		main_image = :main_image,
		updated_at = :updated_at
		WHERE id = :id`,
		postRow,
	)

	if err != nil {
		return Post{}, err
	}
	if err := rows.Close(); err != nil {
		return Post{}, err
	}

	return postFromRow(postRow), nil
}

func (db *Database) NewPost(post Post) (Post, error) {
	postRow := PostRow{
		Author_id:  sql.NullInt64{Int64: post.Author_id, Valid: true},
		Slug:       sql.NullString{String: post.Slug, Valid: true},
		Title:      sql.NullString{String: post.Title, Valid: true},
		Subtitle:   sql.NullString{String: post.Subtitle, Valid: true},
		Content:    sql.NullString{String: post.Content, Valid: true},
		Main_image: sql.NullString{String: post.Main_image, Valid: true},
		Created_at: sql.NullString{String: time.Now().Format(time.DateTime), Valid: true},
	}

	res, err := db.Client.Exec(`INSERT INTO post 
	(author_id, slug, title, subtitle, content, main_image, created_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`,
		postRow.Author_id, postRow.Slug, postRow.Title, postRow.Subtitle, postRow.Content, postRow.Main_image, postRow.Created_at)
	if err != nil {
		return Post{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Post{}, err
	}

	postRow.ID = sql.NullInt64{Int64: id, Valid: true}

	return postFromRow(postRow), nil
}

func postPageFromRow(post_row PostRow, user_row UserRow) PostPage {
	post := postFromRow(post_row)
	user := userFromRow(user_row)

	return PostPage{
		Post: post,
		User: user,
	}
}

func (db *Database) GetPostWithUser(slug string) (PostPage, error) {
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
		return PostPage{}, err
	}

	return postPageFromRow(post_row, user_row), nil
}
