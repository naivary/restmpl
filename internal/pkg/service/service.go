package service

type Service[T any] interface {
	// Unique identifier of the service
	UUID() string

	// Name of the service which may not be unique
	Name() string

	// Detailed description of the service
	Description() string

	// Register is passing the root router
	// so the Service can register itself.
	Register(T)
}
