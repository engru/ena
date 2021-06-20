package binding

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func TestConvertToBool(t *testing.T) {
	type testCase struct {
		Desp   string
		Arg1   []string
		IsErr  bool
		Expect interface{}
	}
	testCases := []testCase{
		{
			Desp:   "Normal 1 test",
			Arg1:   []string{"1"},
			IsErr:  false,
			Expect: true,
		},
		{
			Desp:   "Failed test",
			Arg1:   []string{"2"},
			IsErr:  true,
			Expect: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Desp, func(t *testing.T) {
			g := NewWithT(t)
			v, err := ConvertToBool(context.Background(), tc.Arg1)
			if tc.IsErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(v).To(Equal(tc.Expect))
		})
	}
}
