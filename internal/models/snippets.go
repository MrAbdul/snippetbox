package models

import (
	"database/sql"
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
	return Snippet{}, nil
}
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
