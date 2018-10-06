package figure_parser

import (
	"reflect"
	"testing"
)

func TestFormatNumber(t *testing.T) {
	testCases := []struct {
		in  interface{}
		out interface{}
	}{
		{" a ", " a "},
		{"100000", "100,000"},
		{" 100000 ", "100,000"},
		{" 100000.0101 ", "100,000.01"},
		{1, "1"},
		{1000, "1,000"},
		{1000.01, "1,000.01"},
		{1000.0101, "1,000.01"},
		{[]string{"a"}, []string{"a"}},
	}

	for i, tc := range testCases {
		out := FormatNumber(tc.in)
		if !reflect.DeepEqual(tc.out, out) {
			t.Error(i, ":", "want", tc.out, "got", out)
		}
	}
}
