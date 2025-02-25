package domain

type Publisher interface {
	Publish(event Event) error
}
