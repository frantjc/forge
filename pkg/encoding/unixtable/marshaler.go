package unixtable

type Marshaler interface {
	MarshalUnixTable() ([]byte, error)
}
