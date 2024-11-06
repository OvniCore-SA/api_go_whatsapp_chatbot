package entities

type Monto int64

func (m Monto) Float64() float64 {
	return float64(m) / 100
}

func (m Monto) Int64() int64 {
	return int64(m)
}
