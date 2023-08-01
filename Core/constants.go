package core

type FilterType string

const (
	Equal        FilterType = "equal"
	LessThan     FilterType = "lt"
	LessEqual    FilterType = "lte"
	GreaterThan  FilterType = "gt"
	GreaterEqual FilterType = "gte"
	Contains     FilterType = "contains"
	IContains    FilterType = "icontains"
)

type ContextKey struct {
	Name string
}
