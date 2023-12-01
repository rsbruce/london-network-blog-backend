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
	Title      sql.NullString
	Subtitle   sql.NullString
	Slug       sql.NullString
	Created_at sql.NullString
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
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Slug       string `json:"slug"`
	Created_at string `json:"created_at"`
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
			Title:      row.Title.String,
			Subtitle:   row.Subtitle.String,
			Slug:       row.Slug.String,
			Created_at: row.Created_at.String,
		}
	}

	return items
}

func feedItemsFromRows(rows []FeedItemRow) []FeedItem {
	numItems := len(rows)
	items := make([]FeedItem, numItems)

	for index, row := range rows {
		feedItemPost := FeedItemPost{
			Title:      row.Post.Title.String,
			Subtitle:   row.Post.Subtitle.String,
			Slug:       row.Post.Slug.String,
			Created_at: row.Post.Created_at.String,
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

func (db *Database) GetFeedItemPostsForAuthor(handle string) ([]FeedItemPost, error) {
	var feedItemPosts []FeedItemPostRow
	rows, err := db.Client.Query(
		`SELECT post.title, post.subtitle, post.slug, post.created_at
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
		err = rows.Scan(&feedItemPost.Title, &feedItemPost.Subtitle, &feedItemPost.Slug, &feedItemPost.Created_at)
		if err != nil {
			return nil, err
		}
		feedItemPosts = append(feedItemPosts, feedItemPost)
	}

	return feedItemsPostsFromRows(feedItemPosts), nil
}
