package json

import (
	"bytes"
	"encoding/json"
	"github.com/okcredit/go-common/encoding"
	"io"
)

func NewDecoder(r io.Reader) encoding.Decoder {
	return &decoder{in: r}
}

type decoder struct {
	in io.Reader
}

func (d *decoder) Decode(v interface{}) error {
	if err, isProto := decodeProto(d.in, v); isProto {
		return err
	}

	// std lib decoder
	return json.NewDecoder(d.in).Decode(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(data)).Decode(v)
}
