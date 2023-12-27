package resourcedata

import "log"

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