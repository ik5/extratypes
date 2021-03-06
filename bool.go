package extratypes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Bool contain a boolean data that can be null, and also string
// on JSON and SQL, but value will be converted into bool type
type Bool struct {
	sql.NullBool
	Val bool `json:"val" toml:"val"`
	Nil bool `json:"nil" toml:"nil"`
}

// Scan implements the Scanner interface.
func (b *Bool) Scan(value interface{}) error {
	if value == nil {
		b.Val = false
		b.Nil = true
		return nil
	}

	b.Nil = false
	b.Val = asBool(value)
	return nil
}

func (b Bool) String() string {
	if b.Nil {
		return "nil"
	}
	return strconv.FormatBool(b.Val)
}

// Value implements the driver Valuer interface.
func (b Bool) Value() (driver.Value, error) {
	if b.Nil {
		return nil, nil
	}

	return b.Val, nil
}

// MarshalJSON implement the Marshaler interface
func (b Bool) MarshalJSON() ([]byte, error) {
	if b.Nil {
		return json.Marshal(nil)
	}

	return json.Marshal(b.Val)
}

// UnmarshalJSON implement the un-Marshaler interface
func (b *Bool) UnmarshalJSON(buf []byte) error {
	var v interface{}
	err := json.Unmarshal(buf, &v)
	if err != nil {
		return err
	}

	if v == nil {
		b.Nil = true
		b.Val = false
		return nil
	}

	b.Val = asBool(v)
	return nil
}

// MarshalText implement Text Marshaller interface
func (b Bool) MarshalText() ([]byte, error) {
	if b.Nil {
		return []byte(""), nil
	}

	return asByteSlice(b.String()), nil
}

// UnmarshalText implement the text un-Marshaller interface
func (b *Bool) UnmarshalText(buf []byte) error {
	if buf == nil || bytes.Compare(buf, []byte("")) == 0 ||
		bytes.Compare(buf, []byte("null")) == 0 ||
		bytes.Compare(buf, []byte("nil")) == 0 {

		b.Nil = true
		return nil
	}

	result, err := toType(buf, &b.Val)
	if err != nil {
		return err
	}
	fmt.Printf("buf: %s | result: %t | b.Val: %t\n", buf, result, b.Val)

	b.Nil = result
	return nil
}
