package pkg

import (
	"github.com/jmoiron/sqlx"
)

type MysqlRepository struct {
	*sqlx.DB
}

//NewMysqlRepository create new repository
func NewMysqlRepository(db *sqlx.DB) *MysqlRepository {
	return &MysqlRepository{
		db,
	}
}

func (repo *MysqlRepository) CreateTodoItem(item *TodoItem) error {
	tx := repo.MustBegin()
	result, err := tx.NamedExec("INSERT INTO todo_item(description, due_date) VALUES (:description, :due_date)", item)
	if err != nil {
		return err
	}
	item.Id, err = result.LastInsertId()
	return tx.Commit()
}

func (repo *MysqlRepository) LastTodoItem() (*TodoItem, error) {
	tx := repo.MustBegin()
	item := new(TodoItem)
	err := tx.Get(item, "SELECT id, description, due_date FROM todo_item ORDER BY id DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return item, err
}
