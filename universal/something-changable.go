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

type Closer interface {
	Close(ctx context.Context) error
}

type Closable interface {
	Closer
	Closed(ctx context.Context) (bool, error)
}

type Canceler interface {
	Cancel(ctx context.Context) error
}

type Cancelable interface {
	Canceler
	Canceled(ctx context.Context) (bool, error)
}

type Publisher interface {
	Publish(ctx context.Context) error
}

type Publishable interface {
	Publisher
	IsPublic(ctx context.Context) (bool, error)
}

type Activator interface {
	Activate(ctx context.Context) error
}

type Deactivator interface {
	Deactivate(ctx context.Context) error
}

type Activable interface {
	Activator
	Deactivator
	IsActive(ctx context.Context) (bool, error)
}
