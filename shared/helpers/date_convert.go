package helpers

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return json.Marshal(t.Format("2006-01-02"))
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) ToTime() time.Time {
	return time.Time(d)
}

func (d DateOnly) Value() (driver.Value, error) {
	if d.ToTime().IsZero() {
		return nil, nil
	}
	return d.ToTime(), nil
}

func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		*d = DateOnly(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*d = DateOnly(v)
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		*d = DateOnly(t)
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		*d = DateOnly(t)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into DateOnly", value)
	}
}
