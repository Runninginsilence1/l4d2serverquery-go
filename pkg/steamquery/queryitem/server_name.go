package queryitem

// ServerNameQueryItem 查询服务器名称
type ServerNameQueryItem struct {
	value string
}

func (s *ServerNameQueryItem) PropertyName() string {
	return "Name"
}

func (s *ServerNameQueryItem) Value() string {
	return s.value
}

func (s *ServerNameQueryItem) CompareType() CompareType {
	return CompareType(Contains)
}

func NewServerNameQueryItem(value string) QueryItem {
	return &ServerNameQueryItem{value: value}
}
