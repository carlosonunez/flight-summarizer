package testhelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/carlosonunez/flight-summarizer/types"
)

const FIXTURE_PATH = "../fixtures"

var TIMEZONE_DB_FIXTURE_PATH = FixturePath("timezonedb.csv")

func FixturePath(name string) string {
	fp := filepath.Join(FIXTURE_PATH, name)
	re := regexp.MustCompile("\\.[a-zA-Z0-9]+$")
	if !re.Match([]byte(fp)) {
		fp = fp + ".html"
	}
	return fp
}

func LoadFixture(name string) ([]byte, error) {
	return os.ReadFile(FixturePath(name))
}

func MustParseTime(tStr string) *time.Time {
	t, err := time.Parse(types.ISO8601TimeFormat, tStr)
	if err != nil {
		panic(fmt.Sprintf("time parsing during testing failed: %s", err))
	}
	return &t
}

func matchRegexpDuringUnmarshal(pattern string, data []byte) error {
	re := regexp.MustCompile(pattern)
	if !re.Match(data) {
		return fmt.Errorf("%s doesn't match wanted pattern '%s'", data, pattern)
	}
	return nil
}
