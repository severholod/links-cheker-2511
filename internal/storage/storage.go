package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

var (
	ErrLinksNotFound = errors.New("links not found")
)

func NewStorage(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening storage: %v", err)
	}

	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS links (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	urls TEXT[] NOT NULL,
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

func (s *Storage) SaveUrls(urls []string) (int64, error) {
	stmt, err := s.db.Prepare(`INSERT INTO links (urls) VALUES (?)`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()
	res, err := stmt.Exec(urls)
	if err != nil {
		return 0, fmt.Errorf("failed to save urls: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to save urls: %w", err)
	}
	return id, nil
}
func (s *Storage) GetUrls(id int64) ([]string, error) {
	stmt, err := s.db.Prepare(`SELECT urls FROM links WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()
	var res []string
	err = stmt.QueryRow(id).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLinksNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get urls: %w", err)
	}

	return res, nil
}
