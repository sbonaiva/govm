package util_test

import (
	"testing"

	"github.com/sbonaiva/govm/internal/test"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		name      string
		function  func(msg string, args ...any)
		input     string
		inputArgs []any
		expected  string
	}{
		{
			name:     "Print Success No Args",
			function: util.PrintSuccess,
			input:    "Success",
			expected: "Success\n",
		},
		{
			name:      "Print Success With Args",
			function:  util.PrintSuccess,
			input:     "Success %s",
			inputArgs: []any{"With Args"},
			expected:  "Success With Args\n",
		},
		{
			name:     "Print Warning No Args",
			function: util.PrintWarning,
			input:    "Warning",
			expected: "Warning\n",
		},
		{
			name:      "Print Warning With Args",
			function:  util.PrintWarning,
			input:     "Warning %s",
			inputArgs: []any{"With Args"},
			expected:  "Warning With Args\n",
		},
		{
			name:     "Print Error No Args",
			function: util.PrintError,
			input:    "Error",
			expected: "Error: Error\n",
		},
		{
			name:      "Print Error With Args",
			function:  util.PrintError,
			input:     "Error %s",
			inputArgs: []any{"With Args"},
			expected:  "Error: Error With Args\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := test.CaptureOutput(func() error {
				tt.function(tt.input, tt.inputArgs...)
				return nil
			})
			assert.Equal(t, tt.expected, actual)
		})
	}
}
