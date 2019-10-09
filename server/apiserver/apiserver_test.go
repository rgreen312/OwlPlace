package apiserver_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rgreen312/owlplace/server/apiserver"
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

func TestUpdate(t *testing.T) {
	// api_to_backend_channel := make(chan consensus.BackendMessage)
	// backend_to_api_channel := make(chan consensus.ConsensusMessage)
	// apiserver := apiserver.NewApiServer(api_to_backend_channel, backend_to_api_channel)

	// var dat map[string]interface{}
	renameMe := apiserver.DrawPixelMsg{X: 12, Y: 12, R: 255, G: 255, B: 255, UserID: "testId"}
	// dat["x"] = 12
	// dat["y"] = 12
	// dat["r"] = 255
	// dat["g"] = 255
	// dat["b"] = 255
	result := apiserver.updateMethod(renameMe)
	expected := "put pixel(12,12) (255,255,255,255)"
	if result != expected {
		t.Errorf("Sum was incorrect, got: %f, want: %f.", result, expected)
	}
}

//func TestHeaders(t *testing.T) {
//want := "Hello, world."
//if got := apiserver.Hello(); got != want {
//t.Errorf("Hello() = %q, want %q", got, want)
//}
//}
