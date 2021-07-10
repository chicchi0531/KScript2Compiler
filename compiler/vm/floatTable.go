package vm

type FloatTable struct{
	Values []float32
}

func MakeFloatTable() *FloatTable{
	t := new(FloatTable)
	t.Values = make([]float32, 0)
	return t
}

func (t *FloatTable) Add(value float32) int {
	t.Values = append(t.Values, value)
	return len(t.Values)-1
}