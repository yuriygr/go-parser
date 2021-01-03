package services

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// NewStorage - Создаем экземпляр хранилища
func NewStorage(config *Config) *Storage {
	db, err := sqlx.Connect("mysql", config.Storage.DSN)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetConnMaxLifetime(time.Hour)

	// Unsafe becouse i sleep
	return &Storage{db.Unsafe()}
}

// Storage - Нечто такое обстрактное я хз
// попозже придумаю что тут написать так то
// штука крутая.
type Storage struct {
	db *sqlx.DB
}

// Query - Потому что блин
func (s *Storage) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

// NamedQuery - Потому что блин
func (s *Storage) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return s.db.NamedQuery(query, arg)
}

// Select - Просто потому что могу
func (s *Storage) Select(dest interface{}, query string, args ...interface{}) error {
	return s.db.Select(dest, query, args...)
}

// Get - Потому что блин
func (s *Storage) Get(dest interface{}, query string, args ...interface{}) error {
	return s.db.Get(dest, query, args...)
}

// Exec - Потому что блин
func (s *Storage) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

// NamedExec - Потому что блин
func (s *Storage) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return s.db.NamedExec(query, arg)
}

// PrepareNamed - Потому что блин
func (s *Storage) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed(query)
}

// Preparex - Потому что блин
func (s *Storage) Preparex(query string) (*sqlx.Stmt, error) {
	return s.db.Preparex(query)
}

// Queryx - Потому что блядь да
func (s *Storage) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return s.db.Queryx(query, args)
}
