package types

import "time"
import "database/sql/driver"

const DateFormat string = "02Jan2006"

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

func (nt *NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}

	b := make([]byte, 0, len(DateFormat)+2)
	b = append(b, '"')
	b = nt.Time.AppendFormat(b, DateFormat)
	b = append(b, '"')
	return b, nil
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var err error
	if string(data) == "null" {
		nt.Valid = false
		return err
	}
	nt.Time, err = time.Parse(`"`+DateFormat+`"`, string(data))
	if err == nil {
		nt.Valid = true
	}
	return err
}
