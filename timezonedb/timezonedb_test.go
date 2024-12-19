package timezonedb

import (
	"testing"
	"time"

	"github.com/carlosonunez/flight-summarizer/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeZoneDBDotComFromCSV(t *testing.T) {
	db, err := NewTimeZoneDBDotComDB(testhelpers.TIMEZONE_DB_FIXTURE_PATH)
	require.NoError(t, err)
	exampleTime := time.Unix(1734637838, 0)
	wantOffset := int64(-21600) // -0600
	gotOffset, err := db.LookupUTCOffsetByID("CST", exampleTime)
	require.NoError(t, err)
	assert.Equal(t, wantOffset, gotOffset)
}
