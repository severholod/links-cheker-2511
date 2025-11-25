package http

import "net/http"

type Storage interface {
	SaveUrls(urls []string) (int64, error)
	GetUrls(id int64) ([]string, error)
}

type HTTPHandlers struct {
	storage Storage
}

func NewHTTPHandlers(storage Storage) *HTTPHandlers {
	return &HTTPHandlers{
		storage: storage,
	}
}

func (h *HTTPHandlers) HandleSaveUrls(w http.ResponseWriter, r *http.Request) {

}
func (h *HTTPHandlers) HandleGetUrls(w http.ResponseWriter, r *http.Request) {
	//h.storage.GetUrls()
}
