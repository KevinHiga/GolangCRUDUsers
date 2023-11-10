package utils

import (
	"time"
	dps "github.com/markusmobius/go-dateparser"
)

func ParseDate(date string) (time.Time, error) {
	dt, err := dps.Parse(nil, date)
	return dt.Time, err
}