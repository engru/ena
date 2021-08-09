package validate

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/lsytj0413/ena"
)

// Scope is passed to validate funcs to allow them to continue an ongoing validate.
// If multiple validators exist in the system, Scope will allow you to use the correct one
// from a validate function--that is, the one your validate function was called by.
type Scope interface {
	// Call Validate to valiate sub-objects. Note that if you call it with your own exact
	// parameters, you'll run out of stack space before anything useful happens.
	Validate(in interface{}) error

	// Meta returns any information originally passed to Validate.
	Meta() *Meta
}

// Meta is supplied by Scheme, when it calls Validate.
type Meta struct {
	// Context is an optional field that callers may use to pass info to validate functions.
	Context interface{}
}

// scope contains information about an ongoing validate.
type scope struct {
	validator Interface
	meta      *Meta
}

// Validate continues a validate.
func (s *scope) Validate(in interface{}) error {
	return s.validator.Validate(in, s.meta.Context)
}

// Meta returns the meta object that was originally passed to Validate.
func (s *scope) Meta() *Meta {
	return s.meta
}

// nolint
// ValidateFunc validate the object
type ValidateFunc = func(in interface{}, scope Scope) error

type option struct {
	// WithSlice will register the **auto generate** []Type ValidateFunc
	WithSlice bool

	// WithPtrSlice will register the **auto generate** []*Type ValidateFunc
	WithPtrSlice bool
}

// Option is some configuration that modifies options for a register.
type Option interface {
	Apply(*option)
}

// WithSliceValidateFunc set the WithSlice field
type WithSliceValidateFunc bool

// Apply applies this configuration to the given option
func (w WithSliceValidateFunc) Apply(opt *option) {
	opt.WithSlice = bool(w)
}

// WithPtrSliceValidateFunc set the WithPtrSlice field
type WithPtrSliceValidateFunc bool

// Apply applies this configuration to the given option
func (w WithPtrSliceValidateFunc) Apply(opt *option) {
	opt.WithPtrSlice = bool(w)
}

// Interface knowns how to validate one type
type Interface interface {
	// RegisterValidateFunc registers a function that validate the object.
	RegisterValidateFunc(in interface{}, fn ValidateFunc, opts ...Option) error

	// MustRegisterValidateFunc registers a function that validate the object.
	// It will panic if err occurred.
	MustRegisterValidateFunc(in interface{}, fn ValidateFunc, opts ...Option)

	// Validate will validate the in object. in must be pointers.
	Validate(in interface{}, context interface{}) error
}

type typeInfo struct {
	in reflect.Type
}

type validateFuncs struct {
	untyped map[typeInfo]ValidateFunc
}

type validatorImpl struct {
	validateFuncs validateFuncs
}

// NewValidator return the validate impl
func NewValidator() Interface {
	return &validatorImpl{
		validateFuncs: validateFuncs{
			untyped: make(map[typeInfo]ValidateFunc),
		},
	}
}

func (v *validatorImpl) RegisterValidateFunc(in interface{}, fn ValidateFunc, opts ...Option) error {
	tIn := reflect.TypeOf(in)
	if tIn.Kind() != reflect.Ptr {
		return errors.Errorf("the type %T must be a pointer to register as an validate", in)
	}

	options := &option{}
	for _, opt := range opts {
		opt.Apply(options)
	}

	v.validateFuncs.untyped[typeInfo{
		in: tIn,
	}] = fn

	if options.WithSlice {
		arrType := typeInfo{
			in: reflect.PtrTo(reflect.SliceOf(tIn.Elem())),
		}
		v.validateFuncs.untyped[arrType] = validateSlice
	}

	if options.WithPtrSlice {
		arrType := typeInfo{
			in: reflect.PtrTo(reflect.SliceOf(tIn)),
		}
		v.validateFuncs.untyped[arrType] = validatePtrSlice
	}
	return nil
}

func (v *validatorImpl) MustRegisterValidateFunc(in interface{}, fn ValidateFunc, opts ...Option) {
	err := v.RegisterValidateFunc(in, fn, opts...)
	if err != nil {
		panic(errors.Wrapf(err, "MustRegisterValidateFunc failed"))
	}
}

func (v *validatorImpl) Validate(in interface{}, context interface{}) error {
	tInfo := typeInfo{
		in: reflect.TypeOf(in),
	}
	meta := &Meta{
		Context: context,
	}
	scope := &scope{
		validator: v,
		meta:      meta,
	}

	if fn, ok := v.validateFuncs.untyped[tInfo]; ok {
		return fn(in, scope)
	}

	sv, err := ena.EnforcePtr(in)
	if err != nil {
		return err
	}
	return errors.Errorf("validate (%s): unknown validator", sv.Type())
}

func validateSlice(in interface{}, scope Scope) error {
	tIn := reflect.TypeOf(in)
	if tIn.Kind() != reflect.Ptr {
		panic(errors.Errorf("the type %T must be a pointer", in))
	}
	if tIn.Elem().Kind() != reflect.Slice {
		panic(errors.Errorf("the type %T must be a pointer to slice", in))
	}

	vIn := reflect.ValueOf(in)
	for i := 0; i < vIn.Elem().Len(); i++ {
		if err := scope.Validate(vIn.Elem().Index(i).Addr().Interface()); err != nil {
			return err
		}
	}

	return nil
}

func validatePtrSlice(in interface{}, scope Scope) error {
	tIn := reflect.TypeOf(in)
	if tIn.Kind() != reflect.Ptr {
		panic(errors.Errorf("the type %T must be a pointer", in))
	}
	if tIn.Elem().Kind() != reflect.Slice {
		panic(errors.Errorf("the type %T must be a pointer to slice", in))
	}
	if tIn.Elem().Elem().Kind() != reflect.Ptr {
		panic(errors.Errorf("the type %T must be a pointer to slice with pointer", in))
	}

	vIn := reflect.ValueOf(in)
	for i := 0; i < vIn.Elem().Len(); i++ {
		if err := scope.Validate(vIn.Elem().Index(i).Interface()); err != nil {
			return err
		}
	}

	return nil
}

var (
	// DefaultValidator is the validate function for use
	DefaultValidator Interface
)

func init() {
	DefaultValidator = NewValidator()
}
