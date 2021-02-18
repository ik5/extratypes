package extratypes

import (
	"errors"
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

func TestAsInt8Int8(t *testing.T) {
	src := int8(10)
	dest := asInt8(src)

	if dest != src {
		t.Errorf("dest [%d] != src[%d]", dest, src)
	}
}

func TestAsInt8Int16Big(t *testing.T) {
	src := int16(256)
	dest := asInt8(src)

	if int16(dest) == src {
		t.Errorf("dest [%d] == src[%d]", dest, src)
	}

	if dest != maxInt8 {
		t.Errorf("dest [%d] != %d", dest, maxInt8)
	}
}

func TestAsInt8Int16Small(t *testing.T) {
	src := int16(-256)
	dest := asInt8(src)

	if int16(dest) == src {
		t.Errorf("dest [%d] == src[%d]", dest, src)
	}

	if dest != minInt8 {
		t.Errorf("dest [%d] != %d", dest, maxInt8)
	}
}

func TestAsInt8Int16Equal(t *testing.T) {
	src := int16(10)
	dest := asInt8(src)

	if int16(dest) != src {
		t.Errorf("dest [%d] != src[%d]", dest, src)
	}
}

func TestAsInt8Uint16Big(t *testing.T) {
	src := uint16(256)
	dest := asInt8(src)

	if uint16(dest) == src {
		t.Errorf("dest [%d] == src[%d]", dest, src)
	}

	if dest != maxInt8 {
		t.Errorf("dest [%d] != %d", dest, maxInt8)
	}
}

func TestAsInt8Uint16Equal(t *testing.T) {
	src := uint16(10)
	dest := asInt8(src)

	if uint16(dest) != src {
		t.Errorf("dest [%d] != src[%d]", dest, src)
	}
}

func TestAsInt8StringEmpty(t *testing.T) {
	src := ""
	dest := asInt8(src)

	if dest != 0 {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsInt8StringSignedNormal(t *testing.T) {
	src := "-10"
	dest := asInt8(src)

	if dest != -10 {
		t.Errorf("dest [%d] is not -10", dest)
	}
}

func TestAsInt8StringSignedInvalid(t *testing.T) {
	src := "-a10"
	dest := asInt8(src)

	if dest != 0 {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsInt8StringSignedBig(t *testing.T) {
	src := "-1000"
	dest := asInt8(src)

	if dest != minInt8 {
		t.Errorf("dest [%d] is not %d", dest, minInt8)
	}
}

func TestAsInt8StringUnsignedNormal(t *testing.T) {
	src := "10"
	dest := asInt8(src)

	if dest != 10 {
		t.Errorf("dest [%d] is not 10", dest)
	}
}

func TestAsInt8StringUnsignedBig(t *testing.T) {
	src := "1000"
	dest := asInt8(src)

	if dest != maxInt8 {
		t.Errorf("dest [%d] is not %d", dest, maxInt8)
	}
}

func TestAsInt8StringUnsignedInvalid(t *testing.T) {
	src := "a10"
	dest := asInt8(src)

	if dest != 0 {
		t.Errorf("dest [%d] is not 0", dest)
	}
}

func TestAsInt8Unknown(t *testing.T) {
	var src struct{}
	dest := asInt8(src)

	if dest != 0 {
		t.Errorf("dest [%d] is not 0", dest)
	}
}
