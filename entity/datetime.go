package entity

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type DateTime struct {
	time.Time
}

func (t DateTime) FirstDayOfMonth() DateTime {
	year, month, _ := t.Date()
	currentLocation := t.Location()

	return DateTime{time.Date(year, month, 1, 0, 0, 0, 0, currentLocation)}
}

func (t DateTime) LastDayOfMonth() DateTime {
	return DateTime{Time: t.FirstDayOfMonth().AddDate(0, 1, 0).Time.Add(-time.Nanosecond)}
}

func (t DateTime) Between(before DateTime, after DateTime) bool {
	return (t.After(before) || t.Equal(before)) && (t.Before(after))
}

func (t DateTime) After(u DateTime) bool {
	return t.Time.After(u.Time)
}

func (t DateTime) Equal(u DateTime) bool {
	return t.Time.Equal(u.Time)
}

func (t DateTime) Before(u DateTime) bool {
	return t.Time.Before(u.Time)
}

func (t DateTime) Add(d time.Duration) DateTime {
	return DateTime{t.Time.Add(d)}
}

func (t DateTime) Sub(u DateTime) time.Duration {
	return t.Time.Sub(u.Time)
}

func (t DateTime) AddDate(years int, months int, days int) DateTime {
	return DateTime{t.Time.AddDate(years, months, days)}
}

func (t DateTime) Min(u DateTime) DateTime {
	if t.Before(u) {
		return t
	} else {
		return u
	}
}

func (t DateTime) Max(u DateTime) DateTime {
	if t.Before(u) {
		return u
	} else {
		return t
	}
}

// Scan https://gorm.io/docs/data_types.html#Implements-Customized-Data-Type
func (t *DateTime) Scan(b interface{}) (err error) {
	switch x := b.(type) {
	case time.Time:
		t.Time = x
	default:
		err = fmt.Errorf("unsupported scan type %T", b)
	}
	return
}

// Value https://gorm.io/docs/data_types.html#Implements-Customized-Data-Type
func (t DateTime) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time.Format("2006-01-02 15:04:05"), nil
}
