package offer

import (
	"fmt"
	"time"

	"github.com/lukaszkaleta/saas-go/universal"
)

type GlobalOffersLogger struct {
	next GlobalOffers
}

func (g *GlobalOffersLogger) ById(id int64) (Offer, error) {
	//TODO implement me
	panic("implement me")
}

func NewGlobalOffersLogger(next GlobalOffers) GlobalOffers {
	return &GlobalOffersLogger{
		next: next,
	}
}

func (g *GlobalOffersLogger) NearBy(position *universal.RadarModel) ([]Offer, error) {
	defer func(start time.Time) {
		fmt.Printf("Searching for offers took %v\n", time.Since(start))
	}(time.Now())
	return g.next.NearBy(position)
}
