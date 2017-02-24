package accounting

import (
	"testing"
)

func TestOrderValuesFromPrice(t *testing.T) {
	// Tests the expected output for prices.
	var tests = []struct {
		arg  string
		want string
	}{
		{"123.39", "$53.67"},
	}

	for _, test := range tests {

		if got := OrderValuesFromPrice(test.arg).InitialPaymentFmt; got != test.want {
			t.Errorf("OrderValuesFromPrice(%q) = %v", test.arg, got)
		}
	}
}
