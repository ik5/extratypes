package extratypes

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

var (
	boolMap = map[string]bool{
		"true": true, "false": false,
		"yes": true, "no": false,
		"t": true, "f": false,
		"y": true, "n": false,
		"1": true, "0": false, "-1": false,
	}
)

// toType copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
// If src is nil, then the function return true, and dest remains as-is.
func toType(src, dest interface{}) (bool, error) {
	return false, nil
}

func asBool(src interface{}) bool {
	if src == nil {
		return false
	}

	switch src.(type) {
	case string:
		s := strings.ToLower(src.(string))
		status, ok := boolMap[s]
		if !ok {
			return false
		}
		return status
	case []byte:
		s := strings.ToLower(string(src.([]byte)))
		status, ok := boolMap[s]
		if !ok {
			return false
		}
		return status

	}

	val := reflect.ValueOf(src)
	v := val.Kind()
	switch v {
	case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		return i > 0
	case reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := val.Uint()
		return i > 0
	case reflect.Float32, reflect.Float64:
		f := math.Floor(val.Float())
		return int64(f) > 0
	}

	return false
}

func asInt(src interface{}, minRange, maxRange int64) interface{} {
	val := reflect.ValueOf(src)
	v := val.Kind()
	switch v {
	case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		if i > maxRange {
			return maxRange
		}
		if i < minRange {
			return minRange
		}
		return i
	case reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := val.Uint()
		if i > uint64(maxRange) {
			return maxRange
		}
		return int64(i)
	case reflect.Float32, reflect.Float64:
		f := math.Floor(val.Float())
		return asInt(int64(f), minRange, maxRange)
	case reflect.String:
		s := val.String()
		if s == "" {
			return int64(0)
		}

		if s[0] == '-' { // signed
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return int64(0)
			}

			return asInt(i, minRange, maxRange)
		}

		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return int64(0)
		}
		return asInt(i, minRange, maxRange)
	}

	return int64(0)

}

func asUint(src interface{}, minRange, maxRange uint64) interface{} {
	val := reflect.ValueOf(src)
	v := val.Kind()
	switch v {
	case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		if uint64(i) > maxRange {
			return maxRange
		}
		if uint64(i) < minRange {
			return minRange
		}
		return uint64(i)
	case reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := val.Uint()
		if i > maxRange {
			return maxRange
		}
		return i
	case reflect.Float32, reflect.Float64:
		f := math.Floor(val.Float())
		return asUint(uint64(f), minRange, maxRange)
	case reflect.String:
		s := val.String()
		if s == "" {
			return uint64(0)
		}

		if s[0] == '-' { // signed
			return minRange
		}

		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return uint64(0)
		}
		return asUint(i, minRange, maxRange)
	}

	return uint64(0)

}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v

	case []byte:
		return string(v)
	}

	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}

	return fmt.Sprintf("%v", src)
}

func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}

	return
}

func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}

	return err
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}

	c := make([]byte, len(b))
	copy(c, b)
	return c
}
