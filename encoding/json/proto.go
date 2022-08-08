package json

import (
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/okcredit/go-common/protox"
	"io"
	"reflect"
)

var (
	DefaultProtoMarshaler = jsonpb.Marshaler{
		EnumsAsInts:  true,
		EmitDefaults: true,
		OrigName:     true,
	}

	DefaultProtoUnmarshaler = jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
)

// ProtoArray provides custom json marshalling for array of proto messages 
type protoArray []proto.Message

var _ json.Marshaler = protoArray(nil)

func (pbArr protoArray) MarshalJSON() ([]byte, error) {
	data := make([]byte, 0)
	data = append(data, []byte("[")...)

	lastIndex := len(pbArr) - 1
	for index, pb := range pbArr {
		pbData, err := Marshal(pb)
		if err != nil {
			return nil, err
		}
		data = append(data, pbData...)

		if index != lastIndex {
			data = append(data, []byte(",")...)
		}
	}

	data = append(data, []byte("]")...)
	return data, nil
}

// encoder
func encodeProto(out io.Writer, v interface{}) (err error, isProto bool) {
	// check if v is a proto message
	if pb, ok := protox.Parse(v); ok {
		return DefaultProtoMarshaler.Marshal(out, pb), true
	}

	// check if v is an array of proto message
	if pbArr, ok := protox.ParseArray(v); ok {
		return json.NewEncoder(out).Encode(protoArray(pbArr)), true
	}

	// not proto related
	return nil, false
}

// decoder
func decodeProto(in io.Reader, v interface{}) (err error, isProto bool) {
	// check if v is a proto message
	if pb, ok := protox.Parse(v); ok {
		return DefaultProtoUnmarshaler.Unmarshal(in, pb), true
	}

	// check if v is a pointer to array of proto messages
	v_type := reflect.TypeOf(v)
	if v_type.Kind() == reflect.Ptr && (v_type.Elem().Kind() == reflect.Slice || v_type.Elem().Kind() == reflect.Array) {
		arr_type := v_type.Elem()
		elem_type := arr_type.Elem()

		arr := reflect.MakeSlice(reflect.SliceOf(elem_type), 0, 0)

		jsonDecoder := json.NewDecoder(in)

		// ignore opening of array, [
		_, err := jsonDecoder.Token()
		if err != nil {
			return err, true
		}

		for jsonDecoder.More() {

			elem := reflect.Zero(elem_type)

			if elem_type.Kind() == reflect.Struct {
				// pointer to elem is a proto message

				pb, _ := protox.Parse(reflect.Zero(elem_type).Interface())
				err := jsonpb.UnmarshalNext(jsonDecoder, pb)
				if err != nil {
					return err, true
				}
				elem = reflect.ValueOf(pb).Elem()

			} else {
				// elem is a proto message

				pb, _ := protox.Parse(reflect.Zero(elem_type.Elem()).Interface())
				err := jsonpb.UnmarshalNext(jsonDecoder, pb)
				if err != nil {
					return err, true
				}
				elem = reflect.ValueOf(pb)
			}

			arr = reflect.Append(arr, elem)
		}

		reflect.ValueOf(v).Elem().Set(arr)

		return nil, true
	}

	// not proto related
	return nil, false
}
