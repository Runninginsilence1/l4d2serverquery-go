package queryitem

import "strconv"

// CurrentPlayerCountQueryItem 查询现在的玩家
type CurrentPlayerCountQueryItem struct {
	value       int
	compareType CompareType
}

func (c *CurrentPlayerCountQueryItem) PropertyName() string {
	return "Players"
}

func (c *CurrentPlayerCountQueryItem) Value() string {
	value := strconv.Itoa(c.value)
	return value
}

func (c *CurrentPlayerCountQueryItem) CompareType() CompareType {
	return c.compareType
}

func NewCurrentPlayerCountQueryItem(value int, compareType NumberCompareType) QueryItem {
	//if compareType != GreaterThanOrEqual {
	//	panic("对于数值运算, 无效的类型")
	//}
	return &CurrentPlayerCountQueryItem{value: value, compareType: CompareType(compareType)}
}

// MaximumPlayerCountQueryItem 查询服务器容纳玩家数量
type MaximumPlayerCountQueryItem struct {
	value       int
	compareType CompareType
}

func (m *MaximumPlayerCountQueryItem) PropertyName() string {
	return "Slots"
}

func (m *MaximumPlayerCountQueryItem) Value() string {
	value := strconv.Itoa(m.value)
	return value
}

func (m *MaximumPlayerCountQueryItem) CompareType() CompareType {
	return m.compareType
}
