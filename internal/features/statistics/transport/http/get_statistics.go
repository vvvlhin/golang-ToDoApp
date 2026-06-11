package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vvvlhin/golang-ToDoApp/internal/core/domain"
	core_logger "github.com/vvvlhin/golang-ToDoApp/internal/core/logger"
	core_http_response "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/response"
	core_http_utils "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/utils"
)

type GetStatResponse struct {
	TasksCreated           int      `json:"tasks_created"`
	TasksCompleted         int      `json:"tasks_completed"`
	TasksCompletedRate     *float64 `json:"tasks_completed_rate"`
	TasksAvgCompletionTime *string  `json:"tasks_avg_completion_time"`
}

var (
	userIDParam = "user_id"
	fromParam   = "from"
	toParam     = "to"
)

func (h *StatHTTPHandler) GetStats(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get `userID` | `from` | `to` params from path value")
		return
	}

	stat, err := h.statService.GetStats(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	response := toDTOFromDomain(stat)

	responseHandler.JSONResponse(response, http.StatusOK)

}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	userID, err := core_http_utils.GetQueryParamInt(r, userIDParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `userID` param from path value: %w", err)
	}

	from, err := core_http_utils.GetDateQueryParam(r, fromParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `from` param from path value: %w", err)
	}

	to, err := core_http_utils.GetDateQueryParam(r, toParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `to` param from path value: %w", err)
	}

	return userID, from, to, nil
}

func toDTOFromDomain(stat domain.Statistics) GetStatResponse {
	var avgTime *string
	if stat.TasksAvgCompletionTime != nil {
		duration := stat.TasksAvgCompletionTime.String()
		avgTime = &duration
	}
	return GetStatResponse{
		TasksCreated:           stat.TasksCreated,
		TasksCompleted:         stat.TasksCompleted,
		TasksCompletedRate:     stat.TasksCompletedRate,
		TasksAvgCompletionTime: avgTime,
	}
}
