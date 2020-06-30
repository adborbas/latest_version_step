package version

import( 
	"testing"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	scenarios := []struct{
		raw string
		expected *Version
	}{
		{
			raw: "0.0.0",
			expected: &Version{0,0,0},
		},
	}

	require.AssertTrue(true)
}
