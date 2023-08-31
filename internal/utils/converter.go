package utils

import (
	"fmt"
	"time"
)

const timestampPg = "2006-01-02 15:04:05.999999 -0700 MST"

func PgTimestampConverter(pgVal any) (time.Time, error) {
	timestampStr := fmt.Sprint(pgVal)
	timestamp, err := time.Parse(timestampPg, timestampStr)
	if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}
