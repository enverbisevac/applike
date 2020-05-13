package pkg

import "gopkg.in/go-playground/validator.v9"

type Service struct {
	repo Repository
}

//NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

//CreateWorkspace based on name
func (service *Service) CreateTodoItem(item *TodoItem) error {
	validate := validator.New() // this can be handled much better
	err := validate.Struct(item)
	if err != nil {
		return err
	}
	return service.repo.CreateTodoItem(item)
}

func (service *Service) LastTodoItem() (*TodoItem, error) {
	return service.repo.LastTodoItem()
}
