package binding

import (
	"context"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// Converter is used to convert []string to specific type.
// The data **MUST** have at least one element.
type Converter func(ctx context.Context, data []string) (interface{}, error)

var converters = map[reflect.Type]Converter{
	reflect.TypeOf(bool(false)):  ConvertToBool,
	reflect.TypeOf(int(0)):       ConvertToInt,
	reflect.TypeOf(int8(0)):      ConvertToInt8,
	reflect.TypeOf(int16(0)):     ConvertToInt16,
	reflect.TypeOf(int32(0)):     ConvertToInt32,
	reflect.TypeOf(int64(0)):     ConvertToInt64,
	reflect.TypeOf(uint(0)):      ConvertToUint,
	reflect.TypeOf(uint8(0)):     ConvertToUint8,
	reflect.TypeOf(uint16(0)):    ConvertToUint16,
	reflect.TypeOf(uint32(0)):    ConvertToUint32,
	reflect.TypeOf(uint64(0)):    ConvertToUint64,
	reflect.TypeOf(float32(0)):   ConvertToFloat32,
	reflect.TypeOf(float64(0)):   ConvertToFloat64,
	reflect.TypeOf(string("")):   ConvertToString,
	reflect.TypeOf(new(bool)):    ConvertToBoolP,
	reflect.TypeOf(new(int)):     ConvertToIntP,
	reflect.TypeOf(new(int8)):    ConvertToInt8P,
	reflect.TypeOf(new(int16)):   ConvertToInt16P,
	reflect.TypeOf(new(int32)):   ConvertToInt32P,
	reflect.TypeOf(new(int64)):   ConvertToInt64P,
	reflect.TypeOf(new(uint)):    ConvertToUintP,
	reflect.TypeOf(new(uint8)):   ConvertToUint8P,
	reflect.TypeOf(new(uint16)):  ConvertToUint16P,
	reflect.TypeOf(new(uint32)):  ConvertToUint32P,
	reflect.TypeOf(new(uint64)):  ConvertToUint64P,
	reflect.TypeOf(new(float32)): ConvertToFloat32P,
	reflect.TypeOf(new(float64)): ConvertToFloat64P,
	reflect.TypeOf(new(string)):  ConvertToStringP,
	reflect.TypeOf([]bool{}):     ConvertToBoolSlice,
	reflect.TypeOf([]int{}):      ConvertToIntSlice,
	reflect.TypeOf([]float64{}):  ConvertToFloat64Slice,
	reflect.TypeOf([]string{}):   ConvertToStringSlice,
}

// ConverterFor gets converter for specified type.
func ConverterFor(typ reflect.Type) Converter {
	return converters[typ]
}

// RegisterConverter registers a converter for specified type. New converter
// overrides old one.
func RegisterConverter(typ reflect.Type, converter Converter) {
	converters[typ] = converter
}

// ConvertToBool converts []string to bool. Only the first data is used.
func ConvertToBool(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseBool(origin)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type bool", origin)
	}
	return target, nil
}

// ConvertToBoolP converts []string to *bool. Only the first data is used.
func ConvertToBoolP(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToBool(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(bool)
	return &value, nil
}

// ConvertToInt converts []string to int. Only the first data is used.
func ConvertToInt(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type int", origin)
	}
	return int(target), nil
}

// ConvertToIntP converts []string to *int. Only the first data is used.
func ConvertToIntP(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToInt(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(int)
	return &value, nil
}

// ConvertToInt8 converts []string to int8. Only the first data is used.
func ConvertToInt8(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 8)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type int8", origin)
	}
	return int8(target), nil
}

// ConvertToInt8P converts []string to *int8. Only the first data is used.
func ConvertToInt8P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToInt8(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(int8)
	return &value, nil
}

// ConvertToInt16 converts []string to int16. Only the first data is used.
func ConvertToInt16(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 16)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type int16", origin)
	}
	return int16(target), nil
}

// ConvertToInt16P converts []string to *int16. Only the first data is used.
func ConvertToInt16P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToInt16(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(int16)
	return &value, nil
}

// ConvertToInt32 converts []string to int32. Only the first data is used.
func ConvertToInt32(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 32)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type int32", origin)
	}
	return int32(target), nil
}

// ConvertToInt32P converts []string to *int32. Only the first data is used.
func ConvertToInt32P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToInt32(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(int32)
	return &value, nil
}

// ConvertToInt64 converts []string to int64. Only the first data is used.
func ConvertToInt64(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseInt(origin, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type int64", origin)
	}
	return target, nil
}

// ConvertToInt64P converts []string to *int64. Only the first data is used.
func ConvertToInt64P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToInt64(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(int64)
	return &value, nil
}

// ConvertToUint converts []string to uint. Only the first data is used.
func ConvertToUint(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type uint", origin)
	}
	return uint(target), nil
}

// ConvertToUintP converts []string to *uint. Only the first data is used.
func ConvertToUintP(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToUint(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(uint)
	return &value, nil
}

// ConvertToUint8 converts []string to uint8. Only the first data is used.
func ConvertToUint8(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 8)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type uint8", origin)
	}
	return uint8(target), nil
}

// ConvertToUint8P converts []string to *uint8. Only the first data is used.
func ConvertToUint8P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToUint8(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(uint8)
	return &value, nil
}

// ConvertToUint16 converts []string to uint16. Only the first data is used.
func ConvertToUint16(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 16)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type uint16", origin)
	}
	return uint16(target), nil
}

// ConvertToUint16P converts []string to *uint16. Only the first data is used.
func ConvertToUint16P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToUint16(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(uint16)
	return &value, nil
}

// ConvertToUint32 converts []string to uint32. Only the first data is used.
func ConvertToUint32(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 32)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type uint32", origin)
	}
	return uint32(target), nil
}

// ConvertToUint32P converts []string to *uint32. Only the first data is used.
func ConvertToUint32P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToUint32(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(uint32)
	return &value, nil
}

// ConvertToUint64 converts []string to uint64. Only the first data is used.
func ConvertToUint64(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseUint(origin, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type uint64", origin)
	}
	return target, nil
}

// ConvertToUint64P converts []string to *uint64. Only the first data is used.
func ConvertToUint64P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToUint64(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(uint64)
	return &value, nil
}

// ConvertToFloat32 converts []string to float32. Only the first data is used.
func ConvertToFloat32(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseFloat(origin, 32)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type float32", origin)
	}
	return float32(target), nil
}

// ConvertToFloat32P converts []string to *float32. Only the first data is used.
func ConvertToFloat32P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToFloat32(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(float32)
	return &value, nil
}

// ConvertToFloat64 converts []string to float64. Only the first data is used.
func ConvertToFloat64(ctx context.Context, data []string) (interface{}, error) {
	origin := data[0]
	target, err := strconv.ParseFloat(origin, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "Cann't convert %s to type float64", origin)
	}
	return target, nil
}

// ConvertToFloat64P converts []string to *float64. Only the first data is used.
func ConvertToFloat64P(ctx context.Context, data []string) (interface{}, error) {
	ret, err := ConvertToFloat64(ctx, data)
	if err != nil {
		return nil, err
	}
	value := ret.(float64)
	return &value, nil
}

// ConvertToString return the first element in []string.
func ConvertToString(ctx context.Context, data []string) (interface{}, error) {
	return data[0], nil
}

// ConvertToStringP return the first element's pointer in []string.
func ConvertToStringP(ctx context.Context, data []string) (interface{}, error) {
	return &data[0], nil
}

// ConvertToBoolSlice converts all elements in data to bool, and return []bool
func ConvertToBoolSlice(ctx context.Context, data []string) (interface{}, error) {
	ret := make([]bool, len(data))
	for i := range data {
		r, err := ConvertToBool(ctx, data[i:i+1])
		if err != nil {
			return nil, err
		}
		ret[i] = r.(bool)
	}
	return ret, nil
}

// ConvertToIntSlice converts all elements in data to int, and return []int
func ConvertToIntSlice(ctx context.Context, data []string) (interface{}, error) {
	ret := make([]int, len(data))
	for i := range data {
		r, err := ConvertToInt(ctx, data[i:i+1])
		if err != nil {
			return nil, err
		}
		ret[i] = r.(int)
	}
	return ret, nil
}

// ConvertToFloat64Slice converts all elements in data to float64, and return []float64
func ConvertToFloat64Slice(ctx context.Context, data []string) (interface{}, error) {
	ret := make([]float64, len(data))
	for i := range data {
		r, err := ConvertToFloat64(ctx, data[i:i+1])
		if err != nil {
			return nil, err
		}
		ret[i] = r.(float64)
	}
	return ret, nil
}

// ConvertToStringSlice return all strings in data.
func ConvertToStringSlice(ctx context.Context, data []string) (interface{}, error) {
	return data, nil
}
