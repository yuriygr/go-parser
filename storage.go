package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// NewStorage - init new storage
func NewStorage() *Storage {
	db, err := sqlx.Connect("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalln(err)
	}
	db.SetConnMaxLifetime(time.Hour)
	return &Storage{db}
}

// Storage - Нечто такое обстрактное я хз
// попозже придумаю что тут написать так то
// штука крутая.
type Storage struct {
	db *sqlx.DB
}

type FileInsert struct {
	ID        int    `db:"v.id"`
	URL       string `db:"v.url"`
	Type      string `db:"v.type"`
	Filesize  int    `db:"v.filesize"`
	Md5       string `db:"v.md5"`
	Name      string `db:"v.name"`
	Board     string `db:"v.board"`
	Thread    string `db:"v.thread"`
	CreatedAt int64  `db:"v.created_at"`
}

// InsertFiles - Вставка данных че еще тут написать то
func (s *Storage) InsertFiles(f []*FileInsert) error {
	sqlStr := "INSERT INTO video (url, type, filesize, md5, name, board, thread, created_at) VALUES "

	for _, file := range f {
		sqlStr += fmt.Sprintf("('%s', '%s', '%d', '%s', '%s', '%s', '%s', '%d'),", file.URL, file.Type, file.Filesize, file.Md5, file.Name, file.Board, file.Thread, file.CreatedAt)
	}

	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	// update duplicate
	sqlStr += " ON DUPLICATE KEY UPDATE url = url"

	//format all vals at once
	_, err := s.db.Exec(sqlStr)

	return err
}

// GetFiles - Получает файлы
func (s *Storage) GetFiles() ([]*FileInsert, error) {
	files := []*FileInsert{}
	sql := "select v.* from video as v"

	err := s.db.Select(&files, sql)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// DeleteFile - Удаляет файл из базы данных
func (s *Storage) DeleteFile(id int) error {
	sql := "DELETE FROM video WHERE id = ?"

	_, err := s.db.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}
