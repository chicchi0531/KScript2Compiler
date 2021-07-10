package vm

type StringTable struct{
	Values []string
}

func MakeStringTable() *StringTable{
	t := new(StringTable)
	t.Values = make([]string, 0)
	return t
}

func (s *StringTable) Add(str string) int {
	s.Values = append(s.Values, str)
	return len(s.Values)-1
}