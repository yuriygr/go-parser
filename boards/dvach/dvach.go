package dvach

import (
	"fmt"

	"github.com/yuriygr/go-parser/boards"
	"github.com/yuriygr/go-parser/container"
	"github.com/yuriygr/go-parser/services"
)

// BuildBoard - Создание доски
func BuildBoard(container *container.Container) boards.Board {
	adapter := &Adapter{
		Client: *container.Client,
	}

	board := boards.Board{
		Name:       "Два.ч",
		URLAddress: "https://2ch.hk/",
		Adapter:    adapter,
	}

	return board
}

// Adapter - Адаптер
type Adapter struct {
	// HTTP клиент
	Client services.Client
}

// GetBoard - Получение раздела
func (a *Adapter) GetBoard(request boards.BoardRequest) (boards.CompactBoard, error) {
	boardURL := fmt.Sprintf(a.GetLink("board"), request.Board)

	boardStruct := boards.CompactBoard{}
	if err := a.Client.GetJSON(boardURL, &boardStruct); err != nil {
		return boardStruct, fmt.Errorf("Board parse error (%v)", err)
	}

	return boardStruct, nil
}

// GetThread - Получение треда
func (a *Adapter) GetThread(request boards.ThreadRequest) ([]boards.ThreadWithPosts, error) {
	threadURL := fmt.Sprintf(a.GetLink("thread"), request.Board, request.ThreadNumber, request.ThreadNumber)

	threadStruct := []boards.ThreadWithPosts{}
	if err := a.Client.GetJSON(threadURL, &threadStruct); err != nil {
		return threadStruct, fmt.Errorf("Thread parse error (%v)", err)
	}

	return threadStruct, nil
}

// GetLink - Создает ссылку на пост
func (a *Adapter) GetLink(types string) string {
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
