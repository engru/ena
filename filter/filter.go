// MIT License

// Copyright (c) 2018 soren yang

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package filter provide dynamic filter container to group filters
package filter

import (
	"reflect"

	"github.com/pkg/errors"
)

// Filter is the interface for do generic filter
type Filter interface {
	// Enabled return true if the filter will actually filter obj,
	// it return false if nothing FilterFunc have been registered.
	Enabled() bool

	// RegisterFilterFunc will append the fn and call it when Do is called.
	RegisterFilterFunc(fn FilterFunc)

	// Do will actually filter the object,
	// it return true if the obj should be keeped.
	Do(obj interface{}) bool

	// DoSlice will filter the array object, it will return an array object with the filtered object.
	// NOTE: the obj **MUST** be an array or ptr to array.
	//   1. If the obj is an array, return array
	//   2. If the obj is ptr to array, return ptr to array, and it always not return nil ptr
	DoSlice(obj interface{}) (interface{}, error)
}

// nolint
// FilterFunc is the object specific filter func,
// it return true if the obj should be keeped.
type FilterFunc func(obj interface{}) bool

// filterImpl implement the Filter interface
type filterImpl struct {
	fns []FilterFunc
}

func (f *filterImpl) Enabled() bool {
	return len(f.fns) != 0
}

func (f *filterImpl) RegisterFilterFunc(fn FilterFunc) {
	f.fns = append(f.fns, fn)
}

func (f *filterImpl) Do(obj interface{}) bool {
	for _, fn := range f.fns {
		if !fn(obj) {
			return false
		}
	}

	return true
}

func (f *filterImpl) DoSlice(obj interface{}) (interface{}, error) {
	if !f.Enabled() {
		return obj, nil
	}

	vin := reflect.ValueOf(obj)
	switch k := vin.Kind(); k {
	case reflect.Ptr:
		if vin.Elem().Kind() != reflect.Slice {
			return nil, errors.Errorf("expect ptr to slice, got ptr to %v", vin.Elem().Kind())
		}

		return f.doSlice(vin.Elem(), true)
	case reflect.Slice:
		return f.doSlice(vin, false)
	default:
		return nil, errors.Errorf("expect ptr or slice, got %v", k)
	}
}

func (f *filterImpl) doSlice(v reflect.Value, isPtr bool) (interface{}, error) {
	ret := reflect.MakeSlice(v.Type(), 0, v.Len()/2)

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if f.Do(item.Interface()) {
			ret = reflect.Append(ret, item)
		}
	}

	if !isPtr {
		return ret.Interface(), nil
	}

	ret1 := reflect.New(v.Type())
	ret1.Elem().Set(ret)
	return ret1.Interface(), nil
}

// NewFilter return and filter object
func NewFilter(fns ...FilterFunc) Filter {
	f := &filterImpl{
		fns: fns,
	}

	return f
}
