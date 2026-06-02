package finance

import (
	"context"
	"encoding/xml"
	"sort"
	"time"
)

type DAC7Message struct {
	XMLName  xml.Name    `xml:"DAC7_Message"`
	Header   HeaderXML   `xml:"Header"`
	Platform PlatformXML `xml:"PlatformOperator"`
	Reports  []ReportXML `xml:"Report"`
}

type ReportXML struct {
	Seller SellerXML `xml:"Seller"`
}

type PlatformXML struct {
	PlatformName string `xml:"PlatformName"`
	TIN          string `xml:"TIN"`
	CountryCode  string `xml:"CountryCode"`
}

type HeaderXML struct {
	MessageRefID string    `xml:"MessageRefId"`
	CreationDate time.Time `xml:"CreationDate"`
	CountryCode  string    `xml:"CountryCode"` // NO

	MessageType string `xml:"MessageType"` // DAC7
}

type SellerXML struct {
	SellerID int64 `xml:"SellerId"`

	Identity IdentityXML `xml:"Identity"`
	Address  AddressXML  `xml:"Address"`
	Tax      TaxXML      `xml:"TaxIdentification"`

	Activity ActivityXML `xml:"Activity"`
}

type IdentityXML struct {
	FirstName string `xml:"FirstName"`
	LastName  string `xml:"LastName"`
	Email     string `xml:"Email"`
	Phone     string `xml:"Phone"`
}

type AddressXML struct {
	Line1    string `xml:"Street"`
	Line2    string `xml:"Street2,omitempty"`
	City     string `xml:"City"`
	PostCode string `xml:"PostCode"`
	Region   string `xml:"Region"`
	Country  string `xml:"CountryCode"`
}

type TaxXML struct {
	TIN            string `xml:"TIN"`
	IssuingCountry string `xml:"TIN_Issuer"`

	DateOfBirth *string `xml:"DateOfBirth,omitempty"`
}

type ActivityXML struct {
	Currency string `xml:"Currency"`

	NumberOfActivities int64 `xml:"NumberOfActivities"`

	Consideration ConsiderationXML `xml:"Consideration"`

	Fees FeesXML `xml:"Fees"`
}

type FeesXML struct {
	PlatformFee int64 `xml:"PlatformFee"`
}

type ConsiderationXML struct {
	Gross int64 `xml:"GrossAmount"`
	Net   int64 `xml:"NetAmount"`
}

type EarningsXML struct {
	Currency string `xml:"Currency"`

	NumberOfActivities int64 `xml:"NumberOfActivities"`

	Gross int64 `xml:"GrossEarnings"`
	Fees  int64 `xml:"Fees"`
	Net   int64 `xml:"NetEarnings"`

	Q1 int64 `xml:"Q1"`
	Q2 int64 `xml:"Q2"`
	Q3 int64 `xml:"Q3"`
	Q4 int64 `xml:"Q4"`
}

type Dac7Xml struct {
	repo DacReporting
}

func NewDac7Xml(repo DacReporting) *Dac7Xml {
	return &Dac7Xml{repo: repo}
}

func (s *Dac7Xml) Generate(ctx context.Context, start, end time.Time) ([]byte, error) {

	rows, err := s.repo.Dac7Report(ctx, start, end)
	if err != nil {
		return nil, err
	}

	fees, err := s.repo.PlatformFees(ctx, start, end)
	if err != nil {
		return nil, err
	}

	feeMap := make(map[int64]int64)
	for _, f := range fees {
		feeMap[f.SellerID] = f.TotalPlatformFeesMinor
	}

	reports := make([]ReportXML, 0, len(rows))

	for _, r := range rows {

		report := ReportXML{
			Seller: SellerXML{
				SellerID: r.SellerID,

				Identity: IdentityXML{
					FirstName: r.PersonFirstName,
					LastName:  r.PersonLastName,
					Email:     r.PersonEmail,
					Phone:     r.PersonPhone,
				},

				Address: AddressXML{
					Line1:    r.AddressLine1,
					Line2:    r.AddressLine2,
					City:     r.AddressCity,
					PostCode: r.AddressPostalCode,
					Region:   r.AddressDistrict,
					Country:  r.CountryCode,
				},

				Tax: TaxXML{
					TIN:            r.TaxIdentificationNumber,
					IssuingCountry: r.CountryCode,
					DateOfBirth:    formatDOB(r.DateOfBirth),
				},

				Activity: ActivityXML{
					Currency:           r.Currency,
					NumberOfActivities: r.NumberOfActivities,

					Consideration: ConsiderationXML{
						Gross: r.TotalGrossEarningsMinor,
						Net:   r.TotalNetEarningsMinor,
					},

					Fees: FeesXML{
						PlatformFee: feeMap[r.SellerID],
					},
				},
			},
		}

		reports = append(reports, report)
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Seller.SellerID < reports[j].Seller.SellerID
	})

	doc := DAC7Message{
		Header: HeaderXML{
			MessageRefID: "NABORLY-" + time.Now().Format("20060102-150405"),
			CreationDate: time.Now().UTC(),
			CountryCode:  "NO",
			MessageType:  "DAC7",
		},

		Platform: PlatformXML{
			PlatformName: "Naborly",
			TIN:          "YOUR_PLATFORM_TIN",
			CountryCode:  "NO",
		},

		Reports: reports,
	}

	return xml.MarshalIndent(doc, "", "  ")
}

func formatDOB(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format("2006-01-02")
	return &s
}
