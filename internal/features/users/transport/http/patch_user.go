package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vvvlhin/golang-ToDoApp/internal/core/domain"
	core_logger "github.com/vvvlhin/golang-ToDoApp/internal/core/logger"
	core_http_request "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/request"
	core_http_response "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/response"
	core_http_types "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/types"
	core_http_utils "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Максим Максимович"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79998887766"`
}

type PatchUserResponse UserDTOResponse

// PatchUser godoc
//
//	@Summary		Изменить данные пользователя
//	@Description	Изменить данные пользователя по указанному ID
//	@Description	### Логика обновления полей (Three-state logic):
//	@Description	1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
//	@Description	2. **Явно передано значение**: `phone_number: "+79998887766"` - устанавливает новое значение в БД
//	@Description	3. **Передан null**: `phone_number: null` - очищает поле в БД (set to NULL)
//	@Description	Ограничения `full_name` - не может быть null
//	@Tags			users
//  @Accept			json
//	@Produce		json
//	@Param			id	path			int									true	"ID изменяемого пользователя"
//	@Param			request	body		PatchUserRequest					true	"PatchUser тело запроса"
//	@Success		200		{object}	PatchUserResponse					"Успешное изменение данных пользователя"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Bad request"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"User not found"
//	@Failure		409		{object}	core_http_response.ErrorResponse	"Conflict"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal server error"
//	@Router			/users/{id} [patch]

func (p *PatchUserRequest) Validate() error {
	if p.FullName.Set {
		if p.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*p.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100")
		}

		if p.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*p.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15")
			}

			if !strings.HasPrefix(*p.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must starts with '+'")
			}
		}
	}

	return nil
}

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	var request PatchUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.userService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	log.Debug(
		fmt.Sprintf(
			"Fields:\nFull_name: %v\nPhone_Number: %v\n",
			request.FullName,
			request.PhoneNumber,
		),
	)

	response := PatchUserResponse(userDtoFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
