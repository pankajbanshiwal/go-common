package json

import (
	"bytes"
	"encoding/json"
	"github.com/okcredit/go-common/encoding"
	"io"
	
)

func NewEncoder(w io.Writer) encoding.Encoder {
	return &encoder{out: w}
}

type encoder struct {
	out io.Writer
}

func (e *encoder) Encode(v interface{}) error {
	if err, isProto := encodeProto(e.out, v); isProto {
		return err
	}

	// std lib encoder
	return json.NewEncoder(e.out).Encode(v)
}

func Marshal(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
