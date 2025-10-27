package universal

type ModelAware interface {
	Model() HasId
}

type HasId interface {
	GetId() int64
}

func ModelsToMap[T ModelAware](array []T) map[int64]T {
	idMap := make(map[int64]T)
	for _, modelAware := range array {
		idMap[modelAware.Model().GetId()] = modelAware
	}
	return idMap
}

func HasIdToMap[T HasId](array []T) map[int64]T {
	idMap := make(map[int64]T)
	for _, modelAware := range array {
		idMap[modelAware.GetId()] = modelAware
	}
	return idMap
}
