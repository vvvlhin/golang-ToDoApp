package statistics_postgres_repository

import core_postgres_pool "github.com/vvvlhin/golang-ToDoApp/internal/core/repository/postges/conn"

type StatRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatRepository(
	pool core_postgres_pool.Pool,
) *StatRepository {
	return &StatRepository{
		pool: pool,
	}
}
