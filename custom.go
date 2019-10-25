package custom

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

//CustomJSONB is a custom tipe to handle JSON fields of DBs (for example PostgreSQL)
type CustomJSONB map[string]interface{}

func (jsonb *CustomJSONB) UnmarshalJSON(b []byte) error {

	var i interface{}
	var ok bool

	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	*jsonb, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed")
	}

	return nil
}

// func (c CustomJSONB) MarshalJSON() ([]byte, error) {

// 	return json.Marshal(c)
// }

//Value is one of the methods that the custom type CustomJSONB must implement
func (jsonb CustomJSONB) Value() (driver.Value, error) {

	j, err := json.Marshal(jsonb)
	return j, err
}

//Scan is a method that the custom type CustomJSONB must implement
func (jsonb *CustomJSONB) Scan(src interface{}) error {

	//check if src is of type []byte
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Scan - Type assertion .([]byte) failed")
	}

	//unmarshall the []byte value to an interface{}
	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	//Attempt a type assertion to pass the unmarshalled value to the pointer of the customJSONB
	*jsonb, ok = i.(map[string]interface{})

	//if type assertion fails, return the zero value
	if !ok {
		*jsonb = CustomJSONB{}
	}

	return nil

}
