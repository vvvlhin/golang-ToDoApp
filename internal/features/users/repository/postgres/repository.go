package users_postgres_repository

import core_postgres_pool "github.com/vvvlhin/golang-ToDoApp/internal/core/repository/postges/conn"

type UsersRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(
	pool core_postgres_pool.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
