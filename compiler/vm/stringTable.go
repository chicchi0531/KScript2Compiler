package vm

type StringTable struct{
	Values []string
}

func (s *StringTable) Add(str string) int {
	s.Values = append(s.Values, str)
	return len(s.Values)-1
}