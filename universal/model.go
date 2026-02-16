package universal

import "context"

type ModelAware[M Idable] interface {
	Idable
	Model(ctx context.Context) (*M, error)
}

func ModelsArray[T ModelAware[M], M Idable](ctx context.Context, array []T) ([]*M, error) {
	models := make([]*M, len(array))
	for i, modelAware := range array {
		model, err := modelAware.Model(ctx)
		if err != nil {
			return nil, err
		}
		models[i] = model
	}
	return models, nil
}

func ModelsToMap[T ModelAware[T]](ctx context.Context, array []T) (map[int64]T, error) {
	idMap := make(map[int64]T)
	for _, modelAware := range array {
		model, err := modelAware.Model(ctx)
		if err != nil {
			return nil, err
		}
		x := *model
		idMap[x.ID()] = modelAware
	}
	return idMap, nil
}

func IdableToMap[T Idable](array []T) map[int64]T {
	idMap := make(map[int64]T)
	for _, idable := range array {
		idMap[idable.ID()] = idable
	}
	return idMap
}
