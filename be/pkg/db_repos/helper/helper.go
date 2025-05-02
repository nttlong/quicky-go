package helper

type IHelper interface {
	Connect() error

	Close() error
	GetConnectionString() string
}
