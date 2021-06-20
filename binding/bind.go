package binding

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/pkg/errors"
)

// Request present the input value
type Request interface {
	// Header return the content from http header
	Header(name string) ([]string, error)

	// Query return the content from http query
	Query(name string) ([]string, error)

	// Body return the content from http body
	Body() ([]byte, error)

	// ContentType return the body content-type
	ContentType() string

	// Path return the path variable
	Path(name string) (string, error)
}

type httpRequest struct {
	req   *http.Request
	paths map[string]string
}

// WrapHTTPRequest return the Request instance
func WrapHTTPRequest(req *http.Request) Request {
	return &httpRequest{
		req:   req,
		paths: make(map[string]string),
	}
}

func (h *httpRequest) Header(name string) ([]string, error) {
	return h.req.Header.Values(name), nil
}

func (h *httpRequest) Query(name string) ([]string, error) {
	return h.req.URL.Query()[name], nil
}

func (h *httpRequest) Body() ([]byte, error) {
	return ioutil.ReadAll(h.req.Body)
}

func (h *httpRequest) ContentType() string {
	return h.req.Header.Get("content-type")
}

func (h *httpRequest) Path(name string) (string, error) {
	return h.paths[name], nil
}

// Binder is the interface for bind request param to struct
type Binder interface {
	// Bind will auto-unpack the param from request to struct.
	// 1. the value must been an ptr to struct
	Bind(ctx context.Context, req Request, value interface{}) error
}

var (
	defaultBinder Binder = &bindImpl{}
)

// Bind will auto-unpack the param from request to struct.
func Bind(ctx context.Context, req Request, value interface{}) error {
	return defaultBinder.Bind(ctx, req, value)
}

type bindImpl struct {
}

func (b *bindImpl) dereferencePtrStruct(value interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr {
		return val, errors.Errorf("Expect pointer to struct, but got %v", val.Kind())
	}
	if val.IsNil() {
		return val, errors.Errorf("Expect pointer to struct, but got nil pointer")
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return val, errors.Errorf("Expect pointer to struct, but got ptr to %v", val.Kind())
	}
	return val, nil
}

func (b *bindImpl) enum(val reflect.Value, req Request, index []int, typ reflect.Type) error {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		err := b.bindField(val, req, append(index, i), field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *bindImpl) bindField(val reflect.Value, req Request, index []int, field reflect.StructField) error {
	tag := field.Tag.Get("bind")
	if tag == "" {
		return nil
	}

	bTag, err := parseBindParameterTag(tag)
	if err != nil {
		return errors.Wrapf(err, "Field %v", field.Name)
	}

	switch bTag.source {
	case "Path":
		converter := ConverterFor(field.Type)
		if converter == nil {
			return errors.Errorf("Unknown converter for %T", field.Type)
		}
		value, err := req.Path(bTag.Name())
		if err != nil {
			return errors.Wrapf(err, "Cann't retrieve path %v", bTag.Name())
		}
		if len(value) == 0 && bTag.defaultValue != nil {
			value = *bTag.defaultValue
		}

		v, err := converter(context.TODO(), []string{value})
		if err != nil {
			return errors.Wrapf(err, "Convert for %T with value %v", field.Type, value)
		}
		val.FieldByIndex(index).Set(reflect.ValueOf(v))
	case "Header":
		converter := ConverterFor(field.Type)
		if converter == nil {
			return errors.Errorf("Unknown converter for %T", field.Type)
		}
		value, err := req.Header(bTag.Name())
		if err != nil {
			return errors.Wrapf(err, "Cann't retrieve header %v", bTag.Name())
		}
		if len(value) == 0 && bTag.defaultValue != nil {
			value = append(value, *bTag.defaultValue)
		}

		v, err := converter(context.TODO(), value)
		if err != nil {
			return errors.Wrapf(err, "Convert for %T with value %v", field.Type, value)
		}
		val.FieldByIndex(index).Set(reflect.ValueOf(v))
	case "Query":
		converter := ConverterFor(field.Type)
		if converter == nil {
			return errors.Errorf("Unknown converter for %T", field.Type)
		}
		value, err := req.Query(bTag.Name())
		if err != nil {
			return errors.Wrapf(err, "Cann't retrieve header %v", bTag.Name())
		}
		if len(value) == 0 && bTag.defaultValue != nil {
			value = append(value, *bTag.defaultValue)
		}

		v, err := converter(context.TODO(), value)
		if err != nil {
			return errors.Wrapf(err, "Convert for %T with value %v", field.Type, value)
		}
		val.FieldByIndex(index).Set(reflect.ValueOf(v))
	case "Auto":
		var err error
		if field.Type.Kind() == reflect.Ptr {
			if val.FieldByIndex(index).IsNil() {
				val.FieldByIndex(index).Set(reflect.ValueOf(reflect.New(field.Type.Elem()).Interface()))
			}

			err = b.enum(val, req, index, field.Type.Elem())
		} else {
			err = b.enum(val, req, index, field.Type)
		}
		if err != nil {
			return err
		}
	case "Body":
		b, err := req.Body()
		if err != nil {
			return err
		}

		value := reflect.New(field.Type)
		err = json.Unmarshal(b, value.Interface())
		if err != nil {
			return err
		}
		val.FieldByIndex(index).Set(reflect.ValueOf(value.Elem().Interface()))
	default:
		return errors.Errorf("Unknown source %v", bTag.source)
	}

	return nil
}

func (b *bindImpl) Bind(ctx context.Context, req Request, value interface{}) error {
	val, err := b.dereferencePtrStruct(value)
	if err != nil {
		return err
	}

	return b.enum(val, req, []int{}, val.Type())
}
