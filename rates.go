package extratypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Rate holds information about type of sending rates
type Rate struct {
	Count uint64
	Sleep Duration
}

// NewRate parse raw string and translate it to sending rate
// if there was an issue with parsing, then an error returns
// if the raw input is empty, then an empty Rate is returned
func NewRate(raw string) (Rate, error) {
	var result Rate
	err := result.SetFromRaw(raw)
	if err != nil {
		return Rate{}, err
	}
	return result, nil
}

// SetFromRaw takes data from raw string and set it to Rate.
// if there was an issue with parsing, then an error returns
// if the raw input is empty, then an empty Rate is returned
func (r *Rate) SetFromRaw(raw string) error {
	if len(raw) == 0 {
		return nil
	}

	splitted := strings.Split(string(raw), "/")
	if len(splitted) != 2 {
		return fmt.Errorf("Invalid Rate structure for %s", raw)
	}

	var err error

	r.Count, err = strconv.ParseUint(splitted[0], 10, 64)
	if err != nil {
		return err
	}

	r.Sleep.Duration, err = time.ParseDuration(splitted[1])
	if err != nil {
		return err
	}

	return nil
}

// IsEmpty rate returns true
func (r Rate) IsEmpty() bool {
	return r.Count == 0 &&
		r.Sleep.Duration == 0
}

func (r Rate) String() string {
	return fmt.Sprintf("%d/%s", r.Count, r.Sleep)
}

// MarshalJSON returns a json value represented by Rate
func (r Rate) MarshalJSON() ([]byte, error) {
	if r.IsEmpty() {
		return json.Marshal("")
	}

	return json.Marshal(r.String())
}

// UnmarshalJSON to Rate
func (r *Rate) UnmarshalJSON(b []byte) error {
	var v interface{}
	var err error

	err = json.Unmarshal(b, &v)

	if err != nil {
		return err
	}

	val := reflect.ValueOf(v)
	kind := val.Kind()

	switch kind {
	case reflect.String:
		s := val.String()
		if s == "" {
			r.Count = 0
			r.Sleep.Duration = 0
			return nil
		}
		if s[0] == '"' {
			s = strings.ReplaceAll(s, "\"", "")
		}

		err = r.SetFromRaw(s)
		if err != nil {
			r.Count = 0
			r.Sleep.Duration = 0
		}

		return err

	case reflect.Map:
		iface := val.Interface()

		m, success := iface.(map[string]interface{})
		if !success {
			return fmt.Errorf("Invalid map type: %#v", iface)
		}

		l := len(m)
		if l == 0 {
			return errors.New("no content found")
		}

		if l > 1 {
			return fmt.Errorf("Length %d is too big", l)
		}

		for _, value := range m {
			val2 := reflect.ValueOf(value)
			kind2 := val2.Kind()
			switch kind2 {
			case reflect.String:
				s := val2.String()
				if s[0] == '"' {
					s = strings.ReplaceAll(s, "\"", "")
				}

				err = r.SetFromRaw(s)
				if err != nil {
					r.Count = 0
					r.Sleep.Duration = 0
				}

				return err

			default:
				return fmt.Errorf("Invalid kind of value: %s", kind2)
			}
		}
		return errors.New("No key/value for map")
	default:
		return fmt.Errorf("Invalid data type for %s, expected string/map, but got %s",
			b, kind)
	}

}
