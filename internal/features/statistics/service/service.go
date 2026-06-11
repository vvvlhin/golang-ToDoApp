package statistics_service

import (
	"context"
	"time"

	"github.com/vvvlhin/golang-ToDoApp/internal/core/domain"
)

type StatService struct {
	statRepository StatRepository
}

type StatRepository interface {
	GetStats(
		ctx context.Context,
		userID *int,
		fromParam *time.Time,
		toParam *time.Time,
	) ([]domain.Task, error)
}

func NewStatService(
	statRepository StatRepository,
) StatService {
	return StatService{
		statRepository: statRepository,
	}
}
