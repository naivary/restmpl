package dependency

type Pinger interface {
	Ping() error
}

type dep[T any] struct {
	v    T
	ping func(T) error
}

func (d dep[T]) Ping() error {
	if d.ping == nil {
		return nil
	}
	return d.ping(d.v)
}

func New[T any](ping func(T) error, v T) dep[T] {
	return dep[T]{
		v:    v,
		ping: ping,
	}
}
