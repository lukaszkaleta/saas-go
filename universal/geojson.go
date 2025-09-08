package universal

// Generic GeoJSON structures with typed Properties.
// This is intentionally minimal and independent from external libraries.
// It supports Feature and FeatureCollection with any geometry and typed properties.

// Point geometry in integer coordinates to match universal.PositionModel (Lat,Lon int)
type GeoPoint struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

// Feature is a generic GeoJSON Feature with typed properties.
// T is the type of the Properties payload.
type GeoFeature[T any] struct {
	Type       string `json:"type"`
	Geometry   any    `json:"geometry"`
	Properties T      `json:"properties"`
}

// FeatureCollection is a generic GeoJSON FeatureCollection with typed properties on each feature.
// T is the type used for Properties across all Features.
type GeoFeatureCollection[T any] struct {
	Type     string          `json:"type"`
	Features []GeoFeature[T] `json:"features"`
}

// NewGeoPoint creates a Point geometry from lon/lat (GeoJSON uses [lon, lat]).
func NewGeoPoint(lon, lat float64) GeoPoint {
	return GeoPoint{Type: "Point", Coordinates: [2]float64{lon, lat}}
}

// NewFeature creates a Feature with provided geometry and properties.
func NewGeoFeature[T any](geometry any, properties T) GeoFeature[T] {
	return GeoFeature[T]{Type: "Feature", Geometry: geometry, Properties: properties}
}

// NewFeatureCollection constructs a FeatureCollection from features.
func NewGeoFeatureCollection[T any](features []GeoFeature[T]) GeoFeatureCollection[T] {
	return GeoFeatureCollection[T]{Type: "FeatureCollection", Features: features}
}
