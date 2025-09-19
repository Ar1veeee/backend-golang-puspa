package child

import "backend-golang/internal/domain/repositories"

type Dependencies struct {
	ChildRepo repositories.ChildRepository
	//Validator       Validator
	Mapper Mapper
}

func NewDependencies(
	observationRepo repositories.ChildRepository,
) *Dependencies {
	return &Dependencies{
		ChildRepo: observationRepo,
		Mapper:    NewChildMapper(),
	}
}
