package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/vvvlhin/golang-ToDoApp/internal/core/domain"
	core_logger "github.com/vvvlhin/golang-ToDoApp/internal/core/logger"
	core_http_request "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/request"
	core_http_response "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/response"
	core_http_types "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/types"
	core_http_utils "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/utils"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Выгулять собаку"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"В 15:00"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean" example:"false"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be NULL")
		}
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

// PatchUser godoc
//
//	@Summary		Изменить задачу пользователя
//	@Description	Изменить задачу пользователя по указанному ID задачи
//	@Description	### Логика обновления полей (Three-state logic):
//	@Description	1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
//	@Description	2. **Явно передано значение**: `description: "В 15:00"` - устанавливает новое значение в БД
//	@Description	3. **Передан null**: `description: null` - очищает поле в БД (set to NULL)
//	@Description	Ограничения `title` и `completed` - не может быть null
//	@Tags			tasks
//  @Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"ID изменяемой задачи"
//	@Param			request	body		PatchTaskRequest					true	"PatchUser тело запроса"
//	@Success		200		{object}	PatchTaskResponse					"Успешное изменение задачи пользователя"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Bad request"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"Task not found"
//	@Failure		409		{object}	core_http_response.ErrorResponse	"Conflict"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal server error"
//	@Router			/tasks/{id} [patch]

func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")
		return
	}

	var request PatchTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskPatch := taskPatchTaskRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func taskPatchTaskRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
