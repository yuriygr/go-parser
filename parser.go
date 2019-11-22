package main

import (
	"fmt"
	"log"
)

// BoardParserParams - параметры для парсинга
type BoardParserParams struct {
	boardSlug string
	fileType  string
}

// BoardParser - Отвечает за парс доски
type BoardParser struct {
	client     *ClientInstance
	storage    *Storage
	imageboard string
	boards     []string
}

// Run - метод, вызываемый кроном в Job
func (bp *BoardParser) Run() {
	log.Printf("===> 30 minutes has passed, let's parse %s!", bp.imageboard)

	// Lets parse boards
	for _, board := range bp.boards {

		// Prepare params
		params := &BoardParserParams{boardSlug: board, fileType: "Video"}

		// Run parser task
		files := bp.ParseFilesTask(params)
		if err := bp.storage.InsertFiles(files); err != nil {
			log.Println("Insert files", err)
		}

		log.Printf("Board %s/%s/ parsed successfully!", bp.imageboard, params.boardSlug)
	}

	log.Print("===> All boards parsed, see you soon!")
}

// ParseFilesTask - Функция парсинга разделов и поиска файлов
func (bp *BoardParser) ParseFilesTask(params *BoardParserParams) []*FileInsert {
	files := []*FileInsert{}

	threads := bp.ParseBoard(params.boardSlug).Threads
	for _, thread := range threads {

		posts := bp.ParseThread(params.boardSlug, thread.Num)
		for _, post := range posts {
			for _, file := range post.Files {
				if file.GetType() == params.fileType {
					preparedFile := &FileInsert{
						URL:       fmt.Sprintf(bp.GetLink("file"), file.Path),
						Type:      file.GetMimeType(),
						Filesize:  file.Size,
						Md5:       file.Md5,
						Name:      file.RealName(),
						Board:     params.boardSlug,
						Thread:    thread.Num,
						CreatedAt: post.Timestamp,
					}
					files = append(files, preparedFile)
				}
			}
		}

	}

	return files
}

/**
 * Base functions
 */

// ParseBoard - Парсит страницу раздела (компактную, на данный момент)
func (bp *BoardParser) ParseBoard(boardSlug string) *CompactBoard {
	boardURL := fmt.Sprintf(bp.GetLink("board"), boardSlug)
	boardStruct := &CompactBoard{}
	if err := bp.client.GetJSON(boardURL, &boardStruct); err != nil {
		log.Println("Board parse error", err)
	}
	return boardStruct
}

// ParseThread - Парсит тред и возвращает структуру
func (bp *BoardParser) ParseThread(boardSlug string, threadNum string) []*ThreadWithPosts {
	threadURL := fmt.Sprintf(bp.GetLink("thread"), boardSlug, threadNum, threadNum)
	threadStruct := []*ThreadWithPosts{}
	if err := bp.client.GetJSON(threadURL, &threadStruct); err != nil {
		log.Println("Thread parse error", err)
	}
	return threadStruct
}

// GetLink - Создает ссылку на пост
func (bp *BoardParser) GetLink(types string) string {
	switch types {
	case "board":
		return "https://2ch.hk/%s/threads.json"

	case "thread":
		return "https://2ch.hk/makaba/mobile.fcgi?task=get_thread&board=%s&thread=%s&num=%s"

	case "file":
		return "https://2ch.hk%s"

	default:
		return types
	}
}
