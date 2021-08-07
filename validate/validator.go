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

// Interface knowns how to validate one type
type Interface interface {
	// RegisterValidateFunc registers a function that validate the object.
	RegisterValidateFunc(in interface{}, fn ValidateFunc) error

	// MustRegisterValidateFunc registers a function that validate the object.
	// It will panic if err occurred.
	MustRegisterValidateFunc(in interface{}, fn ValidateFunc)

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

func (v *validatorImpl) RegisterValidateFunc(in interface{}, fn ValidateFunc) error {
	tIn := reflect.TypeOf(in)
	if tIn.Kind() != reflect.Ptr {
		return errors.Errorf("the type %T must be a pointer to register as an validate", in)
	}

	v.validateFuncs.untyped[typeInfo{
		in: tIn,
	}] = fn
	return nil
}

func (v *validatorImpl) MustRegisterValidateFunc(in interface{}, fn ValidateFunc) {
	err := v.RegisterValidateFunc(in, fn)
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

var (
	// DefaultValidator is the validate function for use
	DefaultValidator Interface
)

func init() {
	DefaultValidator = NewValidator()
}
