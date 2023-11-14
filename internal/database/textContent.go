package database

import (
	"database/sql"
	"fmt"

	"rsbruce/blogsite-api/internal/textContent"
)

type TextContentRow struct {
	Content sql.NullString
}

func textContentFromRow(tcr TextContentRow) textContent.TextContent {
	return textContent.TextContent{
		Content: tcr.Content.String,
	}
}

func (db *Database) GetTextContent(slug string) (textContent.TextContent, error) {

	row := db.Client.QueryRow("SELECT content FROM text_content WHERE slug = \"about\"")
	var textContentRow TextContentRow
	if err := row.Scan(&textContentRow.Content); err != nil {
		return textContent.TextContent{}, fmt.Errorf("getTextContent %v", err)
	}

	return textContentFromRow(textContentRow), nil
}
