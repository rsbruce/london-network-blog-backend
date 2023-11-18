package database

import (
	"database/sql"
	"fmt"
	"rsbruce/blogsite-api/internal/models"
)

type SlugRow struct {
	Slug sql.NullString
}

func slugFromRows(rows []SlugRow) models.Slugs {

	numItems := len(rows)
	slugs_slice := make([]string, numItems)
	slugs_obj := models.Slugs{Slugs: slugs_slice}

	for index, row := range rows {
		slugs_obj.Slugs[index] = row.Slug.String
	}

	return slugs_obj
}

func (db *Database) GetSlugsForUser(handle string) (models.Slugs, error) {
	rows, err := db.Client.Query(
		`SELECT post.slug 
        FROM post 
        INNER JOIN user on post.author_id = user.id 
        WHERE user.handle = ?`, handle)
	var slug_rows []SlugRow

	if err != nil {
		return models.Slugs{}, fmt.Errorf("GetSlugsForUser %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		slug_row := SlugRow{}

		if err := rows.Scan(
			&slug_row.Slug,
		); err != nil {
			return models.Slugs{}, fmt.Errorf("getPost %v", err)
		}
		slug_rows = append(slug_rows, slug_row)
	}

	return slugFromRows(slug_rows), nil
}
