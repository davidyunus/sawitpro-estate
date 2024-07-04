package helper

import (
	"database/sql"
)

type key string

type Manager struct {
	db  *sql.DB
	key key
}

func (m *Manager) GetKey() key {
	return m.key
}

func NewManager(db *sql.DB, k key) *Manager {
	return &Manager{
		db:  db,
		key: k,
	}
}
