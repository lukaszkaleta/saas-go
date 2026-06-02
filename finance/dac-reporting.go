package finance

import (
	"context"
	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

type SellerEarnings struct {
	SellerID int64
	Q1       int64
	Q2       int64
	Q3       int64
	Q4       int64
	Total    int64
}

type Dac7ReportRow struct {
	SellerID                int64
	PersonFirstName         string
	PersonLastName          string
	PersonEmail             string
	PersonPhone             string
	AddressLine1            string
	AddressLine2            string
	AddressCity             string
	AddressPostalCode       string
	AddressDistrict         string
	TaxIdentificationNumber string
	DateOfBirth             *time.Time
	CountryCode             string
	Currency                string
	NumberOfActivities      int64
	TotalGrossEarningsMinor int64
	TotalFeesMinor          int64
	TotalNetEarningsMinor   int64
}

type PlatformFees struct {
	SellerID               int64
	Currency               string
	TotalPlatformFeesMinor int64
}

type DacReporting interface {
	SellerEarnings(ctx context.Context, sellerID int64, year int) (SellerEarnings, error)
	Dac7Report(ctx context.Context, dateRange universal.DateRange) ([]Dac7ReportRow, error)
	PlatformFees(ctx context.Context, dateRange universal.DateRange) ([]PlatformFees, error)
}
