package service

type Service[T any] interface {
	UUID() string
	Name() string

	// Detailed description of the service
	Description() string

	// Register is passing the root router
	// so the Service can register itself.
	Register(T)
}
