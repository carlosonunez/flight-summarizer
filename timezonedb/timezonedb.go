package timezonedb

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/hashicorp/go-memdb"
)

var DB_TABLE_NAME = "timezone_data"

// TimeZoneDBDotComEntry is a time zone entry in a time zone database from
// timezonedb.com
type TimeZoneDBDotComEntry struct {
	// ID is a unique nonce used as an index.
	ID int64

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
	csvFile string
	entries []*TimeZoneDBDotComEntry
	db      *memdb.MemDB
}

func (db *TimeZoneDBDotComDB) LookupUTCOffsetByID(ID string, start time.Time) (int64, error) {
	txn := db.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get(DB_TABLE_NAME, "abbreviation", ID)
	if err != nil {
		return 0, err
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		ent := obj.(*TimeZoneDBDotComEntry)
		if start.Unix() <= ent.TimeStart {
			return ent.UTCOffsetSeconds, nil
		}
	}

	return 0, fmt.Errorf("start time for time zone ID '%s' is too early", ID)
}

func NewTimeZoneDBDotComDB(csvFile string) (*TimeZoneDBDotComDB, error) {
	db := TimeZoneDBDotComDB{csvFile: csvFile}
	if err := newInMemoryDB(&db); err != nil {
		return &db, err
	}
	if err := populateDB(&db); err != nil {
		return &db, err
	}
	return &db, nil
}

func populateDB(db *TimeZoneDBDotComDB) error {
	// Ensure that tz entries are cleared to save memory.
	defer func() { db.entries = nil }()
	f, err := os.OpenFile(db.csvFile, os.O_RDONLY, os.ModeType)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := gocsv.UnmarshalFile(f, &db.entries); err != nil {
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
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"abbreviation": {
						Name:    "abbreviation",
						Indexer: &memdb.StringFieldIndex{Field: "Abbreviation"},
					},
					"time_start": {
						Name:    "time_start",
						Indexer: &memdb.StringFieldIndex{Field: "TimeStart"},
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
