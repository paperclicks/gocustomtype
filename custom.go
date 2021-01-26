package customtype

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
	case "xxxx-xx-xx xx:xx:xx +xxxx", "xxxx-xx-xx xx:xx:xx -xxxx":
		t, err := time.Parse("2006-01-02 15:04:05 -0700", s)
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

type SqlNullBool struct {
	sql.NullBool
}

// MarshalJSON for NullBool
func (nb *SqlNullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON for NullBool
func (nb *SqlNullBool) UnmarshalJSON(b []byte) error {
	var intValue int
	var boolValue bool

	if string(b) == "null" {
		// The key was set to null
		nb.Valid = false
		return nil
	}

	//attempt to unmarshal to bool, then to int
	if err := json.Unmarshal(b, &boolValue); err != nil {

		if err = json.Unmarshal(b, &intValue); err != nil {
			nb.Valid = false
			return err
		} else {

			nb.Valid = true
			switch intValue {
			case 0:
				nb.Bool = false
			case 1:
				nb.Bool = true
			default:
				nb.Valid = false
			}
		}

	} else {

		nb.Valid = true
		nb.Bool = boolValue
	}

	return nil
}

type SqlNullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *SqlNullString) MarshalJSON() ([]byte, error) {

	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *SqlNullString) UnmarshalJSON(b []byte) error {
	var val string

	if string(b) == "null" {
		// The key was set to null so simply invalid the object
		ns.Valid = false
		return nil
	}
	if string(b) == "" {
		ns.Valid = true
		ns.String = ""
		return nil
	}

	err := json.Unmarshal(b, &val)
	if err != nil {
		ns.Valid = false
		return err
	}
	ns.Valid = true
	ns.String = val

	return err
}

type SqlNullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullFloat64
func (nf *SqlNullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON for NullFloat64
func (nf *SqlNullFloat64) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		// The key was set to null
		nf.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = (err == nil)
	return err
}

type SqlNullInt32 struct {
	sql.NullInt32
}

// MarshalJSON for NullInt64
func (ni *SqlNullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int32)
}

// UnmarshalJSON for NullInt64
func (ni *SqlNullInt32) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		// The key was set to null
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.Int32)
	ni.Valid = (err == nil)
	return err
}

type SqlNullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *SqlNullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullInt64
func (ni *SqlNullInt64) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		// The key was set to null
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

type SqlNullTime struct {
	sql.NullTime
}

func (nt *SqlNullTime) UnmarshalJSON(b []byte) error {

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
		nt.Time = t
	case "xxxx-xx-xx":
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return err
		}
		nt.Time = t
	case "xxxx-xx-xx xx:xx:xx +xxxx", "xxxx-xx-xx xx:xx:xx -xxxx":
		t, err := time.Parse("2006-01-02 15:04:05 -0700", s)
		if err != nil {
			return err
		}
		nt.Time = t
	}

	return nil
}

func (nt *SqlNullTime) MarshalJSON() ([]byte, error) {

	if !nt.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
		return json.Marshal(t)
	}

	return json.Marshal(nt.Time)

}

//JSONB is a custom type to handle JSON columns ind DBs which are defined as string
type JSONB map[string]interface{}

func (jsonb *JSONB) UnmarshalJSON(b []byte) error {

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
func (jsonb JSONB) Value() (driver.Value, error) {

	j, err := json.Marshal(jsonb)
	return j, err
}

//Scan is a method that the custom type CustomJSONB must implement
func (jsonb *JSONB) Scan(src interface{}) error {

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

	//Attempt a type assertion to pass the unmarshalled value to the pointer of the customJSONB. If it fails return the zero value of CustomJSONB
	*jsonb, ok = i.(map[string]interface{})
	if !ok {
		*jsonb = JSONB{}
	}

	return nil

}

//JSON is a custom type to handle database fields that are suposed to contain json content.

type JSON struct {
	string
}

//Value returns the underlying value encapsulated by the custom type.
//It is used by sql to extract the value to be inserted in DB
func (cjson JSON) Value() (driver.Value, error) {

	m := json.RawMessage(cjson.string)

	j, err := m.MarshalJSON()
	if err != nil {
		return j,err
	}
	return j, err
}

//Scan sets the value of the underlying type using the va
func (cjson *JSON) Scan(src interface{}) error {

	switch src.(type) {
	case nil:
		cjson.string = ""
		return nil
	case string:
		s:=src.(string)
		r:=json.RawMessage(s)
		b,err:=r.MarshalJSON()
		if err != nil {
			fmt.Errorf("invalid json %s", s)
		}
		cjson.string = string(b)
		return nil
	}



	return fmt.Errorf("src is of type %s, not string", reflect.TypeOf(src))

}


func (cjson *JSON) MarshalJSON() ([]byte, error) {
	s:=cjson.string
	//if the conent is an empty string then represent it as an empty object,{}
	if cjson.string==""{
		s="{}"
	}
	r := json.RawMessage(s)

	b, err := r.MarshalJSON()
	if err != nil {
		return r, err

	}
	return b, err

}

func (cjson *JSON) UnmarshalJSON(b []byte) error {

	s:=string(b)
	if s=="null"{
		s=""
	}

	//check if valid json
	r:=json.RawMessage(s)
	b,err:=r.MarshalJSON()
	if err != nil {
		return fmt.Errorf("invalid json string %s. %s",s,err.Error())
	}
	cjson.string = s

	return nil
}
