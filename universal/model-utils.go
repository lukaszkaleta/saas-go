package universal

type ModelAware interface {
	Model() HasId
}

type HasId interface {
	GetId() int64
}

func ModelsToMap(array []ModelAware) map[int64]ModelAware {
	idMap := make(map[int64]ModelAware)
	for _, modelAware := range array {
		idMap[modelAware.Model().GetId()] = modelAware
	}
	return idMap
}

func HasIdToMap(array []HasId) map[int64]HasId {
	idMap := make(map[int64]HasId)
	for _, modelAware := range array {
		idMap[modelAware.GetId()] = modelAware
	}
	return idMap
}
