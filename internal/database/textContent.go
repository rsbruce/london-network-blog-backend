package database

import (
	"database/sql"
	"fmt"

	"rsbruce/blogsite-api/internal/models"
)

type TextContentRow struct {
	Content sql.NullString
}

func textContentFromRow(tcr TextContentRow) models.TextContent {
	return models.TextContent{
		Content: tcr.Content.String,
	}
}

func (db *Database) GetTextContent(slug string) (models.TextContent, error) {

	row := db.Client.QueryRow("SELECT content FROM text_content WHERE slug = ?", slug)
	var textContentRow TextContentRow
	if err := row.Scan(&textContentRow.Content); err != nil {
		return models.TextContent{}, fmt.Errorf("getTextContent %v", err)
	}

	return textContentFromRow(textContentRow), nil
}
