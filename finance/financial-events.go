package finance

import (
	"context"
)

type FinancialEvents interface {
	EscrowHold(ctx context.Context, entry LedgerEntry) (int64, error)
	PayoutRelease(ctx context.Context, entry LedgerEntry) (int64, error)
	PlatformFee(ctx context.Context, entry LedgerEntry) (int64, error)
	Payout(ctx context.Context, entry LedgerEntry) (int64, error)
	Refund(ctx context.Context, entry LedgerEntry) (int64, error)
	Chargeback(ctx context.Context, entry LedgerEntry) (int64, error)
	Adjustment(ctx context.Context, entry LedgerEntry) (int64, error)
}

type financialEvents struct {
	recorder interface {
		Record(ctx context.Context, entry LedgerEntry) (int64, error)
	}
}

func NewFinancialEvents(recorder interface {
	Record(ctx context.Context, entry LedgerEntry) (int64, error)
}) FinancialEvents {
	return &financialEvents{recorder: recorder}
}

func (e *financialEvents) EscrowHold(ctx context.Context, entry LedgerEntry) (int64, error) {
	entry.Type = EventEscrowHold
	return e.recorder.Record(ctx, entry)
}

func (e *financialEvents) PayoutRelease(ctx context.Context, entry LedgerEntry) (int64, error) {
	entry.Type = EventPayoutRelease
	return e.recorder.Record(ctx, entry)
}

func (e *financialEvents) PlatformFee(ctx context.Context, entry LedgerEntry) (int64, error) {
	entry.Type = EventPlatformFee
	return e.recorder.Record(ctx, entry)
}

func (e *financialEvents) Payout(ctx context.Context, entry LedgerEntry) (int64, error) {
	entry.Type = EventPayout
	return e.recorder.Record(ctx, entry)
}

func (e *financialEvents) Refund(ctx context.Context, entry LedgerEntry) (int64, error) {
	entry.Type = EventRefund
	return e.recorder.Record(ctx, entry)
}

func (e *financialEvents) Chargeback(ctx context.Context, entry LedgerEntry) (int64, error) {
	entry.Type = EventChargeback
	return e.recorder.Record(ctx, entry)
}

func (e *financialEvents) Adjustment(ctx context.Context, entry LedgerEntry) (int64, error) {
	entry.Type = EventAdjustment
	return e.recorder.Record(ctx, entry)
}
