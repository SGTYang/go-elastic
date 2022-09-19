package post

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SGTYang/gorest/tree/dev/gorest/domain"
	"github.com/SGTYang/gorest/tree/dev/gorest/elastic"
)

type Handler struct {
	service service
}

func NewHandler(storage elastic.PostStorer) Handler {
	return Handler{
		service: service{storage: storage},
	}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	res, err := h.service.create(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrConflict:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.Header().Set("content-Type", "Application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	bdy, _ := json.Marshal(res)
	_, _ = w.Write(bdy)
}
