package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/vvvlhin/golang-ToDoApp/internal/core/logger"
	core_http_response "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/response"
	core_http_utils "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/utils"
)

type DeleteTaskResponse TaskDTOResponse

// DeleteTask godoc
//
//	@Summary		Удалить задачу
//	@Description	Удалить задачу у конкретного пользователя
//	@Tags			tasks
//	@Param			id 		path 		int 						true	"ID удаляемой задачи"
//	@Success		204														"Успешно удаленный пользователь"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Bad request"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"Task not found"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal server error"
//	@Router			/tasks/{id} [delete]

func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task id")
		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete task")
		return
	}

	responseHandler.NoContentResponse()
}
