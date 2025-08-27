package pg

type Schema interface {
	Create() error
}
