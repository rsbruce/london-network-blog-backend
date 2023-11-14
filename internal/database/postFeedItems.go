package database

import (
	"database/sql"
	"fmt"

	"rsbruce/blogsite-api/internal/postFeedItem"
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

func feedItemsPostsFromRows(rows []FeedItemPostRow) []postFeedItem.FeedItemPost {
	numItems := len(rows)
	items := make([]postFeedItem.FeedItemPost, numItems)

	for index, row := range rows {
		items[index] = postFeedItem.FeedItemPost{
			Title:      row.Title.String,
			Subtitle:   row.Subtitle.String,
			Slug:       row.Slug.String,
			Created_at: row.Created_at.String,
		}
	}

	return items
}

func feedItemsFromRows(rows []FeedItemRow) []postFeedItem.FeedItem {
	numItems := len(rows)
	items := make([]postFeedItem.FeedItem, numItems)

	for index, row := range rows {
		feedItemPost := postFeedItem.FeedItemPost{
			Title:      row.Post.Title.String,
			Subtitle:   row.Post.Subtitle.String,
			Slug:       row.Post.Slug.String,
			Created_at: row.Post.Created_at.String,
		}
		feedItemUser := postFeedItem.FeedItemUser{
			Display_picture: row.User.Display_picture.String,
			Display_name:    row.User.Display_name.String,
			Handle:          row.User.Handle.String,
		}
		items[index] = postFeedItem.FeedItem{
			Post: feedItemPost,
			User: feedItemUser,
		}
	}

	return items
}

func (db *Database) GetLatestPostFeed() ([]postFeedItem.FeedItem, error) {
	var feedItemRows []FeedItemRow

	rows, err := db.Client.Query(
		`SELECT post.title, post.subtitle, post.created_at, post.slug, user.display_name, user.display_picture, user.handle 
                FROM post 
                INNER JOIN user on post.author_id = user.id
                ORDER BY created_at DESC
                LIMIT 5;`)
	if err != nil {
		return nil, fmt.Errorf("latestPosts: %v", err)
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
		return nil, fmt.Errorf("latestPosts %v", err)
	}

	return feedItemsFromRows(feedItemRows), nil
}

func (db *Database) GetFeedItemPostsForAuthor(handle string) ([]postFeedItem.FeedItemPost, error) {
	var feedItemPosts []FeedItemPostRow
	rows, err := db.Client.Query(
		`SELECT post.title, post.subtitle, post.slug, post.created_at
        FROM post
        INNER JOIN user
        ON post.author_id = user.id
        WHERE user.handle = ?`, handle)
	if err != nil {
		return nil, fmt.Errorf("postItemsByUserHandle %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		feedItemPost := FeedItemPostRow{}
		err = rows.Scan(&feedItemPost.Title, &feedItemPost.Subtitle, &feedItemPost.Slug, &feedItemPost.Created_at)
		if err != nil {
			return nil, fmt.Errorf("postItemsByUserHandle %v", err)
		}
		feedItemPosts = append(feedItemPosts, feedItemPost)
	}

	return feedItemsPostsFromRows(feedItemPosts), nil
}
