package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yuriygr/go-parser/boards"
	"github.com/yuriygr/go-parser/container"
)

// NewRepository - Создаем новый парсер
func NewRepository(container *container.Container) *Repository {
	return &Repository{container.Storage}
}

// Repository - Его величество Репозиторий
type Repository struct {
	storage *sqlx.DB
}

// CreateFile - Вставка данных че еще тут написать то
func (r *Repository) CreateFile(file boards.FileInsert) (err error) {
	_, err = r.storage.NamedExec(createFile, file)
	return err
}

// InsertFiles - Вставка данных че еще тут написать то
func (r *Repository) InsertFiles(f []*boards.FileInsert) error {
	sqlStr := "INSERT INTO video (url, type, filesize, width, height, md5, name, board, thread, created_at) VALUES "

	for _, file := range f {
		sqlStr += fmt.Sprintf("('%s', '%s', '%d', '%d', '%d', '%s', '%s', '%s', '%s', '%d'),", file.URL, file.Type, file.Filesize, file.Width, file.Height, file.Md5, file.Name, file.Board, file.Thread, file.CreatedAt)
	}

	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	// update duplicate
	sqlStr += " ON DUPLICATE KEY UPDATE url = url"

	//format all vals at once
	_, err := r.storage.Exec(sqlStr)

	return err
}

// GetFiles - Получает файлы
func (r *Repository) GetFiles() ([]*boards.FileInsert, error) {
	files := []*boards.FileInsert{}

	err := r.storage.Select(&files, selectFiles)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// DeleteFile - Удаляет файл из базы данных
func (r *Repository) DeleteFile(id int) error {
	_, err := r.storage.Exec(deleteFileByID, id)
	if err != nil {
		return err
	}

	return nil
}
