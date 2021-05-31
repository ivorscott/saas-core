package tests

import (
	"bytes"
	"github.com/devpies/devpie-client-core/users/api/handlers"
	"strings"
	"testing"
)

func TestParseOrigins(t *testing.T) {
	testcases := []struct {
		arg  string
		want []string
	}{
		{"http://localhost:3000, https://devpie.local:3000", []string{"http://localhost:3000", "https://devpie.local:3000"}},
		{"   http://localhost:3000", []string{"http://localhost:3000"}},
		{"   http://localhost:3000 , https://devpie.local:3000   ", []string{"http://localhost:3000", "https://devpie.local:3000"}},
		{"", []string{}},
	}

	for _, v := range testcases {
		got := handlers.ParseOrigins(v.arg)
		if !bytes.Equal([]byte(strings.Join(got, ",")), []byte(strings.Join(v.want, ","))) {
			t.Errorf("ParseOrigins(%v) = %v; want %v", v.arg, got, v.want)
		}
	}
}
