package sumamodels

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected CustomDate
		hasError bool
	}{
		{"\"Jan 2, 2006, 15:04:05 PM\"", CustomDate(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)), false},
		{"Feb 29, 2020, 12:00:00 PM", CustomDate(time.Date(2020, time.February, 29, 12, 0, 0, 0, time.UTC)), false},
		{"\"Invalid Date\"", CustomDate{}, true},
	}

	for _, tt := range tests {
		var cd CustomDate
		err := cd.UnmarshalJSON([]byte(tt.input))
		if tt.hasError {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.expected, cd)
		}
	}
}

func TestCustomDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		input    CustomDate
		expected string
	}{
		{CustomDate(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)), "\"2006-01-02T15:04:05Z\""},
		{CustomDate(time.Date(2020, time.February, 29, 12, 0, 0, 0, time.UTC)), "\"2020-02-29T12:00:00Z\""},
	}

	for _, tt := range tests {
		result, err := tt.input.MarshalJSON()
		require.NoError(t, err)
		assert.JSONEq(t, tt.expected, string(result))
	}
}
