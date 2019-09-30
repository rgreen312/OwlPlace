package apiserver_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type handlerTest struct {
	name          string
	handler       func(http.ResponseWriter, *http.Request)
	req           *http.Request
	recorder      *httptest.ResponseRecorder
	expectWritten []byte
}

func (h *handlerTest) Verify(t *testing.T) {
	apiserver.Hello(h.recorder, h.req)
	recordedBytes := h.recorder.Body.Bytes()
	if bytes.Compare(recordedBytes, h.expectWritten) != 0 {
		t.Errorf("Expected: %s, Got: %s", h.expectWritten, recordedBytes)
	}
}

func TestHello(t *testing.T) {
	tests := []*handlerTest{
		&handlerTest{
			name:          "hello_test",
			handler:       apiserver.Hello,
			req:           &http.Request{},
			recorder:      httptest.NewRecorder(),
			expectWritten: []byte("hello\n"),
		},
	}

	for _, test := range tests {
		test.Verify(t)
	}
}

//func TestHeaders(t *testing.T) {
//want := "Hello, world."
//if got := apiserver.Hello(); got != want {
//t.Errorf("Hello() = %q, want %q", got, want)
//}
//}
