package database

import (
	"database/sql"
	"fmt"
)

type TextContentRow struct {
	Content sql.NullString
}

type TextContent struct {
	Content string `json:"content"`
}

func textContentFromRow(tcr TextContentRow) TextContent {
	return TextContent{
		Content: tcr.Content.String,
	}
}

func (db *Database) GetTextContent(slug string) (TextContent, error) {

	row := db.Client.QueryRow("SELECT content FROM text_content WHERE slug = ?", slug)
	var textContentRow TextContentRow
	if err := row.Scan(&textContentRow.Content); err != nil {
		return TextContent{}, fmt.Errorf("getTextContent %v", err)
	}

	return textContentFromRow(textContentRow), nil
}
