package summarizers

import (
	"fmt"
	"testing"

	"github.com/carlosonunez/flight-summarizer/pkg/summarizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupSummarizer(t *testing.T) {
	want := &summarizer.ExampleFlightSummarizer{}
	got, err := Lookup("test")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%T", want), fmt.Sprintf("%T", got))
}
