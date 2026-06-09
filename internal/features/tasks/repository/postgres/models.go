package tasks_postgres_repository

import (
	"time"

	"github.com/vvvlhin/golang-ToDoApp/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Vesion       int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func NewTaskModel(
	taskModel TaskModel,
) *TaskModel {
	return &TaskModel{
		ID:           taskModel.ID,
		Vesion:       taskModel.Vesion,
		Title:        taskModel.Title,
		Description:  taskModel.Description,
		Completed:    taskModel.Completed,
		CreatedAt:    taskModel.CreatedAt,
		CompletedAt:  taskModel.CompletedAt,
		AuthorUserID: taskModel.AuthorUserID,
	}
}

func taskDomainFromModel(task TaskModel) domain.Task {
	return domain.NewTask(
		task.ID,
		task.Vesion,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)

}

func taskDomainsFromModels(tasks []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(tasks))

	for i, model := range tasks {
		domains[i] = taskDomainFromModel(model)
	}

	return domains
}
