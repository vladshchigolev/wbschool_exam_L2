package clock

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClock_HostChecker(t *testing.T) {

	testCases := []struct {
		name    string
		isValid bool
		host    string
	}{
		{
			name:    "Default Host",
			isValid: true,
			host:    DefaultHost,
		},
		{
			name:    "Alternative Host",
			isValid: true,
			host:    AlternativeHost,
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			clock, _ := New(test.host)
			prec, loc := clock.CurrentTime()

			assert.NotNil(t, prec)
			assert.NotNil(t, loc)

		})
	}
}
