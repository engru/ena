package helper

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestConvert(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		in := map[string]string{
			"k": "v",
		}
		var out string
		g.Expect(Convert_map_string_string_to_string(&in, &out, nil)).ToNot(HaveOccurred())
		g.Expect(out).To(Equal(`{"k":"v"}`))
	})

	t.Run("nil test", func(t *testing.T) {
		g := NewWithT(t)

		in := map[string]string{}
		var out string
		g.Expect(Convert_map_string_string_to_string(&in, &out, nil)).ToNot(HaveOccurred())
		g.Expect(out).To(Equal(``))
	})
}
