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

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `
	SELECT id, title, content, created, expires
		FROM snippets
		WHERE id = ? AND datetime('now') < expires
	;`

	row := m.DB.QueryRow(query, id)

	snippet, err := ParseRowToSnippet(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

func (m *SnippetModel) Latest(count int) ([]*Snippet, error) {
	query := `
	SELECT id, title, content, created, expires
		FROM snippets
		WHERE datetime('now') < expires
		ORDER BY id DESC
		LIMIT ?
	;`

	rows, err := m.DB.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	latestSnippets := make([]*Snippet, 0, 10)

	for rows.Next() {
		snippet, err := ParseRowsToSnippet(rows)
		if err != nil {
			return nil, err
		}

		latestSnippets = append(latestSnippets, snippet)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return latestSnippets, nil
}
