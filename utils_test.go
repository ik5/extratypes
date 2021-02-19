package extratypes

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func TestCloneBytesNil(t *testing.T) {
	result := cloneBytes(nil)
	if result != nil {
		t.Errorf("coneByte with nil, must be nil, but %T was found", result)
	}
}

func TestCloneBytes(t *testing.T) {
	src := []byte{'H', 'e', 'l', 'l', 'o'}
	dest := cloneBytes(src)
	if !reflect.DeepEqual(src, dest) {
		t.Errorf("src: %+v not equal to dest: %+v", src, dest)
	}
}

func TestStrConvErrorSimple(t *testing.T) {
	err := errors.New("new error")
	destErr := strconvErr(err)

	if _, ok := destErr.(error); !ok {
		t.Errorf("destErr is not of type err, but of type: %T", destErr)
	}

	if destErr.Error() != err.Error() {
		t.Errorf("err: %s not equal to destErr: %s", err, destErr)
	}
}

func TestStrConvErrorNumError(t *testing.T) {
	err := &strconv.NumError{
		Func: "foo",
		Num:  "0",
		Err:  errors.New("new error"),
	}
	destErr := strconvErr(err)

	if _, ok := destErr.(error); !ok {
		t.Errorf("destErr is not of type err, but of type: %T", destErr)
	}

	if destErr.Error() != err.Err.Error() {
		t.Errorf("err: [%s] not equal to destErr: [%s]", err, destErr)
	}
}

func TestAsBytesInt(t *testing.T) {
	val := reflect.ValueOf(10)
	buf := []byte{10}
	result, ok := asBytes(buf, val)

	if !ok {
		t.Error("Invalid convert")
	}

	if result[0] != buf[0] {
		t.Errorf("result[0] [%v] is not equal to buf[0] [%v]", result[0], buf[0])
	}
}

func TestAsBytesUInt(t *testing.T) {
	val := reflect.ValueOf(uint(10))
	buf := []byte{10}
	result, ok := asBytes(buf, val)

	if !ok {
		t.Error("Invalid convert")
	}

	if result[0] != buf[0] {
		t.Errorf("result[0] [%v] is not equal to buf[0] [%v]", result[0], buf[0])
	}
}

func TestAsBytesFloat32(t *testing.T) {
	val := reflect.ValueOf(float32(10))
	buf := []byte{10}
	result, ok := asBytes(buf, val)

	if !ok {
		t.Error("Invalid convert")
	}

	if result[0] != buf[0] {
		t.Errorf("result[0] [%v] is not equal to buf[0] [%v]", result[0], buf[0])
	}
}

func TestAsBytesFloat64(t *testing.T) {
	val := reflect.ValueOf(float64(10))
	buf := []byte{10}
	result, ok := asBytes(buf, val)

	if !ok {
		t.Error("Invalid convert")
	}

	if result[0] != buf[0] {
		t.Errorf("result[0] [%v] is not equal to buf[0] [%v]", result[0], buf[0])
	}
}

func TestAsBytesBool(t *testing.T) {
	val := reflect.ValueOf(true)
	buf := []byte{1}
	result, ok := asBytes(buf, val)

	if !ok {
		t.Error("Invalid convert")
	}

	if result[0] != buf[0] {
		t.Errorf("result[0] [%v] is not equal to buf[0] [%v]", result[0], buf[0])
	}
}

func TestAsBytesString(t *testing.T) {
	val := reflect.ValueOf("hello")
	buf := []byte{'h', 'e', 'l', 'l', 'o'}
	result, ok := asBytes(buf, val)

	if !ok {
		t.Error("Invalid convert")
	}

	if !reflect.DeepEqual(result, append(buf, buf...)) {
		t.Errorf("result [%v] is not equal to buf [%v]", result, buf)
	}
}

func TestAsBytesPtr(t *testing.T) {
	val := reflect.ValueOf(nil)
	buf := []byte{'h'}
	result, ok := asBytes(buf, val)

	if ok {
		t.Error("Invalid convert")
	}

	if reflect.DeepEqual(result, buf) {
		t.Errorf("result[0] [%v] is equal to buf[0] [%v]", result, buf)
	}
}

func TestAsStringString(t *testing.T) {
	src := "hello"
	dest := asString(src)

	if dest != src {
		t.Errorf("dest '%s' not equal to src '%s'", dest, src)
	}
}

func TestAsStringByteSlice(t *testing.T) {
	src := []byte{'h', 'e', 'l', 'l', 'o'}
	dest := asString(src)

	if dest != string(src) {
		t.Errorf("dest '%s' not equal to src '%s'", dest, src)
	}

}

func TestAsStringInt(t *testing.T) {
	src := 10
	dest := asString(src)
	if dest != "10" {
		t.Errorf("dest '%s' is not '%d'", dest, src)
	}
}

func TestAsStringUInt(t *testing.T) {
	src := uint(10)
	dest := asString(src)
	if dest != "10" {
		t.Errorf("dest '%s' is not '%d'", dest, src)
	}
}

func TestAsStringFloat32(t *testing.T) {
	src := float32(10)
	dest := asString(src)
	if dest != "10" {
		t.Errorf("dest '%s' is not '%f'", dest, src)
	}
}

func TestAsStringFloat64(t *testing.T) {
	src := float64(10)
	dest := asString(src)
	if dest != "10" {
		t.Errorf("dest '%s' is not '%f'", dest, src)
	}
}

func TestAsStringBool(t *testing.T) {
	src := true
	dest := asString(src)

	if dest != "true" {
		t.Errorf("Expected dest to be '%t', got: '%s'", src, dest)
	}
}

func TestAsStringPtr(t *testing.T) {
	var src interface{} = nil
	dest := asString(src)
	if dest != "<nil>" {
		t.Errorf("dest expected to be '<nil>' but got '%s'", dest)
	}
}

func TestAsIntInt8(t *testing.T) {
	src := int8(10)
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != int64(src) {
		t.Errorf("dest [%d] != src[%d]", dest, src)
	}
}

func TestAsIntInt16Big(t *testing.T) {
	src := int64(256)
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) == src {
		t.Errorf("dest [%d] == src[%d]", dest, src)
	}

	if dest.(int64) != int64(math.MaxInt8) {
		t.Errorf("dest [%d] != %d", dest, math.MaxInt8)
	}
}

func TestAsIntInt16Small(t *testing.T) {
	src := int64(-256)
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) == src {
		t.Errorf("dest [%d] == src[%d]", dest, src)
	}

	if dest != int64(math.MinInt8) {
		t.Errorf("dest [%d] != %d", dest, math.MinInt8)
	}
}

func TestAsIntInt16Equal(t *testing.T) {
	src := int64(10)
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != src {
		t.Errorf("dest [%d] != src[%d]", dest, src)
	}
}

func TestAsIntUint16Big(t *testing.T) {
	src := uint64(256)
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) == int64(src) {
		t.Errorf("dest [%d] == src[%d]", dest, src)
	}

	if dest != int64(math.MaxInt8) {
		t.Errorf("dest [%d] != %d", dest, math.MaxInt8)
	}
}

func TestAsIntUint16Equal(t *testing.T) {
	src := uint64(10)
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != int64(src) {
		t.Errorf("dest [%d] != src[%d]", dest, src)
	}
}

func TestAsIntStringEmpty(t *testing.T) {
	src := ""
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != 0 {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsIntStringSignedNormal(t *testing.T) {
	src := "-10"
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != int64(-10) {
		t.Errorf("dest [%d] is not -10", dest)
	}
}

func TestAsIntStringSignedInvalid(t *testing.T) {
	src := "-a10"
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != int64(0) {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsIntStringSignedBig(t *testing.T) {
	src := "-1000"
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != int64(math.MinInt8) {
		t.Errorf("dest [%d] is not %d", dest, math.MinInt8)
	}
}

func TestAsIntStringUnsignedNormal(t *testing.T) {
	src := "10"
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != int64(10) {
		t.Errorf("dest [%d] is not 10", dest)
	}
}

func TestAsIntStringUnsignedBig(t *testing.T) {
	src := "1000"
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != int64(math.MaxInt8) {
		t.Errorf("dest [%d] is not %d", dest, math.MaxInt8)
	}
}

func TestAsIntStringUnsignedInvalid(t *testing.T) {
	src := "a10"
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != int64(0) {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsIntUnknown(t *testing.T) {
	var src struct{}
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != 0 {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsIntFloatNormal(t *testing.T) {
	src := 22.511
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != 22 {
		t.Errorf("dest [%d] != 22", dest)
	}
}

func TestAsIntFloatBig(t *testing.T) {
	src := 1000.5
	dest := asInt(src, math.MinInt8, math.MaxInt8)

	if dest.(int64) != int64(math.MaxInt8) {
		t.Errorf("dest [%d] != %d", dest, math.MaxInt8)
	}
}

//////////////////////////////////////////////////////////////////////////////
func TestAsUintInt8(t *testing.T) {
	src := int8(10)
	dest := asUint(src, 0, math.MaxUint8)

	if dest.(uint64) != uint64(src) {
		t.Errorf("dest [%d] != src[%d]", dest, src)
	}
}

func TestAsUIntInt16Big(t *testing.T) {
	src := int64(256)
	dest := asUint(src, 0, math.MaxUint8)

	if dest.(uint64) == uint64(src) {
		t.Errorf("dest [%d] == src[%d]", dest, src)
	}

	if dest.(uint64) != uint64(math.MaxUint8) {
		t.Errorf("dest [%d] != %d", dest, math.MaxInt8)
	}
}

func TestAsUintUint16Equal(t *testing.T) {
	src := int64(10)
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != uint64(src) {
		t.Errorf("dest [%d] != src[%d]", dest, src)
	}
}

func TestAsUintUint16Big(t *testing.T) {
	src := uint64(256)
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) == uint64(src) {
		t.Errorf("dest [%d] == src[%d]", dest, src)
	}

	if dest.(uint64) != uint64(math.MaxInt8) {
		t.Errorf("dest [%d] != %d", dest, math.MaxInt8)
	}
}

func TestAsUintStringEmpty(t *testing.T) {
	src := ""
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != 0 {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsUintStringSignedNormal(t *testing.T) {
	src := "-10"
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != uint64(0) {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsUintStringSignedInvalid(t *testing.T) {
	src := "-a10"
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != uint64(0) {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsUintStringSignedBig(t *testing.T) {
	src := "-1000"
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != uint64(0) {
		t.Errorf("dest [%d] is not %d", dest, math.MinInt8)
	}
}

func TestAsUintStringUnsignedNormal(t *testing.T) {
	src := "10"
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != uint64(10) {
		t.Errorf("dest [%d] is not 10", dest)
	}
}

func TestAsUintStringUnsignedBig(t *testing.T) {
	src := "1000"
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != uint64(math.MaxInt8) {
		t.Errorf("dest [%d] is not %d", dest, math.MaxInt8)
	}
}

func TestAsUintStringUnsignedInvalid(t *testing.T) {
	src := "a10"
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != uint64(0) {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsUintUnknown(t *testing.T) {
	var src struct{}
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != 0 {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsUintFloatNormal(t *testing.T) {
	src := 22.511
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != 22 {
		t.Errorf("dest [%d] != 22", dest)
	}
}

func TestAsUintFloatBig(t *testing.T) {
	src := 1000.5
	dest := asUint(src, 0, math.MaxInt8)

	if dest.(uint64) != uint64(math.MaxInt8) {
		t.Errorf("dest [%d] != %d", dest, math.MaxInt8)
	}
}

func TestAsUintRange(t *testing.T) {
	min := uint64(5)
	src := 1
	dest := asUint(src, min, math.MaxInt8)

	if dest.(uint64) != min {
		t.Errorf("dest [%d] != %d", dest, min)
	}
}
