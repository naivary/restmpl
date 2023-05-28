package service

type Service[T any] interface {
	// Unique identifier of the service. Recommended is
	// to assign an UUID which you generated beforehand
	// to be able to reference the service in other services
	// or external systems.
	ID() string

	// Name of the service which may not be unique
	Name() string

	// Detailed description of the service
	Description() string

	// Register registers the service to
	// root router of type `T` of the env.
	Register(T)
}
