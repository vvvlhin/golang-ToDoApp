package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/vvvlhin/golang-ToDoApp/internal/core/domain"
	core_errors "github.com/vvvlhin/golang-ToDoApp/internal/core/errors"
)

func (s *StatService) GetStats(
	ctx context.Context,
	userID *int,
	fromParam *time.Time,
	toParam *time.Time,
) (domain.Statistics, error) {
	if fromParam != nil && toParam != nil {
		if toParam.Before(*fromParam) || toParam.Equal(*fromParam) {
			return domain.Statistics{}, fmt.Errorf("`to` must be afted `from`: %w", core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.statRepository.GetStats(ctx, userID, fromParam, toParam)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks: %w", err)
	}

	stat := calcStatistics(tasks)

	return stat, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{
			TasksCreated:           0,
			TasksCompleted:         0,
			TasksCompletedRate:     nil,
			TasksAvgCompletionTime: nil,
		}
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0

	var totalCompletedDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletedDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAvgCompletionTime *time.Duration

	if tasksCompleted > 0 && totalCompletedDuration != 0 {
		avg := totalCompletedDuration / time.Duration(tasksCompleted)

		tasksAvgCompletionTime = &avg
	}

	return domain.Statistics{
		TasksCreated:           tasksCreated,
		TasksCompleted:         tasksCompleted,
		TasksCompletedRate:     &tasksCompletedRate,
		TasksAvgCompletionTime: tasksAvgCompletionTime,
	}
}
