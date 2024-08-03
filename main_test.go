package main

import (
	"testing"
)

func SpliceIsSame[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParseInput(t *testing.T) {
	v := NewTerminalView()
	m := NewStateKeeper()
	c := NewController(v, m)
	str := "/j #test-channel"
	want := "JOIN #test-channel\r\n"
	got := c.ParseInput(str)
	if !SpliceIsSame(got, []byte(want)) {
		t.Fatalf("expected: %s got: %s", want, string(got))
	}
}

func TestGetConnectionInfo(t *testing.T) {

	view := NewTerminalView()
	ipAndPort, nick := view.GetConnectionInfo()
	if ipAndPort == "" {
		t.Fatalf("error, ipAndPort empty")
	}

	if nick == "" {
		t.Fatalf("error, nick empty")
	}
}
