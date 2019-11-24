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
	// apiserver.Hello(h.recorder, h.req)
	recordedBytes := h.recorder.Body.Bytes()
	if bytes.Compare(recordedBytes, h.expectWritten) != 0 {
		t.Errorf("Expected: %s, Got: %s", h.expectWritten, recordedBytes)
	}
}
