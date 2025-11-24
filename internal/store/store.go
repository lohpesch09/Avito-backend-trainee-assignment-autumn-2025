package store

import (
	"database/sql"
)

type Store struct {
	DB *sql.DB
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) BeginTx() (*sql.Tx, error) {
	return s.DB.Begin()
}
