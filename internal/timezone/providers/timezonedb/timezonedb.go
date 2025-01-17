package timezonedb

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/hashicorp/go-memdb"
)

const defaultTZDBDCCSVFilePath = "/data/timezonedb.csv"

var DB_TABLE_NAME = "timezone_data"

// TimeZoneDBDotComDBOptions configures how timezonedb.com databases are
// initialized. See DEFAULT_OPTIONS to view defaults.
type TimeZoneDBDotComDBOptions struct {
	// CSVFile changes the CSV file from which timezonedb.com databases will be
	// populated.
	CSVFile string
}

// TimeZoneDBDotComEntry is a time zone entry in a time zone database from
// timezonedb.com
type TimeZoneDBDotComEntry struct {
	// ID is a unique nonce used as an index.
	ID int64 `csv:"id"`

	// IsDaylightSavingsTime indicates whether DST applies for this entry.
	IsDaylightSavingsTime int `csv:"dst"`
	// ZoneName is the friendly name of the zone, like 'America/Houston'
	ZoneName string `csv:"zone_name"`
	// TimeStart is the time at which the Offset corresponding to this entry
	// applies. This database has multiple time samples to properly account for
	// daylight savings time.
	TimeStart int64 `csv:"time_start"`
	// UTCOffsetSeconds is the GMT/UTC offset, in seconds.
	UTCOffsetSeconds int64 `csv:"gmt_offset"`
	// CountryCode is a two-char ID of the country the zone belongs to, like 'US'
	CountryCode string `csv:"country_code"`
	// Abbreviation is the timezone abbreviation for the zone at a given time,
	// like 'CST'
	Abbreviation string `csv:"abbreviation"`
}

// TimeZoneDBDotComDB is a TimeZoneDatabase built from CSV dumps downloaded
// from timezonedb.com.
type TimeZoneDBDotComDB struct {
	contents []byte
	entries  []*TimeZoneDBDotComEntry
	db       *memdb.MemDB
}

// LookupUTCOffsetByID looks up a UTC offset given a timezone ID, like "CST".
func (db *TimeZoneDBDotComDB) LookupUTCOffsetByID(ID string, start time.Time) (int64, error) {
	if strings.ToLower(ID) == "utc" {
		return 0, nil
	}
	txn := db.db.Txn(false)
	defer txn.Abort()

	it, err := txn.GetReverse(DB_TABLE_NAME, "id")
	if err != nil {
		return 0, err
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		ent := obj.(*TimeZoneDBDotComEntry)
		if ent.Abbreviation != ID {
			continue
		}
		if start.Unix() <= ent.TimeStart {
			return ent.UTCOffsetSeconds, nil
		}
	}

	return 0, fmt.Errorf("start time '%d' for time zone ID '%s' is too early", start.Unix(), ID)
}

type csvReader interface {
	ReadFile(string) ([]byte, error)
}

type defaultCSVReader struct{}

func (f *defaultCSVReader) ReadFile(fp string) ([]byte, error) {
	return os.ReadFile(fp)
}

var DEFAULT_OPTIONS = &TimeZoneDBDotComDBOptions{
	CSVFile: defaultTZDBDCCSVFilePath,
}

// NewTimeZoneDBDotComDB creates a new in-memory timezone database.
func NewTimeZoneDBDotComDB(opts *TimeZoneDBDotComDBOptions) (*TimeZoneDBDotComDB, error) {
	var b []byte
	var err error
	if b, err = os.ReadFile(opts.CSVFile); err != nil {
		return nil, err
	}
	db := TimeZoneDBDotComDB{contents: b}
	if err := newInMemoryDB(&db); err != nil {
		return &db, err
	}
	if err := populateDB(&db); err != nil {
		return &db, err
	}
	return &db, nil
}

func newInMemoryDB(db *TimeZoneDBDotComDB) error {
	var err error
	db.db, err = memdb.NewMemDB(&memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			DB_TABLE_NAME: {
				Name: "timezone_data",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func populateDB(db *TimeZoneDBDotComDB) error {
	// Ensure that tz entries are cleared to save memory.
	defer func() { db.entries = nil }()
	if err := gocsv.Unmarshal(bytes.NewReader(db.contents), &db.entries); err != nil {
		return err
	}
	txn := db.db.Txn(true)
	for _, ent := range db.entries {
		if err := txn.Insert(DB_TABLE_NAME, ent); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}
