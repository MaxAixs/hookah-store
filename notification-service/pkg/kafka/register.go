package kafka

type Register interface {
	RegisterHandler(topic string, handler EventHandler)
}
