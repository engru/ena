package ena

import (
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
