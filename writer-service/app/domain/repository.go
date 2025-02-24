package domain

type Repository interface {
	Write(event Event) error
}
