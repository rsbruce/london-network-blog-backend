package resourcedata

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (svc *Service) GetPost(authorHandle string, slug string) (*Post, error) {
	var postRow PostRow

	err := svc.DbConn.QueryRow(`
		SELECT post.id, post.title, post.subtitle, post.content, post.main_image, post.created_at, post.updated_at, post.deleted_at 
		FROM post
		JOIN user
		WHERE user.handle = ?
		AND post.author_id = user.id
		AND post.slug = ?
	`, authorHandle, slug).Scan(
		&postRow.ID,
		&postRow.Title,
		&postRow.Subtitle,
		&postRow.Content,
		&postRow.Main_image,
		&postRow.Created_at,
		&postRow.Updated_at,
		&postRow.Deleted_at,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return postRow.GetPost(), nil
}

func (svc *Service) CreatePost(post Post) error {
	postRow := post.GetRow()

	_, err := svc.DbConn.Exec(`
		INSERT INTO post 
		(author_id, title, slug, subtitle, content, main_image, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, postRow.Author_id, postRow.Title, postRow.Slug, postRow.Subtitle, postRow.Content, postRow.Main_image, postRow.Created_at)

	return err
}

func (svc *Service) UpdatePost(post Post, oldSlug string) error {
	fields := getFieldsToUpdate(post, []string{"id", "author_id", "created_at", "deleted_at"})
	assignments := make([]string, len(fields))
	for i, field := range fields {
		assignments[i] = fmt.Sprintf("%s = :%s", field, field)
	}

	query := `UPDATE post SET ` +
		strings.Join(assignments, ", ") +
		fmt.Sprintf(` WHERE author_id = :author_id AND slug = "%s"`, oldSlug)

	rows, err := svc.DbConn.NamedQuery(query, post)
	if err != nil {
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}

	return nil
}

func (svc *Service) ArchivePost(authorID int64, slug string) error {
	_, err := svc.DbConn.Exec(`UPDATE post
		SET deleted_at = ?
		WHERE author_id = ? AND slug = ?`,
		time.Now().Format(time.DateTime), authorID, slug)

	return err
}

func (svc *Service) RestorePost(authorID int64, slug string) error {
	_, err := svc.DbConn.Exec(`UPDATE post
		SET deleted_at = NULL
		WHERE author_id = ? AND slug = ?`, authorID, slug)

	return err
}

func (svc *Service) DeletePost(authorID int64, slug string) error {
	_, err := svc.DbConn.Exec(`DELETE FROM post WHERE author_id = ? AND slug = ?`, authorID, slug)

	return err
}
