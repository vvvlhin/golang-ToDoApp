package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/vvvlhin/golang-ToDoApp/internal/core/domain"
	core_http_server "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/server"
)

type StatHTTPHandler struct {
	statService StatService
}

type StatService interface {
	GetStats(
		ctx context.Context,
		userID *int,
		fromParam *time.Time,
		toParam *time.Time,
	) (domain.Statistics, error)
}

func NewStatHTTPHandler(statService StatService) *StatHTTPHandler {
	return &StatHTTPHandler{
		statService: statService,
	}
}

func (h *StatHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStats,
		},
	}
}
