package users_postgres_repository

import "github.com/vvvlhin/golang-ToDoApp/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainFromModel(users UserModel) domain.User {
	return domain.NewUser(
		users.ID,
		users.Version,
		users.FullName,
		users.PhoneNumber,
	)
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		userDomains[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber,
		)
	}

	return userDomains
}
