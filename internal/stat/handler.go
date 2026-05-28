package stat

import (
	"fmt"
	"golang/configs"
	"golang/packages/res"
	"net/http"
	"time"
)

const (
	FilterByDay   = "day"
	FilterByMonth = "month"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

type StatHandler struct {
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, deps *StatHandlerDeps) {
	stat := &StatHandler{
		StatRepository: deps.StatRepository,
	}
	router.HandleFunc("GET /stat", stat.GetStat())
}

func (stat *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const layout = "2006-01-02"

		dateFrom := r.URL.Query().Get("from")
		fmt.Println(dateFrom)
		t, err := time.Parse(layout, dateFrom)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dateTo := r.URL.Query().Get("to")
		t1, err := time.Parse(layout, dateTo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dateBy := r.URL.Query().Get("by")
		if dateBy != FilterByDay && dateBy != FilterByMonth {
			http.Error(w, "Invalid Params.", http.StatusBadRequest)
			return
		}

		stats := stat.StatRepository.GetStats(dateBy, t, t1)
		res.Json(w, stats, 200)
	}
}
