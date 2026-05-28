package link

import (
	"fmt"
	"golang/configs"
	"golang/packages/event"
	"golang/packages/middleware"
	"golang/packages/req"
	"golang/packages/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	Config         *configs.Config
	EventBus       *event.EventBus
}

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}
	router.HandleFunc("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", middleware.IsAuth(handler.Update(), deps.Config))
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.Handle("GET /get_links", middleware.IsAuth(handler.getLinks(), deps.Config))
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LinkCreateRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 502)
			return
		}
		link := NewLink(body.Url)
		for {
			if handler.LinkRepository.CheckDuplicate("hash", link.Hash) {
				break
			}
			link.GenerateHash(10)
		}
		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdLink, http.StatusCreated)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !handler.LinkRepository.CheckDuplicate("id", id) {
			err = handler.LinkRepository.Delete(int(id))
			if err != nil {
				res.Json(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			res.Json(w, "Have not this ID in DB.", http.StatusNotFound)
			return
		}
		res.Json(w, "Deleted.", 200)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(middleware.ContextEmail).(string)
		if ok {
			fmt.Println(email)
		}
		body, err := req.HandleBody[LinkUpdateRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 502)
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !handler.LinkRepository.CheckDuplicate("hash", body.Hash) {
			res.Json(w, "This hash already using.", http.StatusBadRequest)
			return
		}
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			URL:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, link, http.StatusOK)
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			res.Json(w, err.Error(), http.StatusNotFound)
			return
		}
		//	handler.StatRepository.AddClick(link.ID)
		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})
		http.Redirect(w, r, link.URL, http.StatusAccepted)
	}
}

func (handler *LinkHandler) getLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		links := handler.LinkRepository.GetLinks(limit, offset)
		res.Json(w, GetAllLinks{
			Links: links,
			Count: handler.LinkRepository.Count(),
		}, 200)
	}
}
