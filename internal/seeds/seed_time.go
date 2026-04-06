package seeds

import (
	"database/sql/driver"
	"time"
)

// SeedTime parses RFC3339 timestamps from CSV (same idea as users API seeders).
type SeedTime struct {
	time.Time
	isNil bool
}

func (st *SeedTime) UnmarshalCSV(data []byte) (err error) {
	if len(data) == 0 {
		st.isNil = true
		return nil
	}
	st.Time, err = time.Parse(time.RFC3339, string(data))
	if err != nil {
		st.Time, err = time.Parse("2006-01-02T15:04:05.999999Z", string(data))
	}
	return err
}

func (st SeedTime) Value() (driver.Value, error) {
	if st.isNil {
		return nil, nil
	}
	return st.Time, nil
}
