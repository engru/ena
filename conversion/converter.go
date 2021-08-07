package conversion

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/lsytj0413/ena"
)

var (
	// ErrUnsupported defines that the generic conversion fn not support current args,
	//  then the converter will try the next generic conversion fn.
	ErrUnsupported = errors.New("ErrUnsupported")
)

// Scope is passed to conversion funcs to allow them to continue an ongoing conversion.
// If multiple converters exist in the system, Scope will allow you to use the correct one
// from a conversion function--that is, the one your conversion function was called by.
type Scope interface {
	// Call Convert to convert sub-objects. Note that if you call it with your own exact
	// parameters, you'll run out of stack space before anything useful happens.
	Convert(in interface{}, out interface{}) error

	// Meta returns any information originally passed to Convert.
	Meta() *Meta
}

// Meta is supplied by Scheme, when it calls Convert.
type Meta struct {
	// Context is an optional field that callers may use to pass info to conversion functions.
	Context interface{}
}

// scope contains information about an ongoing conversion.
type scope struct {
	converter Converter
	meta      *Meta
}

// Convert continues a conversion.
func (s *scope) Convert(in interface{}, out interface{}) error {
	return s.converter.Convert(in, out, s.meta.Context)
}

// Meta returns the meta object that was originally passed to Convert.
func (s *scope) Meta() *Meta {
	return s.meta
}

// nolint
// ConversionFunc converts the object a into the object b, reusing arrays or objects
// or pointers if necessary. It should return an error if the object cannot be converted
// or if some data is invalid. If you do not wish a and b to share fields or nested
// objects, you must copy a before calling this function.
type ConversionFunc func(in interface{}, out interface{}, scope Scope) error

// Converter knows how to convert one type to another.
type Converter interface {
	// Convert will attempt to convert in into out. Both must be pointers.
	Convert(in interface{}, out interface{}, context interface{}) error

	// RegisterConversionFunc registers a function that converts between a and b by passing objects of those
	// types to the provided function. The function *must* accept objects of a and b - this machinery will not enforce
	// any other guarantee.
	RegisterConversionFunc(in interface{}, out interface{}, fn ConversionFunc) error

	// RegisterGenericConversionFunc registers a function that convert objects.
	// It will been called if there is no typed conversion func exists.
	// NOTE: All generic conversion funcs will been called one by one, if the function returns ErrUnspport
	//    the next one will been called. Otherwise Convert will return the result.
	RegisterGenericConversionFunc(fn ConversionFunc)
}

type typePair struct {
	source reflect.Type
	dest   reflect.Type
}

// nolint
// ConversionFuncs defines the func map
type ConversionFuncs struct {
	untyped map[typePair]ConversionFunc
}

type converterImpl struct {
	conversionFuncs ConversionFuncs
	genericFuncs    []ConversionFunc
}

// NewConverter creates a new Converter object.
func NewConverter() Converter {
	c := &converterImpl{
		conversionFuncs: ConversionFuncs{
			untyped: make(map[typePair]ConversionFunc),
		},
	}
	return c
}

func (c *converterImpl) Convert(in interface{}, out interface{}, context interface{}) error {
	pair := typePair{reflect.TypeOf(in), reflect.TypeOf(out)}
	meta := &Meta{
		Context: context,
	}
	scope := &scope{
		converter: c,
		meta:      meta,
	}

	if fn, ok := c.conversionFuncs.untyped[pair]; ok {
		return fn(in, out, scope)
	}

	for _, fn := range c.genericFuncs {
		err := fn(in, out, scope)
		if err == nil {
			return nil
		}

		if !errors.Is(err, ErrUnsupported) {
			return err
		}
	}

	dv, err := ena.EnforcePtr(out)
	if err != nil {
		return err
	}
	sv, err := ena.EnforcePtr(in)
	if err != nil {
		return err
	}
	return errors.Errorf("converting (%s) to (%s): unknown conversion", sv.Type(), dv.Type())
}

func (c *converterImpl) RegisterConversionFunc(in interface{}, out interface{}, fn ConversionFunc) error {
	tA, tB := reflect.TypeOf(in), reflect.TypeOf(out)
	if tA.Kind() != reflect.Ptr {
		return errors.Errorf("the type %T must be a pointer to register as an conversion", in)
	}
	if tB.Kind() != reflect.Ptr {
		return errors.Errorf("the type %T must be a pointer to register as an conversion", out)
	}
	c.conversionFuncs.untyped[typePair{tA, tB}] = fn
	return nil
}

func (c *converterImpl) RegisterGenericConversionFunc(fn ConversionFunc) {
	c.genericFuncs = append(c.genericFuncs, fn)
}

var (
	// DefaultConverter is the convert function for use
	DefaultConverter Converter
)

func init() {
	DefaultConverter = NewConverter()
}
