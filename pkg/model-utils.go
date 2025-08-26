package universal

type ModelAware[I any] interface {
	Model() I
}
