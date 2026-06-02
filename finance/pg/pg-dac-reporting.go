package pg

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/finance"
	"github.com/lukaszkaleta/saas-go/universal"
)

type PgDacReporting struct {
	db *pg.PgDb
}

func (r *PgDacReporting) SellerEarnings(ctx context.Context, sellerID int64, year int) (finance.SellerEarnings, error) {
	const query = `
		SELECT
			seller_id,
			SUM(CASE WHEN EXTRACT(QUARTER FROM occurred_at)=1 THEN amount ELSE 0 END) q1,
			SUM(CASE WHEN EXTRACT(QUARTER FROM occurred_at)=2 THEN amount ELSE 0 END) q2,
			SUM(CASE WHEN EXTRACT(QUARTER FROM occurred_at)=3 THEN amount ELSE 0 END) q3,
			SUM(CASE WHEN EXTRACT(QUARTER FROM occurred_at)=4 THEN amount ELSE 0 END) q4,
			SUM(amount) total
		FROM financial_ledger
		WHERE seller_id = $1
		  AND type = $2
		  AND occurred_at >= $3
		  AND occurred_at < $4
		GROUP BY seller_id;
	`

	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	var res finance.SellerEarnings
	err := r.db.Pool.QueryRow(ctx, query, sellerID, finance.EventPayout, start, end).Scan(
		&res.SellerID,
		&res.Q1,
		&res.Q2,
		&res.Q3,
		&res.Q4,
		&res.Total,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return finance.SellerEarnings{SellerID: sellerID}, nil
		}
		return finance.SellerEarnings{}, err
	}

	return res, nil
}

func (r *PgDacReporting) Dac7Report(ctx context.Context, dateRange universal.DateRange) ([]finance.Dac7ReportRow, error) {
	const query = `
		SELECT
			u.id AS seller_id,

			u.person_first_name,
			u.person_last_name,
			u.person_email,
			u.person_phone,

			u.address_line_1,
			u.address_line_2,
			u.address_city,
			u.address_postal_code,
			u.address_district,

			u.tax_identification_number,
			u.date_of_birth,
			u.country_code,

			fl.currency,

			COUNT(DISTINCT fl.job_id) AS number_of_activities,

			SUM(fl.amount) AS total_gross_earnings_minor,

			SUM(COALESCE(fl.fee_amount, 0)) AS total_fees_minor,

			SUM(COALESCE(fl.net_amount, fl.amount - COALESCE(fl.fee_amount, 0)))
				AS total_net_earnings_minor

		FROM financial_ledger fl
		JOIN users u
			ON u.id = fl.seller_id

		WHERE
			fl.type = 2 -- payout_release (DAC7 earnings event)
			AND fl.occurred_at >= $1
			AND fl.occurred_at < $2

		GROUP BY
			u.id,
			u.person_first_name,
			u.person_last_name,
			u.person_email,
			u.person_phone,
			u.address_line_1,
			u.address_line_2,
			u.address_city,
			u.address_postal_code,
			u.address_district,

			u.tax_identification_number,
			u.date_of_birth,
			u.country_code,

			fl.currency

		ORDER BY
			u.id;
	`

	rows, err := r.db.Pool.Query(ctx, query, dateRange.From, dateRange.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report []finance.Dac7ReportRow
	for rows.Next() {
		var row finance.Dac7ReportRow
		err := rows.Scan(
			&row.SellerID,
			&row.PersonFirstName,
			&row.PersonLastName,
			&row.PersonEmail,
			&row.PersonPhone,
			&row.AddressLine1,
			&row.AddressLine2,
			&row.AddressCity,
			&row.AddressPostalCode,
			&row.AddressDistrict,
			&row.TaxIdentificationNumber,
			&row.DateOfBirth,
			&row.CountryCode,
			&row.Currency,
			&row.NumberOfActivities,
			&row.TotalGrossEarningsMinor,
			&row.TotalFeesMinor,
			&row.TotalNetEarningsMinor,
		)
		if err != nil {
			return nil, err
		}
		report = append(report, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}

func (r *PgDacReporting) PlatformFees(ctx context.Context, dateRange universal.DateRange) ([]finance.PlatformFees, error) {
	const query = `
		SELECT
			seller_id,
			currency,
			SUM(fee_amount) AS total_platform_fees_minor

		FROM financial_ledger
		WHERE
			type = $1
			AND occurred_at >= $2
			AND occurred_at < $3

		GROUP BY seller_id, currency;
	`

	rows, err := r.db.Pool.Query(ctx, query, finance.EventPayoutRelease, dateRange.From, dateRange.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report []finance.PlatformFees
	for rows.Next() {
		var row finance.PlatformFees
		err := rows.Scan(
			&row.SellerID,
			&row.Currency,
			&row.TotalPlatformFeesMinor,
		)
		if err != nil {
			return nil, err
		}
		report = append(report, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}
