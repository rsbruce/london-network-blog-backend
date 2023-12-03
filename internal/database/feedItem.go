package database

import (
	"database/sql"
	"fmt"
)

type FeedItemRow struct {
	Post FeedItemPostRow
	User FeedItemUserRow
}

type FeedItemPostRow struct {
	ID         sql.NullInt64
	Title      sql.NullString
	Subtitle   sql.NullString
	Slug       sql.NullString
	Created_at sql.NullString
	Deleted_at sql.NullString
}

type FeedItemUserRow struct {
	Display_name    sql.NullString
	Display_picture sql.NullString
	Handle          sql.NullString
}

type FeedItem struct {
	Post FeedItemPost `json:"post"`
	User FeedItemUser `json:"user"`
}

type FeedItemPost struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Slug       string `json:"slug"`
	Created_at string `json:"created_at"`
	Deleted_at string `json:"deleted_at"`
}

type FeedItemUser struct {
	Display_name    string `json:"display_name"`
	Display_picture string `json:"display_picture"`
	Handle          string `json:"handle"`
}

func feedItemsPostsFromRows(rows []FeedItemPostRow) []FeedItemPost {
	numItems := len(rows)
	items := make([]FeedItemPost, numItems)

	for index, row := range rows {
		items[index] = FeedItemPost{
			ID:         row.ID.Int64,
			Title:      row.Title.String,
			Subtitle:   row.Subtitle.String,
			Slug:       row.Slug.String,
			Created_at: row.Created_at.String,
			Deleted_at: row.Deleted_at.String,
		}
	}

	return items
}

func feedItemsFromRows(rows []FeedItemRow) []FeedItem {
	numItems := len(rows)
	items := make([]FeedItem, numItems)

	for index, row := range rows {
		feedItemPost := FeedItemPost{
			ID:         row.Post.ID.Int64,
			Title:      row.Post.Title.String,
			Subtitle:   row.Post.Subtitle.String,
			Slug:       row.Post.Slug.String,
			Created_at: row.Post.Created_at.String,
			Deleted_at: row.Post.Deleted_at.String,
		}
		feedItemUser := FeedItemUser{
			Display_picture: row.User.Display_picture.String,
			Display_name:    row.User.Display_name.String,
			Handle:          row.User.Handle.String,
		}
		items[index] = FeedItem{
			Post: feedItemPost,
			User: feedItemUser,
		}
	}

	return items
}

func (db *Database) GetLatestPostFeed() ([]FeedItem, error) {
	var feedItemRows []FeedItemRow

	rows, err := db.Client.Query(
		`SELECT post.title, post.subtitle, post.created_at, post.slug, user.display_name, user.display_picture, user.handle 
                FROM post 
                INNER JOIN user on post.author_id = user.id
                ORDER BY created_at DESC
                LIMIT 5;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		feedItemRow := FeedItemRow{Post: FeedItemPostRow{}, User: FeedItemUserRow{}}
		if err := rows.Scan(
			&feedItemRow.Post.Title,
			&feedItemRow.Post.Subtitle,
			&feedItemRow.Post.Created_at,
			&feedItemRow.Post.Slug,
			&feedItemRow.User.Display_name,
			&feedItemRow.User.Display_picture,
			&feedItemRow.User.Handle); err != nil {
			return nil, fmt.Errorf("latestPosts %v", err)
		}
		feedItemRows = append(feedItemRows, feedItemRow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feedItemsFromRows(feedItemRows), nil
}

func (db *Database) GetActiveFeedItemPostsForAuthor(handle string) ([]FeedItemPost, error) {
	var feedItemPosts []FeedItemPostRow
	rows, err := db.Client.Query(
		`SELECT post.id, post.title, post.subtitle, post.slug, post.created_at
        FROM post
        INNER JOIN user
        ON post.author_id = user.id
        WHERE user.handle = ? AND deleted_at IS NULL
		ORDER BY created_at DESC`, handle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		feedItemPost := FeedItemPostRow{}
		err = rows.Scan(&feedItemPost.ID, &feedItemPost.Title, &feedItemPost.Subtitle, &feedItemPost.Slug, &feedItemPost.Created_at)
		if err != nil {
			return nil, err
		}
		feedItemPosts = append(feedItemPosts, feedItemPost)
	}

	return feedItemsPostsFromRows(feedItemPosts), nil
}

func (db *Database) GetAllFeedItemPostsForAuthor(handle string) ([]FeedItemPost, error) {
	var feedItemPosts []FeedItemPostRow
	rows, err := db.Client.Query(
		`SELECT post.id, post.title, post.subtitle, post.slug, post.created_at, post.deleted_at
        FROM post
        INNER JOIN user
        ON post.author_id = user.id
        WHERE user.handle = ?
		ORDER BY created_at DESC`, handle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		feedItemPost := FeedItemPostRow{}
		err = rows.Scan(&feedItemPost.ID, &feedItemPost.Title, &feedItemPost.Subtitle, &feedItemPost.Slug, &feedItemPost.Created_at, &feedItemPost.Deleted_at)
		if err != nil {
			return nil, err
		}
		feedItemPosts = append(feedItemPosts, feedItemPost)
	}

	return feedItemsPostsFromRows(feedItemPosts), nil
}
