package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/vvvlhin/golang-ToDoApp/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
		FROM todoapp.tasks
		%s
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2
	`

	args := []any{limit, offset}

	if userID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id=$3")
		args = append(args, userID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	row, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return []domain.Task{}, fmt.Errorf("select tasks: %w", err)
	}
	defer row.Close()

	var taskModels []TaskModel

	for row.Next() {
		var taskModel TaskModel

		err = row.Scan(
			&taskModel.ID,
			&taskModel.Vesion,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan err: %w", err)
		}

		taskModels = append(taskModels, taskModel)

	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	tasksDTO := taskDomainsFromModels(taskModels)

	return tasksDTO, nil

}
