package finance

import (
	"context"
)

type EventType int

const (
	EventEscrowHold    EventType = 1
	EventPayoutRelease EventType = 2
	EventPlatformFee   EventType = 3
	EventPayout        EventType = 4
	EventRefund        EventType = 5
	EventChargeback    EventType = 6
	EventAdjustment    EventType = 7
)

type FinancialLedger interface {
	Record(ctx context.Context, entry LedgerEntry) (int64, error)
	JobView(ctx context.Context, jobID int64) ([]LedgerEntry, error)
	SellerReporting(sellerID int64) SellerReporting
	DacReporting() DacReporting

	Events() FinancialEvents
}
