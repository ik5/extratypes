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
