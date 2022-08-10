package types

import (
	"fmt"
	"strconv"
	"time"
)

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("%d", time.Time(t).Unix())
	return []byte(stamp), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*t = Time(time.Unix(n, 0))
	return nil
}

func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}
