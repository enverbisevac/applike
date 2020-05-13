package pkg

type Repository interface {
	CreateTodoItem(item *TodoItem) error
	LastTodoItem() (*TodoItem, error)
}
