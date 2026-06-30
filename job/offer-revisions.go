package job

import (
	"context"

	"github.com/lukaszkaleta/saas-go/universal"
)

type OfferRevisions interface {
	universal.Creator[OfferRevisionModel, OfferRevision]
	universal.Lister[OfferRevision]

	ById(ctx context.Context, id int64) (OfferRevision, error)
	FromUser(ctx context.Context, id int64) (OfferRevision, error)
	NewestFromWorker(ctx context.Context) (OfferRevision, error)
	NewestFromOwner(ctx context.Context) (OfferRevision, error)
	Accepted(ctx context.Context) (OfferRevision, error)
	Newest(ctx context.Context) (OfferRevision, error)
}
