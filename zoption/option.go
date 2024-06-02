package zoption

type Option[T any] interface {
	Apply(T) error
}

func Build[T any](param T, opts ...Option[T]) error {

	for _, opt := range opts {
		if err := opt.Apply(param); err != nil {
			return err
		}
	}
	return nil
}

type FuncOption[T any] func(T) error

func (f FuncOption[T]) Apply(t T) error {
	return f(t)
}
