package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/vvvlhin/golang-ToDoApp/internal/core/logger"
	core_http_response "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/response"
	core_http_utils "github.com/vvvlhin/golang-ToDoApp/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

// GetUsers godoc
//
//	@Summary		Получить список пользователей
//	@Description	Получить список всех пользователей с опциональной пагинацией
//	@Tags			users
//	@Produce		json
//	@Param			limit	query		int									false	"Размер страницы с пользователем"
//	@Param			offset	query		int									false	"Смещение относительно 1-ого пользователя"
//	@Success		200		{object}	GetUsersResponse					"Успешное получение списка пользователей"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Bad request"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"User not found"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal server error"
//	@Router			/users [get]

func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit'/'offset param'",
		)
		return
	}

	userDomains, err := h.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(userDtoFromDomains(userDomains))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetQueryParamInt(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetQueryParamInt(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
