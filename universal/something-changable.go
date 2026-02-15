package universal

import "context"

type Acceptor interface {
	Accept(ctx context.Context) error
}

type Acceptable interface {
	Acceptor
	Accepted() (bool, error)
}

type Rejecter interface {
	Reject(ctx context.Context) error
}

type Rejectable interface {
	Rejecter
	Rejected() (bool, error)
}

type Releaser interface {
	Release(ctx context.Context) error
}

type Releasable interface {
	Releaser
	Released() (bool, error)
}
