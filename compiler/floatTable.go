package compiler

type FloatTable struct{
	values []float32
}

func (t *FloatTable) Add(value float32) int {
	t.values = append(t.values, value)
	return len(t.values)-1
}