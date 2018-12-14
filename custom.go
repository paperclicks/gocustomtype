package custom

import (
	"regexp"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {

	//remove any extra " from the date string
	s := strings.Trim(string(b), "\"")

	regex := regexp.MustCompile("[0-9]")
	mask := regex.ReplaceAllString(s, "x")

	//try to parse from different formats
	switch mask {
	case "xxxx-xx-xx xx:xx:xx":
		t, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return err
		}
		ct.Time = t
	case "xxxx-xx-xx":
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return err
		}
		ct.Time = t

	}

	return nil
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {

	return []byte(ct.Time.Format("2006-01-02 15:04:05")), nil
}
