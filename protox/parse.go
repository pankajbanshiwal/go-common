package protox

import (
	"github.com/golang/protobuf/proto"
	"reflect"
)

// Parse parses v into a proto message if:
// 		v implements proto.Message, or
// 		pointer to v implements proto.Message
func Parse(v interface{}) (proto.Message, bool) {
	// check if v itself is a proto message
	if _, ok := v.(proto.Message); ok {
		return v.(proto.Message), true
	}

	// check if v is a struct and pointer to v is a proto message
	v_type := reflect.TypeOf(v)
	if v_type.Kind() == reflect.Struct {
		// v is a struct

		ptr_type := reflect.PtrTo(v_type)
		pb_type := reflect.TypeOf((*proto.Message)(nil)).Elem()

		if ptr_type.Implements(pb_type) {
			// pointer to v is a proto message

			pb_ptr := reflect.New(v_type)         // pointer to given struct
			pb_ptr.Elem().Set(reflect.ValueOf(v)) // set underlying value of pointer to thr given struct's value
			pb, ok := (pb_ptr.Interface()).(proto.Message)
			return pb, ok
		}
	}

	return nil, false
}

// ParseArray parses v into a slice of proto messages if
// 		v is an array/slice of T, where T implements proto.Message, or
//		v is an array/slick of T, where pointer to T implements proto.Message, or
//		v is an array/slick of interface{}, where each element is T or pointer to T, such that pointer to T implements proto.Message
func ParseArray(v interface{}) ([]proto.Message, bool) {
	v_type := reflect.TypeOf(v)
	if v_type.Kind() == reflect.Ptr {
		// dereference pointer
		v_type = v_type.Elem()
	}

	if v_type.Kind() == reflect.Slice || v_type.Kind() == reflect.Array {
		// v (or its underlying value) is an array

		v_val := reflect.ValueOf(v)
		pbArr := make([]proto.Message, 0)
		for index := 0; index < v_val.Len(); index = index + 1 {
			pb, ok := Parse(v_val.Index(index).Interface())
			if !ok {
				return nil, false
			}
			pbArr = append(pbArr, pb)
		}

		return pbArr, true
	}

	return nil, false
}
