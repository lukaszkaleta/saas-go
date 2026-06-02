package finance

import (
	"context"
)

type EventType string

const (
	EventEscrowHold    EventType = "escrow_hold"
	EventPayoutRelease EventType = "payout_release"
	EventPlatformFee   EventType = "platform_fee"
	EventPayout        EventType = "payout"
	EventRefund        EventType = "refund"
	EventChargeback    EventType = "chargeback"
	EventAdjustment    EventType = "adjustment"
)

type FinancialLedger interface {
	Record(ctx context.Context, entry LedgerEntry) (int64, error)

	// Event helper methods
	EscrowHold(ctx context.Context, entry LedgerEntry) (int64, error)
	PayoutRelease(ctx context.Context, entry LedgerEntry) (int64, error)
	PlatformFee(ctx context.Context, entry LedgerEntry) (int64, error)
	Payout(ctx context.Context, entry LedgerEntry) (int64, error)
	Refund(ctx context.Context, entry LedgerEntry) (int64, error)
	Chargeback(ctx context.Context, entry LedgerEntry) (int64, error)
	Adjustment(ctx context.Context, entry LedgerEntry) (int64, error)
}
