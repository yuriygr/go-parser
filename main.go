package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"github.com/yuriygr/go-loggy"
	"github.com/yuriygr/go-parser/boards"
	"github.com/yuriygr/go-parser/boards/dvach"
	"github.com/yuriygr/go-parser/container"
	"github.com/yuriygr/go-parser/repository"
)

// Каналы для обработки начальных данных закупок и обновления закупок
var filesChain = make(chan boards.FileInsert)

// Канал логов
var infoChain = make(chan string)
var successChain = make(chan string)
var errorsChain = make(chan error)

func main() {
	container := container.NewContainer()
	repository := repository.NewRepository(container)
	dvach := dvach.BuildBoard(container)

	// Run tusk on startup
	go parseFiles(dvach, []string{"b", "po"})
	//go actualizeFiles(repository, container)

	// Use cron it's overkill, but
	// @link https://stackoverflow.com/a/28886762
	c := cron.New()
	c.AddFunc("@every 30m", func() {
		// Основная рабочая, идеальная горутина
		go parseFiles(dvach, []string{"b", "po", "mov", "tv", "a", "v", "cg", "vg", "kpop"})
	})
	c.AddFunc("@every 40m", func() {
		go actualizeFiles(repository, container)
	})
	c.Start()

	// Сохраняем в базу данных полученные файлы
	go saveFiles(repository)

	// Работаем с извинениями горутины
	handleChain(container.Logger)

	// Очищаем каналы
	<-filesChain
	<-infoChain
	<-successChain
	<-errorsChain
}

func parseFiles(imageboard boards.Board, bs []string) {
	// Уведомляем, что задача пошла
	infoChain <- fmt.Sprintf("[Parser] 30 minutes has passed, let's parse %s!", imageboard.Name)

	// Lets parse boards
	for _, b := range bs {

		board, _ := imageboard.GetBoard(boards.BoardRequest{Board: b})
		for _, thread := range board.Threads {

			posts, _ := imageboard.GetThread(boards.ThreadRequest{Board: b, ThreadNumber: thread.Num})
			for _, post := range posts {

				for _, file := range post.Files {
					if file.GetType() == "Video" {
						filesChain <- boards.FileInsert{
							URL:       fmt.Sprintf(imageboard.GetLink("file"), file.Path),
							Type:      file.GetMimeType(),
							Filesize:  file.Size,
							Width:     file.Width,
							Height:    file.Height,
							Md5:       file.Md5,
							Name:      file.RealName(),
							Board:     b,
							Thread:    thread.Num,
							CreatedAt: post.Timestamp,
						}
					}
				}

			}

		}

		infoChain <- fmt.Sprintf("[Parser] /%s/ parsed, let's parse other board", b)
	}

	infoChain <- fmt.Sprint("[Parser] All boards parsed, see you soon")
}

// actualizeFiles - Обновляет ссылки, проверяет, нет ли удаленных
func actualizeFiles(repository *repository.Repository, container *container.Container) {
	// Уведомляем, что задача пошла
	infoChain <- fmt.Sprint("[Actualizer] 40 minutes has passed, let's truncate some links")

	files, err := repository.GetFiles()
	if err != nil {
		errorsChain <- fmt.Errorf("[Actualizer] Get files error (%v)", err)
	}

	for _, file := range files {
		if err := container.Client.IsExist(file.URL); err != nil {
			if err := repository.DeleteFile(file.ID); err != nil {
				errorsChain <- fmt.Errorf("[Actualizer] Delete file error (%v)", err)
			} else {
				successChain <- fmt.Sprintf("[Actualizer] File %s gone", file.URL)
			}
		}
	}

	infoChain <- fmt.Sprint("[Actualizer] All files have been updated, see you in 40 minutes")
}

// saveFiles - Сохраняем в базу данных полученные файлы
func saveFiles(repository *repository.Repository) {
	for {
		select {
		case file := <-filesChain:
			err := repository.CreateFile(file)
			if err != nil {
				errorsChain <- fmt.Errorf("Файл %s не добавлен (%v)", file.URL, err)
				continue
			}
			successChain <- fmt.Sprintf("Файл %s добавлен", file.URL)
		}
	}
}

// handleChain - Обрабатывает результат выполнения горутин-логгеров
func handleChain(logger *loggy.Logger) {
	for {
		select {
		case entry := <-successChain:
			logger.Success(entry)
		case entry := <-errorsChain:
			logger.Error(entry)
		case entry := <-infoChain:
			logger.Info(entry)
		}
	}
}
