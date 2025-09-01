package offer

import "github.com/lukaszkaleta/saas-go/universal"

type GlobalOffers interface {
	NearBy(position *universal.RadarModel) ([]Offer, error)
	ById(id int64) (Offer, error)
}
