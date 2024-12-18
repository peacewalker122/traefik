package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/traefik/traefik/v2/pkg/testhelpers"
)

type staticTransport struct {
	res *http.Response
}

func (t *staticTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.res, nil
}

func BenchmarkProxy(b *testing.B) {
	res := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	w := httptest.NewRecorder()
	req := testhelpers.MustNewRequest(http.MethodGet, "http://foo.bar/", nil)

	pool := newBufferPool()
	handler, _ := buildProxy(pointer(false), nil, &staticTransport{res}, pool)

	b.ReportAllocs()
	for range b.N {
		handler.ServeHTTP(w, req)
	}
}
