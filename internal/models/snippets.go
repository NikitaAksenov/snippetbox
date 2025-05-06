package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := fmt.Sprintf(`
		INSERT INTO snippets (title, content, created, expires) VALUES (
			?,
			?,
			datetime('now'),
			datetime('now', '%d days')
		);`, expires)

	result, err := m.DB.Exec(stmt, title, content)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), err
}

const DateTime string = "2006-01-02 15:04:05"

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `
	SELECT id, title, content, created, expires
		FROM snippets
		WHERE id = ? AND datetime('now') < expires
	;`

	row := m.DB.QueryRow(query, id)

	snippet := Snippet{}

	var createdStr, expiresStr string
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &createdStr, &expiresStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
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

func (m *SnippetModel) Latest(count int) ([]*Snippet, error) {
	return nil, nil
}
