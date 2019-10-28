package apiserver_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rgreen312/owlplace/server/apiserver"
	"github.com/rgreen312/owlplace/server/consensus"
)

type handlerTest struct {
	name          string
	handler       func(http.ResponseWriter, *http.Request)
	req           *http.Request
	recorder      *httptest.ResponseRecorder
	expectWritten []byte
}

func (h *handlerTest) Verify(t *testing.T) {
	// apiserver.Hello(h.recorder, h.req)
	recordedBytes := h.recorder.Body.Bytes()
	if bytes.Compare(recordedBytes, h.expectWritten) != 0 {
		t.Errorf("Expected: %s, Got: %s", h.expectWritten, recordedBytes)
	}
}

// func TestHello(t *testing.T) {
// 	tests := []*handlerTest{
// 		&handlerTest{
// 			name:          "hello_test",
// 			handler:       apiserver.Hello,
// 			req:           &http.Request{},
// 			recorder:      httptest.NewRecorder(),
// 			expectWritten: []byte("hello\n"),
// 		},
// 	}

// 	for _, test := range tests {
// 		test.Verify(t)
// 	}
// }

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestUpdate(t *testing.T) {
	fmt.Println("Starting TestUpdate")
	api_to_backend_channel := make(chan consensus.BackendMessage)
	backend_to_api_channel := make(chan consensus.ConsensusMessage)
	api := apiserver.NewApiServer(api_to_backend_channel, backend_to_api_channel)

	// var dat map[string]interface{}
	// renameMe := apiserver.DrawPixelMsg{X: 12, Y: 12, R: 255, G: 255, B: 255, UserID: "testId"}
	fmt.Println("Before calling the UpdateMethod")
	result := api.UpdateMethod(12, 12, 255, 255, 255, "user1")
	fmt.Println("Got result")
	expected := []byte("put pixel(12,12) (255,255,255,255)")
	result = append(result, 0)

	if Equal(result, expected) {
		fmt.Println("Success!")
	} else {
		s := string(result[:])
		t.Errorf("Sum was incorrect, got: %s, want: %s", s, expected)
	}
}

//func TestHeaders(t *testing.T) {
//want := "Hello, world."
//if got := apiserver.Hello(); got != want {
//t.Errorf("Hello() = %q, want %q", got, want)
//}
//}
