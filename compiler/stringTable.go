package compiler

type StringTable struct{
	values []string
}

func (s *StringTable) Add(str string) int {
	s.values = append(s.values, str)
	return len(s.values)-1
}