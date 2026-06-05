package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/vvvlhin/golang-ToDoApp/internal/core/errors"
)

func GetQueryParamInt(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param=%s by key=%s not a valid integer: %v : %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}

// func GetQueryParamTime(r *http.Request, key string) time.Time {
// 	param := r.URL.Query().Get(key)
// }
