package pg

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/finance"
	"github.com/lukaszkaleta/saas-go/job"
	"github.com/lukaszkaleta/saas-go/payment"
)

type LedgerPaymentsCreator struct {
	creator payment.PaymentsCreator
	ledger  finance.FinancialLedger
	job     job.Job
}

func NewLedgerPaymentsCreator(creator payment.PaymentsCreator, ledger finance.FinancialLedger, job job.Job) payment.PaymentsCreator {
	return &LedgerPaymentsCreator{
		creator: creator,
		ledger:  ledger,
		job:     job,
	}
}

func (l *LedgerPaymentsCreator) Create(ctx context.Context, offer any) (payment.Intent, error) {
	intent, err := l.creator.Create(ctx, offer)
	if err != nil {
		return nil, err
	}

	intentModel, err := intent.Model(ctx)
	if err != nil {
		return nil, err
	}

	jobID := intentModel.JobId
	buyerId := intentModel.PayerId
	sellerId := intentModel.PayeeId

	entry := finance.LedgerEntry{
		JobID:      &jobID,
		SellerID:   sellerId,
		BuyerID:    &buyerId,
		Amount:     intentModel.Amount,
		Currency:   intentModel.Currency,
		OccurredAt: time.Now(),
	}

	_, err = l.ledger.Events().PayoutRelease(ctx, entry)
	if err != nil {
		return nil, err
	}

	return intent, nil
}
