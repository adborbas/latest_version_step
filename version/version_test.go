package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	scenarios := []struct {
		raw      string
		expected *Version
	}{
		{
			raw:      "0.0.0",
			expected: &Version{0, 0, 0},
		},
		{
			raw:      "0.0.1",
			expected: &Version{0, 0, 1},
		},
		{
			raw:      "10.0.1",
			expected: &Version{10, 0, 1},
		},
		{
			raw:      "0.0.0.0",
			expected: nil,
		},
		{
			raw:      "something",
			expected: nil,
		},
	}

	for _, scenario := range scenarios {
		actualVersion := New(scenario.raw)
		require.Equal(t, scenario.expected, actualVersion)
	}
}

func Test_IsNewer(t *testing.T) {
	scenarios := []struct {
		left    Version
		right   Version
		isNewer bool
	}{
		{
			left:    Version{0, 0, 0},
			right:   Version{0, 0, 0},
			isNewer: false,
		},
		{
			left:    Version{0, 0, 1},
			right:   Version{0, 0, 0},
			isNewer: true,
		},
		{
			left:    Version{0, 0, 0},
			right:   Version{0, 0, 1},
			isNewer: false,
		},
		{
			left:    Version{0, 1, 0},
			right:   Version{0, 0, 1},
			isNewer: true,
		},
		{
			left:    Version{0, 1, 0},
			right:   Version{0, 100, 0},
			isNewer: false,
		},
		{
			left:    Version{1, 1, 0},
			right:   Version{1, 1, 0},
			isNewer: false,
		},
		{
			left:    Version{1, 1, 1},
			right:   Version{1, 1, 0},
			isNewer: true,
		},
	}

	for _, scenario := range scenarios {
		require.Equal(t, scenario.isNewer, scenario.left.IsNewer(scenario.right))
	}
}

func Test_ToString(t *testing.T) {
	scenarios := []struct {
		version  Version
		expected string
	}{
		{
			version:  Version{1, 0, 0},
			expected: "1.0.0",
		},
		{
			version:  Version{1, 1, 1},
			expected: "1.1.1",
		},
	}

	for _, scenario := range scenarios {
		actual := scenario.version.String()
		require.Equal(t, scenario.expected, actual)
	}
}
