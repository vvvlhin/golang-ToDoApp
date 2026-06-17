package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/vvvlhin/golang-ToDoApp/internal/core/logger"
	core_http_response "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/response"
	core_http_utils "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/utils"
)

type GetTasksResponse []TaskDTOResponse

var (
	userIDParam = "user_id"
	limitParam  = "limit"
	offsetParam = "offset"
)

// GetTasks godoc

// @Summary Получить все задачи пользователя
// @Description Получить список всех задач пользователя по ID
// @Tags tasks
// @Produce json
// @Param 	id 		query int 	false 								"ID автора задач"
// @Param 	limit 	query int 	false 								"Размер страницы с задачами"
// @Param 	offset 	query int 	false 								"Смещение относительно 1-ой задачи страницы с задачами"
// @Success	200 	{object} 	GetTasksResponse 					"Список задач"
// @Failure	400		{object}	core_http_response.ErrorResponse	"BadRequest"
// @Failure	500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
// @Router	/tasks	[get]

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID / limit / offset query params")
		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasksDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	userID, err := core_http_utils.GetQueryParamInt(r, userIDParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `user_id` query param: %w", err)
	}

	limit, err := core_http_utils.GetQueryParamInt(r, limitParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetQueryParamInt(r, offsetParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
