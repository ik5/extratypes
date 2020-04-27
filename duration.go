package extratypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

// Duration is wrapper for time.Duration with additional methods
type Duration struct {
	time.Duration
	Nil bool
}

// Value that the database usage will see
func (d Duration) Value() (driver.Value, error) {
	return int64(d.Duration), nil
}

// Scan the result from a query and assign it to the struct
func (d *Duration) Scan(v interface{}) error {
	if v == nil {
		d.Duration = 0
		d.Nil = true
		return nil
	}

	val := reflect.ValueOf(v)
	kind := val.Kind()

	switch kind {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.Duration = time.Duration(val.Int())
	case reflect.String:
		var err error
		d.Duration, err = time.ParseDuration(val.String())
		if err != nil {
			return err
		}
	case reflect.Float32, reflect.Float64:
		d.Duration = time.Duration(val.Float())
	default:
		return fmt.Errorf("Invalid type of %T", val)
	}

	return nil
}

// MarshalJSON takes a duration and marshal it as a string
func (d Duration) MarshalJSON() ([]byte, error) {
	if d.Nil {
		return json.Marshal(nil)
	}
	return json.Marshal(d.String())
}

// MarshalText takes duration and marshal it as a string
func (d *Duration) MarshalText() ([]byte, error) {
	if d.Nil {
		return []byte(""), nil
	}
	return []byte(d.Duration.String()), nil
}

// UnmarshalJSON takes a slice of bytes and convert it to Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	val := reflect.ValueOf(v)
	kind := val.Kind()

	switch kind {
	case reflect.Float32, reflect.Float64:
		d.Duration = time.Duration(val.Float())
		return nil
	case reflect.String:
		var err error
		d.Duration, err = time.ParseDuration(val.String())
		if err != nil {
			return err
		}
		return nil
	case reflect.Map:
		iface := val.Interface()
		m, success := iface.(map[string]interface{})
		if !success {
			return errors.New("Invalid map type")
		}

		l := len(m)
		if l == 0 {
			return errors.New("no content found")
		}

		if l > 1 {
			return fmt.Errorf("Length %d is too big", l)
		}

		var err error
		for _, value := range m {
			val2 := reflect.ValueOf(value)
			kind2 := val2.Kind()
			switch kind2 {
			case reflect.Float32, reflect.Float64:
				d.Duration = time.Duration(val2.Float())
				return nil
			case reflect.String:
				d.Duration, err = time.ParseDuration(val2.String())
				if err != nil {
					return err
				}
				return nil
			case reflect.Invalid:
				d.Duration = 0
				d.Nil = true
				return nil

			default:
				return fmt.Errorf("Invalid kind of value: %s", kind2)
			}
		}
	case reflect.Invalid:
		d.Duration = 0
		d.Nil = true
		return nil
	default:
		return fmt.Errorf("Invalid duration type: %T, %+v", v, v)
	}
	return errors.New("Unknown error")
}

// UnmarshalText unmar
func (d *Duration) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		d.Duration = 0
		d.Nil = true
		return nil
	}

	var err error
	d.Duration, err = time.ParseDuration(string(b))

	if err != nil {
		d.Duration = 0
		d.Nil = true
		return err
	}
	return nil
}

func (d Duration) String() string {
	if d.Nil {
		return "nil"
	}
	return d.Duration.String()
}
