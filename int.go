package extratypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Int struct contains int data type that can be null, and also string
// on JSON and SQL, but value will be converted to int type
type Int struct {
	Val int
	Nil bool
}

func (i Int) String() string {
	if i.Nil {
		return "nil"
	}

	return fmt.Sprintf("%d", i.Val)
}

// Value interface for db
func (i Int) Value() (driver.Value, error) {
	return i.Val, nil
}

// Scan implement the Scan function from db interface
func (i *Int) Scan(v interface{}) error {
	isNil, err := toType(v, &i.Val)
	if err != nil {
		return err
	}

	i.Nil = isNil
	return nil
}

// MarshalJSON takes a Int and marshal it as a string
func (i Int) MarshalJSON() ([]byte, error) {
	if i.Nil {
		return json.Marshal(nil)
	}
	return json.Marshal(i.Val)
}

// UnmarshalJSON takes a slice of bytes and convert it to Int
func (i *Int) UnmarshalJSON(b []byte) error {
	var v interface{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	result, err := toType(v, &i.Val)
	if err != nil {
		return err
	}

	i.Nil = result
	return nil
}

// MarshalText takes a Int and marshal it as a string
func (i Int) MarshalText() ([]byte, error) {
	if i.Nil {
		return []byte(""), nil
	}

	return asByteSlice(i.String()), nil
}

// UnmarshalText takes a slice of bytes and convert it to Int
func (i *Int) UnmarshalText(b []byte) error {
	if b == nil {
		i.Nil = true
		return nil
	}

	if bytes.Compare(b, []byte("")) == 0 {
		i.Nil = true
		return nil
	}

	result, err := toType(b, &i.Val)
	if err != nil {
		return err
	}

	i.Nil = result
	return nil
}
