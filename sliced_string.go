package extratypes

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// SlicedString parses and translate either a string or a slice of strings
// based on serialization JSON/Text/DB.
type SlicedString []string

// UnmarshalJSON for contacts
func (s *SlicedString) UnmarshalJSON(data []byte) error {
	var str interface{}
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	items := reflect.ValueOf(str)
	kind := items.Kind()
	var result SlicedString

	switch kind {
	case reflect.String:
		result = append(result, items.String())

	case reflect.Slice:
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			switch item.Kind() {
			case reflect.String:
				result = append(result, item.String())
			case reflect.Interface:
				sliceItem := reflect.ValueOf(item.Interface())
				sliceKind := sliceItem.Kind()
				switch sliceKind {
				case reflect.String:
					result = append(result, sliceItem.String())
				default:
					return fmt.Errorf("unsupported type '%s' in slice", sliceKind)
				}
			}
		}
	default:
		return fmt.Errorf("unsupported type '%s'", kind)
	}

	*s = make(SlicedString, 0, len(result))
	*s = result
	return nil
}

// Scan implements the Scanner interface.
func (s *SlicedString) Scan(value interface{}) error {
	if value == nil {
		s = nil
		return nil
	}

	items := reflect.ValueOf(value)
	kind := items.Kind()
	var result SlicedString

	switch kind {
	case reflect.String:
		result = append(result, items.String())

	case reflect.Slice:
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			switch item.Kind() {
			case reflect.String:
				result = append(result, item.String())
			case reflect.Interface:
				sliceItem := reflect.ValueOf(item.Interface())
				sliceKind := sliceItem.Kind()
				switch sliceKind {
				case reflect.String:
					result = append(result, sliceItem.String())
				default:
					return fmt.Errorf("unsupported type '%s' in slice", sliceKind)
				}
			}
		}
	default:
		return fmt.Errorf("unsupported type '%s'", kind)
	}

	*s = make(SlicedString, 0, len(result))
	*s = result
	return nil

	return nil
}
