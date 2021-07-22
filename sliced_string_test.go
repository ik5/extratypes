package extratypes

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestSlicedStringJSONMarshal(t *testing.T) {
	type toCheck = struct {
		s        string
		expected SlicedString
		hasError bool
		err      error
	}
	checks := []toCheck{
		toCheck{
			s:        "",
			expected: nil,
			hasError: true,
			err:      errors.New("unexpected end of JSON input"),
		},
		toCheck{
			s:        "abcdef",
			expected: nil,
			hasError: true,
			err:      errors.New("invalid character 'a' looking for beginning of value"),
		},
		toCheck{
			s:        `"abcdef"`,
			expected: SlicedString{"abcdef"},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        `"a"`,
			expected: SlicedString{"a"},
			hasError: false,
			err:      nil,
		},

		toCheck{
			s:        "1",
			expected: nil,
			hasError: true,
			err:      errors.New("unsupported type 'float64'"),
		},
		toCheck{
			s:        "1.1",
			expected: nil,
			hasError: true,
			err:      errors.New("unsupported type 'float64'"),
		},
		toCheck{
			s:        "true",
			expected: nil,
			hasError: true,
			err:      errors.New("unsupported type 'bool'"),
		},
		toCheck{
			s:        `["a", "b", "c"]`,
			expected: SlicedString{"a", "b", "c"},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        `[ 1, "a" ]`,
			expected: nil,
			hasError: true,
			err:      errors.New("unsupported type 'float64' in slice"),
		},
		toCheck{
			s:        `[ "a", 1 ]`,
			expected: nil,
			hasError: true,
			err:      errors.New("unsupported type 'float64' in slice"),
		},
	}

	for _, check := range checks {
		var rec SlicedString
		err := json.Unmarshal([]byte(check.s), &rec)
		if check.hasError && err == nil {
			t.Errorf("Expected error '%s', but non exists", check.err)
			continue
		}

		if check.hasError && check.err == nil {
			t.Errorf("Have error '%s', but non expected?!", err)
			continue
		}

		if !check.hasError && err != nil {
			t.Errorf("Not expected err but '%s' exists", err)
			continue
		}

		if check.hasError && err != nil && check.err != nil && err.Error() != check.err.Error() {
			t.Errorf("Expected err (%T) '%s' but (%T) '%s' exists", check.err, check.err, err, err)
			continue
		}

		if !reflect.DeepEqual(check.expected, rec) {
			t.Errorf("Expected rec '%#v', but got '%#v'", check.expected, rec)
			continue
		}
	}
}

func TestSlicedStringScan(t *testing.T) {
	type toCheck = struct {
		s        interface{}
		expected SlicedString
		hasError bool
		err      error
	}

	validChecks := []toCheck{
		toCheck{
			s:        nil,
			expected: nil,
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        "",
			expected: SlicedString{""},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        "abcdef",
			expected: SlicedString{"abcdef"},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        `"abcdef"`,
			expected: SlicedString{`"abcdef"`},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        `"a"`,
			expected: SlicedString{`"a"`},
			hasError: false,
			err:      nil,
		},

		toCheck{
			s:        "1",
			expected: SlicedString{"1"},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        "1.1",
			expected: SlicedString{"1.1"},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        "true",
			expected: SlicedString{"true"},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        "false",
			expected: SlicedString{"false"},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        `["a", "b", "c"]`,
			expected: SlicedString{`["a", "b", "c"]`},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        `[1, "a" ]`,
			expected: SlicedString{`[1, "a" ]`},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        `["a", 1]`,
			expected: SlicedString{`["a", 1]`},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        []string{"a", "1"},
			expected: SlicedString{"a", "1"},
			hasError: false,
			err:      nil,
		},
		toCheck{
			s:        []interface{}{`["a", "1"]`},
			expected: SlicedString{`["a", "1"]`},
			hasError: false,
			err:      nil,
		},
	}

	t.Run("valid scan", func(t2 *testing.T) {
		for _, check := range validChecks {
			var rec SlicedString
			err := rec.Scan(check.s)
			if check.hasError && err == nil {
				t2.Errorf("Expected error '%s', but non exists", check.err)
				continue
			}

			if check.hasError && check.err == nil {
				t2.Errorf("Have error '%s', but non expected?!", err)
				continue
			}

			if !check.hasError && err != nil {
				t2.Errorf("Not expected err but '%s' exists", err)
				continue
			}

			if check.hasError && err != nil && check.err != nil && err.Error() != check.err.Error() {
				t2.Errorf("Expected err (%T) '%s' but (%T) '%s' exists", check.err, check.err, err, err)
				continue
			}

			if !reflect.DeepEqual(check.expected, rec) {
				t2.Errorf("Expected rec '%#v', but got '%#v'", check.expected, rec)
				continue
			}
		}
	})
}
