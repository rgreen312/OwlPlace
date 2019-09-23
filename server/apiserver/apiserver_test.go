package apiserver_test

import (
	"testing"

	"github.com/rgreen312/owlplace/server/apiserver"
)

func TestHello(t *testing.T) {
	want := "Hello, world."
	if got := apiserver.Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
