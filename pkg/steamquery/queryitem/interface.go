package queryitem

type QueryCondition struct {
	PropertyName string      `json:"propertyName"`
	Value        string      `json:"value"`
	CompareType  CompareType `json:"comparison"`
}

type CompareType int

type (
	NumberCompareType CompareType
	StringCompareType CompareType
)

const (
	Equal NumberCompareType = iota
	NotEqual
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
)

const (
	Contains StringCompareType = 6
)

type QueryItem interface {
	PropertyName() string
	Value() string
	CompareType() CompareType
}
