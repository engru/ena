package validate

import (
	"fmt"
	"testing"
)

func TestRegisterValidateFunc(t *testing.T) {
	type testCase struct {
		desp string

		in   interface{}
		opts []Option
		fn   ValidateFunc

		isErr bool
	}
	testCases := []testCase{
		{
			desp:  "normal test",
			in:    &testCase{},
			fn:    nil,
			isErr: false,
		},
		{
			desp: "register with slice",
			in:   &testCase{},
			opts: []Option{
				WithSliceValidateFunc(true),
			},
			fn:    nil,
			isErr: false,
		},
		{
			desp: "register with ptr slice",
			in:   &testCase{},
			opts: []Option{
				WithPtrSliceValidateFunc(true),
			},
			fn:    nil,
			isErr: false,
		},
		{
			desp: "register with both slice & ptr slice",
			in:   &testCase{},
			opts: []Option{
				WithPtrSliceValidateFunc(true),
				WithSliceValidateFunc(true),
			},
			fn:    nil,
			isErr: false,
		},
		{
			desp:  "in not ptr",
			in:    testCase{},
			fn:    nil,
			isErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			c := NewValidator()

			err := c.RegisterValidateFunc(tc.in, tc.fn, tc.opts...)
			if tc.isErr != (err != nil) {
				t.Fatalf("Expect err %v, got %v", tc.isErr, err)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	type testValidate1 struct {
		Total1 int
	}
	var fn = func(in *testValidate1, _ Scope) error { //nolint
		if in.Total1 <= 0 {
			return fmt.Errorf("validate failed")
		}

		return nil
	}

	t.Run("normal test", func(t *testing.T) {
		c := NewValidator()
		c.RegisterValidateFunc((*testValidate1)(nil), func(in interface{}, scope Scope) error {
			return fn(in.(*testValidate1), scope)
		})

		in := &testValidate1{
			Total1: 1,
		}

		err := c.Validate(in, nil)
		if err != nil {
			t.Fatalf("Expect no err, go %v", err)
		}
	})

	t.Run("not register", func(t *testing.T) {
		c := NewValidator()

		in := &testValidate1{
			Total1: 1,
		}
		err := c.Validate(in, nil)
		if err == nil {
			t.Fatalf("Expect err, go no err")
		}
	})

	t.Run("in not ptr", func(t *testing.T) {
		c := NewValidator()
		c.RegisterValidateFunc((*testValidate1)(nil), func(in interface{}, scope Scope) error {
			return fn(in.(*testValidate1), scope)
		})

		in := testValidate1{
			Total1: 1,
		}

		err := c.Validate(in, nil)
		if err == nil {
			t.Fatalf("Expect err, go no err")
		}
	})

	t.Run("continue a validate with scope", func(t *testing.T) {
		fn1 := func(in *testValidate1, scope Scope) error { //nolint
			return scope.Validate(&(in.Total1))
		}

		fn2 := func(in *int, scope Scope) error {
			if *in <= 0 {
				return fmt.Errorf("validate failed")
			}

			ctx, ok := scope.Meta().Context.(string)
			if ok && ctx == "test-context" {
				return fmt.Errorf("err scope context")
			}
			return nil
		}

		c := NewValidator()
		c.RegisterValidateFunc((*testValidate1)(nil), func(in interface{}, scope Scope) error {
			return fn1(in.(*testValidate1), scope)
		})
		c.RegisterValidateFunc((*int)(nil), func(in interface{}, scope Scope) error {
			return fn2(in.(*int), scope)
		})

		in := &testValidate1{
			Total1: 1,
		}

		err := c.Validate(in, nil)
		if err != nil {
			t.Fatalf("Expect no err, got %v", err)
		}

		err = c.Validate(in, "test-context")
		if err == nil {
			t.Fatalf("Expect err, got no err")
		}

		in.Total1 = -1
		err = c.Validate(in, nil)
		if err == nil {
			t.Fatalf("Expect err, got no err")
		}
	})

	t.Run("with slice", func(t *testing.T) {
		c := NewValidator()
		c.RegisterValidateFunc((*testValidate1)(nil), func(in interface{}, scope Scope) error {
			return fn(in.(*testValidate1), scope)
		}, WithSliceValidateFunc(true))

		in := &[]testValidate1{
			{
				Total1: 1,
			},
		}

		err := c.Validate(in, nil)
		if err != nil {
			t.Fatalf("Expect no err, got %v", err)
		}
	})

	t.Run("with ptr slice", func(t *testing.T) {
		c := NewValidator()
		c.RegisterValidateFunc((*testValidate1)(nil), func(in interface{}, scope Scope) error {
			return fn(in.(*testValidate1), scope)
		}, WithPtrSliceValidateFunc(true))

		in := &[]*testValidate1{
			{
				Total1: 1,
			},
		}

		err := c.Validate(in, nil)
		if err != nil {
			t.Fatalf("Expect no err, got %v", err)
		}
	})
}
