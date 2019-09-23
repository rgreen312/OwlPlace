package server_test

import (
	"testing"

	"github.com/rgreen312/owlplace/server"
)

func TestHello(t *testing.T) {
	want := "Hello, world."
	if got := server.Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
