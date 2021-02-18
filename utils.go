package extratypes

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Values from math, do not use it just for that...
const (
	maxInt8                = 1<<7 - 1
	minInt8                = -1 << 7
	maxInt16               = 1<<15 - 1
	minInt16               = -1 << 15
	maxInt32               = 1<<31 - 1
	minInt32               = -1 << 31
	maxInt64               = 1<<63 - 1
	minInt64               = -1 << 63
	maxUint8               = 1<<8 - 1
	maxUint16              = 1<<16 - 1
	maxUint32              = 1<<32 - 1
	maxUint64              = 1<<64 - 1
	maxFloat32             = 3.40282346638528859811704183484516925440e+38   // 2**127 * (2**24 - 1) / 2**23
	smallestNonzeroFloat32 = 1.401298464324817070923729583289916131280e-45  // 1 / 2**(127 - 1 + 23)
	maxFloat64             = 1.797693134862315708145274237317043567981e+308 // 2**1023 * (2**53 - 1) / 2**52
	smallestNonzeroFloat64 = 4.940656458412465441765687928682213723651e-324 // 1 / 2**(1023 - 1 + 52)
)

var errNilPtr = errors.New("destination pointer is nil")

// toType copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
// If src is nil, then the function return true, and dest remains as-is.
func toType(src, dest interface{}) (bool, error) {
	return false, nil
}

func asInt8(src interface{}) int8 {
	val := reflect.ValueOf(src)
	v := val.Kind()
	switch v {
	case reflect.Int8:
		return src.(int8)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		if i > int64(maxInt8) {
			return maxInt8
		}
		if i < int64(minInt8) {
			return minInt8
		}
		return int8(i)
	case reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := val.Uint()
		if i > maxInt8 {
			return maxInt8
		}
		return int8(i)
	case reflect.String:
		s := val.String()
		if s == "" {
			return 0
		}

		if s[0] == '-' { // signed
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return 0
			}

			return asInt8(i)
		}

		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0
		}
		return asInt8(i)
	}

	return 0
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
