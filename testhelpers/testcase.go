package testhelpers

import (
	"fmt"
	"os"

	"github.com/carlosonunez/flight-summarizer/types"
	"github.com/go-yaml/yaml"
)

const TESTS_YAML_FILE = "../tests.yaml"

type TestCase struct {
	Want           *types.FlightSummary `yaml:"want"`
	flightStatuses *testCaseStatuses    `yaml:"statuses"`
	fixture        []byte
}

func (c *TestCase) LoadFixtureForTestCase() ([]byte, error) {
	fp := fmt.Sprintf("%s_departure_%s_arrival",
		*c.flightStatuses.DepartureStatus,
		*c.flightStatuses.ArrivalStatus)
	return LoadFixture(fp)
}

type TestCaseFile struct {
	LiveFlights []*TestCase `yaml:"live_flights"`
}

func LoadTestCases() (*TestCaseFile, error) {
	var cs TestCaseFile
	d, err := os.ReadFile(TESTS_YAML_FILE)
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
