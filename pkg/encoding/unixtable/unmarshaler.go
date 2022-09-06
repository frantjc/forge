package unixtable

type Unmarshaler interface {
	UnmarshalUnixTable([]byte) error
}
