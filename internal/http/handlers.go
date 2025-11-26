package http

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"links-cheker-2511/internal/storage"
	"links-cheker-2511/pkg/link"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Storage interface {
	SaveUrl(id int, url string, status string, checked_at time.Time) error
	GetUrls(id int) ([]storage.Links, error)
}

type HTTPHandlers struct {
	storage Storage
	log     *slog.Logger
}

type SaveUrlsRequest struct {
	Links []string `json:"links"`
}
type SaveUrlsResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}
type GetUrlsRequest struct {
	LinksList []int `json:"links_list"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusError     = "Error"
	DecodeJSONError = "failed to decode request"
	SaveUrlsError   = "failed to save urls"
	GetUrlsError    = "failed to get urls"
)

var (
	linksMutex sync.Mutex
	currentID  int
	wg         sync.WaitGroup
)

func NewHTTPHandlers(storage Storage, log *slog.Logger) *HTTPHandlers {
	return &HTTPHandlers{
		storage: storage,
		log:     log,
	}
}

func (h *HTTPHandlers) HandleSaveUrls(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(
		slog.String("method", "HandleSaveUrls"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req SaveUrlsRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error(DecodeJSONError, wrapError(err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, createErrorResponse(DecodeJSONError))
		return
	}
	if req.Links == nil || len(req.Links) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, createErrorResponse(SaveUrlsError))
		return
	}

	linksMutex.Lock()
	currentID++
	setID := currentID
	linksMutex.Unlock()

	res := make(map[string]string)
	for _, currentLink := range req.Links {
		wg.Add(1)
		go func(lnk string) {
			defer wg.Done()

			status := link.CheckLink(lnk)

			linksMutex.Lock()
			res[currentLink] = status
			linksMutex.Unlock()
			err := h.storage.SaveUrl(setID, lnk, status, time.Now())
			if err != nil {
				log.Error(SaveUrlsError, wrapError(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, createErrorResponse(SaveUrlsError))
				return
			}

		}(currentLink)

	}
	wg.Wait()
	response := SaveUrlsResponse{
		Links:    res,
		LinksNum: setID,
	}
	log.Info("urls saved with ", slog.Int("id", setID))

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response)

}
func (h *HTTPHandlers) HandleGetUrls(w http.ResponseWriter, r *http.Request) {
	log := h.log.With(
		slog.String("method", "HandleGetUrls"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req GetUrlsRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error(DecodeJSONError, wrapError(err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, createErrorResponse(DecodeJSONError))
		return
	}
	if req.LinksList == nil || len(req.LinksList) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, createErrorResponse(SaveUrlsError))
		return
	}
	var allLinks []storage.Links
	for _, currentId := range req.LinksList {
		l, err := h.storage.GetUrls(currentId)
		if err != nil {
			log.Error(GetUrlsError, wrapError(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, createErrorResponse(GetUrlsError))
			return
		}
		allLinks = append(allLinks, l...)
	}

	fmt.Println("allLinks", allLinks)
}

func wrapError(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
func createErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Error:  msg,
		Status: StatusError,
	}
}
