package methods

import (
	"testing"
)

func TestParseOndemandSURLFromScript(t *testing.T) {
	tests := []struct {
		name string
		js   string
		want string
	}{
		{
			name: "find url",
			js:   `123:"ondemand.s",{123:"deadbeef"}`,
			want: "https://abs.twimg.com/responsive-web/client-web/ondemand.s.deadbeefa.js",
		},
		{
			name: "missing chunk",
			js:   `123:"main",{123:"deadbeef"}`,
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ParseOndemandSURLFromScript([]byte(test.js))
			if got != test.want {
				t.Fatalf("unexpected ondemand url: got %q want %q", got, test.want)
			}
		})
	}
}
