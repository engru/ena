package filter

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestDoFilter(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		f := NewFilter()
		g.Expect(f.Enabled()).To(BeFalse())
		g.Expect(f.Do(1)).To(BeTrue())

		f.RegisterFilterFunc(func(obj interface{}) bool {
			o := obj.(int)
			return o > 3 && o < 10
		})
		g.Expect(f.Enabled()).To(BeTrue())
		g.Expect(f.Do(1)).To(BeFalse())
		g.Expect(f.Do(5)).To(BeTrue())
		g.Expect(f.Do(8)).To(BeTrue())

		f.RegisterFilterFunc(func(obj interface{}) bool {
			o := obj.(int)
			return o > 7 && o < 13
		})
		g.Expect(f.Enabled()).To(BeTrue())
		g.Expect(f.Do(1)).To(BeFalse())
		g.Expect(f.Do(5)).To(BeFalse())
		g.Expect(f.Do(8)).To(BeTrue())
	})
}

func TestDoArray(t *testing.T) {
	t.Run("normal slice test", func(t *testing.T) {
		g := NewWithT(t)

		items := []int{1, 2, 3, 4, 5}
		f := NewFilter(func(obj interface{}) bool {
			v := obj.(int)
			return v > 2 && v <= 4
		})
		v, err := f.DoSlice(items)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v.([]int)).To(Equal([]int{3, 4}))
	})

	t.Run("normal ptr slice test", func(t *testing.T) {
		g := NewWithT(t)

		items := []int{1, 2, 3, 4, 5}
		f := NewFilter(func(obj interface{}) bool {
			v := obj.(int)
			return v > 2 && v <= 4
		})
		v, err := f.DoSlice(&items)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v.(*[]int)).To(Equal(&[]int{3, 4}))
	})

	t.Run("normal ptr slice filter none test", func(t *testing.T) {
		g := NewWithT(t)

		items := []int{1, 2, 3, 4, 5}
		f := NewFilter(func(obj interface{}) bool {
			v := obj.(int)
			return v > 100
		})
		v, err := f.DoSlice(&items)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v).ToNot(BeNil())
		g.Expect(v.(*[]int)).To(Equal(&[]int{}))
	})

	t.Run("normal disable test", func(t *testing.T) {
		g := NewWithT(t)

		items := []int{1, 2, 3, 4, 5}
		f := NewFilter()
		v, err := f.DoSlice(&items)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v.(*[]int)).To(Equal(&items))
	})

	t.Run("not slice test", func(t *testing.T) {
		g := NewWithT(t)

		items := 1
		_, err := NewFilter(func(obj interface{}) bool {
			v := obj.(int)
			return v > 2 && v <= 4
		}).DoSlice(items)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("expect ptr or slice"))
	})

	t.Run("not ptr to slice test", func(t *testing.T) {
		g := NewWithT(t)

		items := 1
		_, err := NewFilter(func(obj interface{}) bool {
			v := obj.(int)
			return v > 2 && v <= 4
		}).DoSlice(&items)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("expect ptr to slice"))
	})
}
