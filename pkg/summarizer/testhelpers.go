package summarizer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

const TESTS_YAML_FILE = "../tests.yaml"

type TestCase struct {
	Want           *FlightSummary    `yaml:"want"`
	flightStatuses *testCaseStatuses `yaml:"statuses"`
	fixture        []byte
}

type TestCaseFile struct {
	LiveFlights []*TestCase `yaml:"live_flights"`
}

func LoadTestCases(fixturePath string) (*TestCaseFile, error) {
	var cs TestCaseFile
	d, err := os.ReadFile(SummarizerTestFixturePath(fixturePath))
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(d, &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}

type flightScheduleStatus string

func (s *flightScheduleStatus) UnmarshalYAML(data []byte) error {
	pattern := "^(late|ontime|early)$"
	if err := matchRegexpDuringUnmarshal(pattern, data); err != nil {
		return err
	}
	*s = flightScheduleStatus(string(data))
	return nil
}

type flightStatus string

func (s *flightStatus) UnmarshalYAML(data []byte) error {
	pattern := "^(live)$"
	if err := matchRegexpDuringUnmarshal(pattern, data); err != nil {
		return err
	}
	*s = flightStatus(string(data))
	return nil
}

type testCaseStatuses struct {
	flightStatus    *flightStatus         `yaml:"flight"`
	DepartureStatus *flightScheduleStatus `yaml:"departure"`
	ArrivalStatus   *flightScheduleStatus `yaml:"arrival"`
}

const FIXTURE_PATH = "../fixtures"

func SummarizerTestFixturePath(name string) string {
	fp := filepath.Join(FIXTURE_PATH, name)
	re := regexp.MustCompile("\\.[a-zA-Z0-9]+$")
	if !re.Match([]byte(fp)) {
		fp = fp + ".html"
	}
	return fp
}

func LoadFixture(name string) ([]byte, error) {
	return os.ReadFile(SummarizerTestFixturePath(name))
}

func matchRegexpDuringUnmarshal(pattern string, data []byte) error {
	re := regexp.MustCompile(pattern)
	if !re.Match(data) {
		return fmt.Errorf("%s doesn't match wanted pattern '%s'", data, pattern)
	}
	return nil
}
