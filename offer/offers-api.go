package offer

import (
	"strconv"

	"github.com/lukaszkaleta/saas-go/universal"
)

type Offers interface {
	AddFromPosition(model *universal.PositionModel) (Offer, error)
}

func OfferModels(offers []Offer) []*OfferModel {
	var models []*OfferModel
	for _, modelAware := range offers {
		models = append(models, modelAware.Model()) // note the = instead of :=
	}
	return models
}

func OfferHints(offers []Offer) []*OfferHint {
	var hints []*OfferHint
	for _, o := range offers {
		if o != nil {
			hints = append(hints, o.Model().Hint()) // note the = instead of :=
		}
	}
	return hints
}

func GeoOffers(offers []Offer) universal.GeoFeatureCollection[OfferHint] {
	features := make([]universal.GeoFeature[OfferHint], 0, len(offers))
	for i := range offers {
		m := offers[i]
		pt := universal.NewGeoPoint(m.Model().Position.LonF(), m.Model().Position.LatF())
		features = append(features, universal.NewGeoFeature[OfferHint](strconv.FormatInt(m.Model().Id, 10), pt, *m.Model().Hint()))
	}
	return universal.NewGeoFeatureCollection(features)
}
