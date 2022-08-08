package encoding

type Decoder interface {
	Decode(v interface{}) error
}
