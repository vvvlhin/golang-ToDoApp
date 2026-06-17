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
	TasksCreated           int      `json:"tasks_created" example:"2"`
	TasksCompleted         int      `json:"tasks_completed" example:"1"`
	TasksCompletedRate     *float64 `json:"tasks_completed_rate" example:"50"`
	TasksAvgCompletionTime *string  `json:"tasks_avg_completion_time" example:"46.0000ms"`
}

var (
	userIDParam = "user_id"
	fromParam   = "from"
	toParam     = "to"
)

// GetStatistics godoc

// @Summary Получить статистику
// @Description Получить статистику по ID пользователя с опциональной фильтрацией и/или временному промежутку
// @Tags statistics
// @Produce json
// @Param 	user_id 		query int 			false 										"ID автора задач"
// @Param 	from 	query string 		false 										"Начало промежутка времени задач"
// @Param 	to 		query string 		false 										"Конец промежутка времени задач"
// @Success	200 	{object} 			GetStatResponse 							"Успешное получение статистики"
// @Failure	400		{object}			core_http_response.ErrorResponse			"BadRequest"
// @Failure	500		{object}			core_http_response.ErrorResponse			"Internal Server Error"
// @Router	/statistics	[get]

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
