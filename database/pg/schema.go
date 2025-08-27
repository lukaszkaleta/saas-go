package pg

type Schema interface {
	Create() error
	Drop() error
}
