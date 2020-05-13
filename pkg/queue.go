package pkg

type Queue interface {
	PublishMessage(item TodoItem) error
	ReceiveMessage() (*TodoItem, error)
}
