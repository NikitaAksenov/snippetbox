package models

import (
	"database/sql"
	"time"
)

const DateTime string = "2006-01-02 15:04:05"

func ParseRowToSnippet(row *sql.Row) (*Snippet, error) {
	snippet := Snippet{}

	var createdStr, expiresStr string
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &createdStr, &expiresStr)
	if err != nil {
		return nil, err
	}

	snippet.Created, err = time.Parse(DateTime, createdStr)
	if err != nil {
		return nil, nil
	}

	snippet.Expires, err = time.Parse(DateTime, expiresStr)
	if err != nil {
		return nil, nil
	}

	return &snippet, nil
}

func ParseRowsToSnippet(rows *sql.Rows) (*Snippet, error) {
	snippet := Snippet{}

	var createdStr, expiresStr string
	err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &createdStr, &expiresStr)
	if err != nil {
		return nil, err
	}

	snippet.Created, err = time.Parse(DateTime, createdStr)
	if err != nil {
		return nil, nil
	}

	snippet.Expires, err = time.Parse(DateTime, expiresStr)
	if err != nil {
		return nil, nil
	}

	return &snippet, nil
}
