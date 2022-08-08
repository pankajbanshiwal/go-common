package encoding

type Encoder interface {
	Encode(v interface{}) error
}
