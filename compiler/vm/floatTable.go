package vm

type FloatTable struct{
	Values []float32
}

func (t *FloatTable) Add(value float32) int {
	t.Values = append(t.Values, value)
	return len(t.Values)-1
}