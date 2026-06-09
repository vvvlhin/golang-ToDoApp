package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/vvvlhin/golang-ToDoApp/internal/core/logger"
	core_http_response "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/response"
	core_http_utils "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/utils"
)

type GetTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get `ID` from request")
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(err, "falied to get task")
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}
