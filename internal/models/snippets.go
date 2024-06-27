package models

import (
	"database/sql"
	"errors"
	"time"
)

//we defined a Snippet type to hold the data for indivdual snippets,

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// we define a snippetModel type which will wrap a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title,content,created,expires)
			VALUES (?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	//we can use the LastInsertId to get the id of our newly inserted record in the table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil

}
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id,title,content,created,expires from snippets where expires > UTC_TIMESTAMP() AND id=?`

	row := m.DB.QueryRow(stmt, id)
	var snippet Snippet
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNoRecord error
		// instead (we'll create this in a moment).
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return snippet, nil
}
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
