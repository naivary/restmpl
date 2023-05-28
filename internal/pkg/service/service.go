package service

type Service[T any] interface {
	// Unique identifier of the service.
	ID() string

	// Human friendly name of the service.
	Name() string

	// Detailed description of the service
	Description() string

	// Register registers the service
	// to the public router of type T
	Register(T)

	// Health reutrns the health status
	// of the service. If the error is
	// non nil the service is considered unhealthy.
	// Health(http.Responsewriter, *http.Request) error
}
