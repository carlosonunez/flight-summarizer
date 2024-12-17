package testhelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/carlosonunez/flight-summarizer/types"
)

const FIXTURE_PATH = "./fixtures"

func loadFixture(name string) ([]byte, error) {
	path := filepath.Join(FIXTURE_PATH, strings.ReplaceAll(name, ".html", ""), ".html")
	return os.ReadFile(path)
}

func mustParseTime(tStr string) *time.Time {
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
