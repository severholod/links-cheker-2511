package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Storage struct {
	db *sql.DB
}

type Links struct {
	URL    string
	Status string
}

var (
	ErrLinksNotFound = errors.New("links not found")
)

func NewStorage(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening storage: %v", err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			set_id INTEGER,
			url TEXT,
			status TEXT,
			checked_at DATETIME,
			FOREIGN KEY(set_id) REFERENCES link_sets(id)
		);
	)`)
	if err != nil {
		return nil, fmt.Errorf("error preparing links table: %w", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("error creating links table: %w", err)
	}
	return &Storage{db}, nil
}

func (s *Storage) SaveUrl(id int, url string, status string, checked_at time.Time) error {
	stmt, err := s.db.Prepare(`
		INSERT INTO links (set_id, url, status, checked_at) VALUES (?, ?, ?, ?);
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()
	_, err = stmt.Exec(id, url, status, checked_at)
	if err != nil {
		return fmt.Errorf("failed to save urls: %w", err)
	}

	return nil
}
func (s *Storage) GetUrls(id int) ([]Links, error) {
	stmt, err := s.db.Prepare(`SELECT url, status FROM links WHERE set_id = ? ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var res []Links
	for rows.Next() {
		var link Links
		err := rows.Scan(&link.URL, &link.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		res = append(res, link)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	if len(res) == 0 {
		return nil, ErrLinksNotFound
	}

	return res, nil
}
