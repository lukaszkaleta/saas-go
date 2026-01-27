package universal

type ModelAware interface {
	Model() Idable
}

func ModelsToMap[T ModelAware](array []T) map[int64]T {
	idMap := make(map[int64]T)
	for _, modelAware := range array {
		idMap[modelAware.Model().ID()] = modelAware
	}
	return idMap
}

func IdableToMap[T Idable](array []T) map[int64]T {
	idMap := make(map[int64]T)
	for _, modelAware := range array {
		idMap[modelAware.ID()] = modelAware
	}
	return idMap
}
