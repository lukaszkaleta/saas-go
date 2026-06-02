package pg_test

import (
	"testing"
	"time"

	"github.com/lukaszkaleta/saas-go/database/pg"
	pgfilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/finance"
	financepg "github.com/lukaszkaleta/saas-go/finance/pg"
	pgjob "github.com/lukaszkaleta/saas-go/job/pg"
	"github.com/lukaszkaleta/saas-go/universal"
	pguser "github.com/lukaszkaleta/saas-go/user/pg"
)

func financeSetupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("postgres", "postgres")

	// Drop and recreate table to ensure clean state
	financepg.NewFinanceSchema(db).Drop()
	pgjob.NewJobSchema(db).Drop()
	pguser.NewUserSchema(db).Drop()
	pgfilestore.NewFilestoreSchema(db).Drop()

	// Create required schemas in order
	pgfilestore.NewFilestoreSchema(db).Create()
	pguser.NewUserSchema(db).Create()
	pgjob.NewJobSchema(db).Create()
	financepg.NewFinanceSchema(db).Create()

	return func(tb testing.TB) {
		financepg.NewFinanceSchema(db).Drop()
		pgjob.NewJobSchema(db).Drop()
		pguser.NewUserSchema(db).Drop()
		pgfilestore.NewFilestoreSchema(db).Drop()
	}, db
}

func TestDacReporting_SellerEarnings(t *testing.T) {
	teardown, db := financeSetupTest(t)
	defer teardown(t)

	ctx := t.Context()

	// 1. Create a seller
	var sellerID int64
	err := db.Pool.QueryRow(ctx, "INSERT INTO users (person_first_name) VALUES ('Seller') RETURNING id").Scan(&sellerID)
	if err != nil {
		t.Fatal(err)
	}

	ledger := financepg.NewPgFinancialLedger(db)
	dac := ledger.DacReporting()

	year := 2026
	jan1 := time.Date(year, 1, 1, 10, 0, 0, 0, time.UTC)
	apr1 := time.Date(year, 4, 1, 10, 0, 0, 0, time.UTC)
	jul1 := time.Date(year, 7, 1, 10, 0, 0, 0, time.UTC)
	oct1 := time.Date(year, 10, 1, 10, 0, 0, 0, time.UTC)

	// 2. Insert some events
	// Event 4 = payout -> SHOULD BE COUNTED
	events := []struct {
		amount int64
		at     time.Time
	}{
		{1000, jan1}, // Q1
		{2000, apr1}, // Q2
		{3000, jul1}, // Q3
		{4000, oct1}, // Q4
		{500, jan1},  // Q1 again
	}

	for _, e := range events {
		_, err = db.Pool.Exec(ctx, `
			INSERT INTO financial_ledger (seller_id, type, amount, currency, occurred_at)
			VALUES ($1, $2, $3, $4, $5)`,
			sellerID, finance.EventPayout, e.amount, "NOK", e.at)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Event 2 = payout_release -> SHOULD BE IGNORED
	_, err = db.Pool.Exec(ctx, `
		INSERT INTO financial_ledger (seller_id, type, amount, currency, occurred_at)
		VALUES ($1, $2, $3, $4, $5)`,
		sellerID, finance.EventPayoutRelease, 1000, "NOK", jan1)
	if err != nil {
		t.Fatal(err)
	}

	// Outside year -> SHOULD BE IGNORED
	_, err = db.Pool.Exec(ctx, `
		INSERT INTO financial_ledger (seller_id, type, amount, currency, occurred_at)
		VALUES ($1, $2, $3, $4, $5)`,
		sellerID, finance.EventPayout, 5000, "NOK", time.Date(year-1, 12, 31, 23, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	// 3. Verify
	earnings, err := dac.SellerEarnings(ctx, sellerID, year)
	if err != nil {
		t.Fatal(err)
	}

	if earnings.SellerID != sellerID {
		t.Errorf("expected sellerID %d, got %d", sellerID, earnings.SellerID)
	}
	if earnings.Q1 != 1500 {
		t.Errorf("expected Q1 1500, got %d", earnings.Q1)
	}
	if earnings.Q2 != 2000 {
		t.Errorf("expected Q2 2000, got %d", earnings.Q2)
	}
	if earnings.Q3 != 3000 {
		t.Errorf("expected Q3 3000, got %d", earnings.Q3)
	}
	if earnings.Q4 != 4000 {
		t.Errorf("expected Q4 4000, got %d", earnings.Q4)
	}
	if earnings.Total != 10500 {
		t.Errorf("expected Total 10500, got %d", earnings.Total)
	}
}

func TestSellerReporting_SumInPeriod(t *testing.T) {
	teardown, db := financeSetupTest(t)
	defer teardown(t)

	ctx := t.Context()

	var sellerID int64
	err := db.Pool.QueryRow(ctx, "INSERT INTO users (person_first_name) VALUES ('Seller') RETURNING id").Scan(&sellerID)
	if err != nil {
		t.Fatal(err)
	}

	ledger := financepg.NewPgFinancialLedger(db)
	reporting := ledger.SellerReporting(sellerID)

	now := time.Now().Truncate(time.Second)
	period := universal.DateRange{
		From: now.Add(-time.Hour),
		To:   now.Add(time.Hour),
	}

	// Insert events
	_, err = db.Pool.Exec(ctx, `
		INSERT INTO financial_ledger (seller_id, type, amount, currency, occurred_at)
		VALUES ($1, $2, $3, $4, $5)`,
		sellerID, finance.EventPayoutRelease, 2500, "NOK", now)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Pool.Exec(ctx, `
		INSERT INTO financial_ledger (seller_id, type, amount, currency, occurred_at)
		VALUES ($1, $2, $3, $4, $5)`,
		sellerID, finance.EventPayout, 1000, "NOK", now)
	if err != nil {
		t.Fatal(err)
	}

	sum, err := reporting.SumInPeriod(ctx, period)
	if err != nil {
		t.Fatal(err)
	}

	if sum != 2500 {
		t.Errorf("expected 2500, got %d", sum)
	}
}
