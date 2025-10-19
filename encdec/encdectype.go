package encdec

//go:generate stringer -type=EncDecType
type EncDecType int

const (
	EncDecTypeHraban EncDecType = iota
	EncDecTypeJJ11hh EncDecType = iota
)
