package service

type publisher interface {
	Publish(subject string, message []byte)
}
