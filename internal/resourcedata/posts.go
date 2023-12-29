package resourcedata

import (
	"log"
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
		&postRow.MainImage,
		&postRow.CreatedAt,
		&postRow.UpdatedAt,
		&postRow.DeletedAt,
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
	`, postRow.AuthorID, postRow.Title, postRow.Slug, postRow.Subtitle, postRow.Content, postRow.MainImage, postRow.CreatedAt)

	return err
}

func (svc *Service) UpdatePost(post Post) error {
	postRow := post.GetRow()

	rows, err := svc.DbConn.NamedQuery(
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
