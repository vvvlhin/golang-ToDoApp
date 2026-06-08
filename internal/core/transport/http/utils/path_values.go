package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/vvvlhin/golang-ToDoApp/internal/core/errors"
)

func GetIntPathValue(r *http.Request, path string) (int, error) {
	pathValue := r.PathValue(path)
	if pathValue == "" {
		return 0, fmt.Errorf("no key='%s' in path values: %w", path, core_errors.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf("path value='%s' by key='%s' not a valid integer: %w", err)
	}

	return val, nil
}
