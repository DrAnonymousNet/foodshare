package core

import (
	"reflect"

	"github.com/google/uuid"
)

func IsZero(value interface{}) bool {
	if value == nil {
		return true
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return val.Len() == 0
	case reflect.Bool:
		return !val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return val.IsNil()
	case reflect.Struct:
		if val.Type() == reflect.TypeOf(uuid.UUID{}) {
			// For uuid.UUID, check if all bytes are zero
			zeroUUID := uuid.UUID{}
			return value.(uuid.UUID) == zeroUUID
		}
	}

	return false
}
