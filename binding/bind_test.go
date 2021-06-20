package binding

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/gomega"
)

func TestBind(t *testing.T) {
	t.Run("Header test", func(t *testing.T) {
		g := NewWithT(t)
		type Out struct {
			Value string `bind:"Header,value"`
		}
		o := &Out{}
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Set("value", "vvv")

		g.Expect(Bind(context.Background(), WrapHTTPRequest(req), o)).ToNot(HaveOccurred())
		g.Expect(o).To(Equal(&Out{
			Value: "vvv",
		}))
	})

	t.Run("Query test", func(t *testing.T) {
		g := NewWithT(t)
		type Out struct {
			Value string `bind:"Query,value"`
		}
		o := &Out{}
		req := httptest.NewRequest("GET", "/api?value=vvv", nil)

		g.Expect(Bind(context.Background(), WrapHTTPRequest(req), o)).ToNot(HaveOccurred())
		g.Expect(o).To(Equal(&Out{
			Value: "vvv",
		}))
	})

	t.Run("Body test", func(t *testing.T) {
		g := NewWithT(t)

		out := struct {
			Value1 struct {
				Value string `json:"value"`
			} `bind:"Body"`
		}{}

		req := httptest.NewRequest("GET", "/api?value=vvv", bytes.NewReader([]byte(`{"value":"test"}`)))

		g.Expect(Bind(context.Background(), WrapHTTPRequest(req), &out)).ToNot(HaveOccurred())
		g.Expect(out.Value1.Value).To(Equal("test"))
	})

	t.Run("Body ptr test", func(t *testing.T) {
		g := NewWithT(t)

		out := struct {
			Value1 *struct {
				Value string `json:"value"`
			} `bind:"Body"`
		}{}
		req := httptest.NewRequest("GET", "/api?value=vvv", bytes.NewReader([]byte(`{"value":"test"}`)))

		g.Expect(Bind(context.Background(), WrapHTTPRequest(req), &out)).ToNot(HaveOccurred())
		g.Expect(out.Value1.Value).To(Equal("test"))
	})

	t.Run("Auto test", func(t *testing.T) {
		g := NewWithT(t)

		out := struct {
			Value1 struct {
				Value string `bind:"Query,value"`
			} `bind:"Auto"`
		}{}

		req := httptest.NewRequest("GET", "/api?value=vvv", nil)

		g.Expect(Bind(context.Background(), WrapHTTPRequest(req), &out)).ToNot(HaveOccurred())
		g.Expect(out.Value1.Value).To(Equal("vvv"))
	})

	t.Run("Auto ptr test", func(t *testing.T) {
		g := NewWithT(t)

		out := struct {
			Value1 *struct {
				Value string `bind:"Query,value"`
			} `bind:"Auto"`
		}{}

		req := httptest.NewRequest("GET", "/api?value=vvv", nil)

		g.Expect(Bind(context.Background(), WrapHTTPRequest(req), &out)).ToNot(HaveOccurred())
		g.Expect(out.Value1.Value).To(Equal("vvv"))
	})
}
