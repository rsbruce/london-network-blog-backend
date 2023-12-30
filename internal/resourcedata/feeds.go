package resourcedata

import (
	"log"
)

func (svc *Service) GetFeed(limit int) ([]FeedItem, error) {
	var feed []FeedItem

	if limit == 0 {
		limit = 10
	}

	query := `
		SELECT post.id, post.title, post.subtitle, post.created_at, post.updated_at, user.display_name, user.display_picture, user.handle
		FROM post
		JOIN user ON post.author_id = user.id
		WHERE post.deleted_at IS NULL
		LIMIT ?
		`
	rows, err := svc.DbConn.Query(query, limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var postRow PostRow
		var userRow UserRow
		rows.Scan(
			&postRow.ID,
			&postRow.Title,
			&postRow.Subtitle,
			&postRow.Created_at,
			&postRow.Updated_at,
			&userRow.Display_name,
			&userRow.Display_picture,
			&userRow.Handle,
		)
		feed = append(feed, FeedItem{Post: *postRow.GetPost(), User: *userRow.GetUser()})
	}

	return feed, nil
}

func (svc *Service) GetPersonalFeed(userId int64, limit int) ([]Post, error) {
	if limit == 0 {
		limit = 10
	}

	var feed []Post

	rows, err := svc.DbConn.Query(`
		SELECT id, title, subtitle, created_at, updated_at, deleted_at
		FROM post
		WHERE deleted_at IS NULL AND author_id = ?
		LIMIT ?
	`, userId, limit)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var postRow PostRow
		err = rows.Scan(&postRow.ID, &postRow.Title, &postRow.Subtitle, &postRow.Created_at, &postRow.Updated_at, &postRow.Deleted_at)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		feed = append(feed, *postRow.GetPost())
	}

	return feed, nil

}

func (svc *Service) GetSingleUserFeed(handle string, limit int) ([]Post, error) {
	if limit == 0 {
		limit = 10
	}

	var feed []Post

	rows, err := svc.DbConn.Query(`
		SELECT post.id, post.title, post.subtitle, post.created_at, post.updated_at
		FROM post
		JOIN user on user.id = post.author_id
		WHERE deleted_at IS NULL AND handle = ?
		LIMIT ?
	`, handle, limit)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var postRow PostRow
		err = rows.Scan(&postRow.ID, &postRow.Title, &postRow.Subtitle, &postRow.Created_at, &postRow.Updated_at)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		feed = append(feed, *postRow.GetPost())
	}

	return feed, nil

}
