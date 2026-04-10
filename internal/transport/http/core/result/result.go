package result

type Paged[T any] struct {
	Items    []T
	HasMore  bool
	LastItem *T
}
