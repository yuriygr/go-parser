package boards

// Board - Структура чана
type Board struct {
	Name       string
	URLAddress string
	Adapter    AdapterInterface
}

// GetBoard - Получение раздела
func (b *Board) GetBoard(request BoardRequest) (CompactBoard, error) {
	return b.Adapter.GetBoard(request)
}

// GetThread - Получение треда
func (b *Board) GetThread(request ThreadRequest) ([]ThreadWithPosts, error) {
	return b.Adapter.GetThread(request)
}

// GetLink - Получение ссылки на что-то
func (b *Board) GetLink(request string) string {
	return b.Adapter.GetLink(request)
}

// AdapterInterface - Интерфейс портала
// Каждый адаптер может добывать нужные данные как хочет,
// главное чтобы в нужном формате
type AdapterInterface interface {
	GetBoard(BoardRequest) (CompactBoard, error)
	GetThread(ThreadRequest) ([]ThreadWithPosts, error)
	GetLink(string) string
}

// FileInsert - Структура файла для вставки в БД
type FileInsert struct {
	ID        int    `db:"v.id"`
	URL       string `db:"v.url"`
	Type      string `db:"v.type"`
	Filesize  int    `db:"v.filesize"`
	Width     int    `db:"v.width"`
	Height    int    `db:"v.height"`
	Md5       string `db:"v.md5"`
	Name      string `db:"v.name"`
	Board     string `db:"v.board"`
	Thread    string `db:"v.thread"`
	CreatedAt int64  `db:"v.created_at"`
}

// BoardParserParams - параметры для парсинга
type BoardParserParams struct {
	BoardSlug string
	FileType  string
}

// ThreadRequest -
type ThreadRequest struct {
	Board        string
	ThreadNumber string
}

// BoardRequest -
type BoardRequest struct {
	Board string
}

/**
 * sadas
 */

type CompactBoard struct {
	Board   string `json:"board"`
	Threads []struct {
		Comment    string  `json:"comment"`
		Lasthit    int     `json:"lasthit"`
		Num        string  `json:"num"`
		PostsCount int     `json:"posts_count"`
		Score      float64 `json:"score"`
		Subject    string  `json:"subject"`
		Timestamp  int     `json:"timestamp"`
		Views      int     `json:"views"`
	} `json:"threads"`
}

type ThreadWithPosts struct {
	Banned        int64   `json:"banned"`
	Closed        int64   `json:"closed"`
	Comment       string  `json:"comment"`
	Date          string  `json:"date"`
	Email         string  `json:"email"`
	Endless       int64   `json:"endless"`
	Files         []File  `json:"files"`
	Lasthit       int64   `json:"lasthit"`
	Name          string  `json:"name"`
	Num           string  `json:"num"`
	Op            int8    `json:"op"`
	Parent        string  `json:"parent"`
	Sticky        int8    `json:"sticky"`
	Subject       string  `json:"subject"`
	Tags          *string `json:"tags,omitempty"`
	Timestamp     int64   `json:"timestamp"`
	Trip          string  `json:"trip"`
	TripType      string  `json:"trip_type"`
	UniquePosters *string `json:"unique_posters,omitempty"`
}

type File struct {
	Displayname string `json:"displayname"`
	Fullname    string `json:"fullname"`
	Height      int    `json:"height"`
	Md5         string `json:"md5"`
	Name        string `json:"name"`
	Nsfw        int8   `json:"nsfw"`
	Path        string `json:"path"`
	Size        int    `json:"size"`
	Thumbnail   string `json:"thumbnail"`
	TnHeight    int    `json:"tn_height"`
	TnWidth     int    `json:"tn_width"`
	Type        int    `json:"type"`
	Width       int    `json:"width"`
}

func (f *File) GetType() string {
	switch f.Type {
	case 1, 2:
		return "Image"
	case 6, 10:
		return "Video"
	case 100:
		return "Sticker"
	default:
		return "Image"
	}
}

func (f *File) GetMimeType() string {
	switch f.Type {
	case 6:
		return "video/webm"
	case 10:
		return "video/mp4"
	default:
		return ""
	}
}

func (f *File) RealName() string {
	if f.Fullname != "" {
		return f.Fullname
	}
	if f.Name != "" {
		return f.Name
	}
	return ""
}
