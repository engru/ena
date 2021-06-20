package conversion

import (
	"fmt"
	"testing"
)

func TestEnforcePtr(t *testing.T) {
	type testCase struct {
		desp string
		obj  interface{}

		isErr bool
	}
	testCases := []testCase{
		{
			desp:  "normal test",
			obj:   &testCase{},
			isErr: false,
		},
		{
			desp:  "not pointer",
			obj:   testCase{},
			isErr: true,
		},
		{
			desp:  "invalid kind",
			obj:   nil,
			isErr: true,
		},
		{
			desp:  "nil value",
			obj:   (*testCase)(nil),
			isErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			_, err := EnforcePtr(tc.obj)
			if tc.isErr != (err != nil) {
				t.Fatalf("Expect err %v, got %v", tc.isErr, err)
			}
		})
	}
}

func TestRegisterGenericConversionFunc(t *testing.T) {
	c := NewConverter()

	var fn ConversionFunc

	c.RegisterGenericConversionFunc(fn)

	ci := c.(*converterImpl)
	if len(ci.genericFuncs) != 1 {
		t.Fatalf("Expect generic funcs len 1, got %v", len(ci.genericFuncs))
	}
}

func TestRegisterConversionFunc(t *testing.T) {
	type testCase struct {
		desp string

		in  interface{}
		out interface{}
		fn  ConversionFunc

		isErr bool
	}
	testCases := []testCase{
		{
			desp: "normal test",
			in:   &testCase{},
			out:  &testCase{},
			fn:   nil,

			isErr: false,
		},
		{
			desp: "in not ptr",
			in:   testCase{},
			out:  &testCase{},
			fn:   nil,

			isErr: true,
		},
		{
			desp: "out not ptr",
			in:   &testCase{},
			out:  testCase{},
			fn:   nil,

			isErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			c := NewConverter()

			err := c.RegisterConversionFunc(tc.in, tc.out, tc.fn)
			if tc.isErr != (err != nil) {
				t.Fatalf("Expect err %v, got %v", tc.isErr, err)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	type testConvert1 struct {
		Total1 int
	}
	type testConvert2 struct {
		Total2 int
	}
	var fn = func(in *testConvert1, out *testConvert2, _ Scope) error { //nolint
		out.Total2 = in.Total1
		return nil
	}

	t.Run("normal test", func(t *testing.T) {
		c := NewConverter()
		c.RegisterConversionFunc((*testConvert1)(nil), (*testConvert2)(nil), func(in interface{}, out interface{}, scope Scope) error {
			return fn(in.(*testConvert1), out.(*testConvert2), scope)
		})

		in := &testConvert1{
			Total1: 1,
		}
		out := &testConvert2{
			Total2: 2,
		}
		err := c.Convert(in, out, nil)
		if err != nil {
			t.Fatalf("Expect no err, go %v", err)
		}

		if out.Total2 != in.Total1 {
			t.Fatalf("Expect out %v, got %v", in.Total1, out.Total2)
		}
	})

	t.Run("not register", func(t *testing.T) {
		c := NewConverter()

		in := &testConvert1{
			Total1: 1,
		}
		out := &testConvert2{
			Total2: 2,
		}
		err := c.Convert(in, out, nil)
		if err == nil {
			t.Fatalf("Expect err, go no err")
		}

		if out.Total2 != 2 {
			t.Fatalf("Expect out %v, got %v", 2, out.Total2)
		}
	})

	t.Run("in not ptr", func(t *testing.T) {
		c := NewConverter()
		c.RegisterConversionFunc((*testConvert1)(nil), (*testConvert2)(nil), func(in interface{}, out interface{}, scope Scope) error {
			return fn(in.(*testConvert1), out.(*testConvert2), scope)
		})

		in := testConvert1{
			Total1: 1,
		}
		out := &testConvert2{
			Total2: 2,
		}
		err := c.Convert(in, out, nil)
		if err == nil {
			t.Fatalf("Expect err, go no err")
		}

		if out.Total2 != 2 {
			t.Fatalf("Expect out %v, got %v", 2, out.Total2)
		}
	})

	t.Run("out not ptr", func(t *testing.T) {
		c := NewConverter()
		c.RegisterConversionFunc((*testConvert1)(nil), (*testConvert2)(nil), func(in interface{}, out interface{}, scope Scope) error {
			return fn(in.(*testConvert1), out.(*testConvert2), scope)
		})

		in := &testConvert1{
			Total1: 1,
		}
		out := testConvert2{
			Total2: 2,
		}
		err := c.Convert(in, out, nil)
		if err == nil {
			t.Fatalf("Expect err, go no err")
		}

		if out.Total2 != 2 {
			t.Fatalf("Expect out %v, got %v", 2, out.Total2)
		}
	})

	t.Run("with generic funcs converted unsupported", func(t *testing.T) {
		c := NewConverter()
		c.RegisterGenericConversionFunc(func(in interface{}, out interface{}, scope Scope) error {
			return ErrUnsupported
		})

		in := &testConvert1{
			Total1: 1,
		}
		out := &testConvert2{
			Total2: 2,
		}
		err := c.Convert(in, out, nil)
		if err == nil {
			t.Fatalf("Expect err, got no err")
		}

		if out.Total2 != 2 {
			t.Fatalf("Expect out %v, got %v", 2, out.Total2)
		}
	})

	t.Run("with generic funcs converted err", func(t *testing.T) {
		c := NewConverter()
		c.RegisterGenericConversionFunc(func(in interface{}, out interface{}, scope Scope) error {
			return fmt.Errorf("tested")
		})

		in := &testConvert1{
			Total1: 1,
		}
		out := &testConvert2{
			Total2: 2,
		}
		err := c.Convert(in, out, nil)
		if err == nil {
			t.Fatalf("Expect err, got no err")
		}

		if out.Total2 != 2 {
			t.Fatalf("Expect out %v, got %v", 2, out.Total2)
		}
	})

	t.Run("with generic funcs convert supported", func(t *testing.T) {
		c := NewConverter()
		c.RegisterGenericConversionFunc(func(in interface{}, out interface{}, scope Scope) error {
			return fn(in.(*testConvert1), out.(*testConvert2), scope)
		})

		in := &testConvert1{
			Total1: 1,
		}
		out := &testConvert2{
			Total2: 2,
		}
		err := c.Convert(in, out, nil)
		if err != nil {
			t.Fatalf("Expect no err, go %v", err)
		}

		if out.Total2 != in.Total1 {
			t.Fatalf("Expect out %v, got %v", in.Total1, out.Total2)
		}
	})

	t.Run("with second generic funcs convert supported", func(t *testing.T) {
		c := NewConverter()
		c.RegisterGenericConversionFunc(func(in interface{}, out interface{}, scope Scope) error {
			return ErrUnsupported
		})
		c.RegisterGenericConversionFunc(func(in interface{}, out interface{}, scope Scope) error {
			return fn(in.(*testConvert1), out.(*testConvert2), scope)
		})

		in := &testConvert1{
			Total1: 1,
		}
		out := &testConvert2{
			Total2: 2,
		}
		err := c.Convert(in, out, nil)
		if err != nil {
			t.Fatalf("Expect no err, go %v", err)
		}

		if out.Total2 != in.Total1 {
			t.Fatalf("Expect out %v, got %v", in.Total1, out.Total2)
		}
	})

	t.Run("continue a conversion with scope", func(t *testing.T) {
		fn1 := func(in *testConvert1, out *testConvert2, scope Scope) error { //nolint
			return scope.Convert(&(in.Total1), &(out.Total2))
		}

		fn2 := func(in *int, out *int, scope Scope) error {
			*out = *in
			ctx, ok := scope.Meta().Context.(string)
			if !ok || ctx != "test-context" {
				return fmt.Errorf("err scope context")
			}
			return nil
		}

		c := NewConverter()
		c.RegisterConversionFunc((*testConvert1)(nil), (*testConvert2)(nil), func(in interface{}, out interface{}, scope Scope) error {
			return fn1(in.(*testConvert1), out.(*testConvert2), scope)
		})
		c.RegisterConversionFunc((*int)(nil), (*int)(nil), func(in interface{}, out interface{}, scope Scope) error {
			return fn2(in.(*int), out.(*int), scope)
		})

		in := &testConvert1{
			Total1: 1,
		}
		out := &testConvert2{
			Total2: 2,
		}
		err := c.Convert(in, out, "test-context")
		if err != nil {
			t.Fatalf("Expect no err, got %v", err)
		}

		if out.Total2 != in.Total1 {
			t.Fatalf("Expect out %v, got %v", in.Total1, out.Total2)
		}

	})
}
