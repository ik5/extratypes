package extratypes

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestInvalidRateStructure(t *testing.T) {
	raw := "0"

	_, err := NewRate(raw)
	if err == nil {
		t.Error("Rate should have returned an error, but it is nil")
		return
	}

	if !strings.Contains(
		err.Error(), fmt.Sprintf("Invalid Rate structure for %s", raw),
	) {
		t.Errorf("Expected Invalid Rate structure, got %s instead", err)
		return
	}
}

func TestRateInvalidParseUInt(t *testing.T) {
	raw := "1d/5"

	_, err := NewRate(raw)

	if err == nil {
		t.Error("Rate should have returned an error, but it is nil")
		return
	}

	if !errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("Expected: %+v, but got: %+v", strconv.ErrSyntax, err)
		return
	}
}

func TestRateInvalidDuration(t *testing.T) {
	raw := "10/l"

	_, err := NewRate(raw)
	if err == nil {
		t.Error("Rate should have returned an error, but it is nil")
		return
	}

	if !strings.HasPrefix(err.Error(), "time: invalid duration") {
		t.Errorf("Expected error of 'time: invalid duration l', but got: %s", err)
		return
	}

}

func TestMarshalJSON(t *testing.T) {
	r, err := NewRate("100/1s")
	if err != nil {
		t.Errorf("Unable to call NewRate: %s", err)
		return
	}

	jStr, err := r.MarshalJSON()
	if err != nil {
		t.Errorf("Unable to marshal json: %s", err)
		return
	}

	result := `"100/1s"`

	c := strings.Compare(string(jStr), result)
	if c != 0 {
		t.Errorf("Invalid marshal: %d (%s | %s)", c, jStr, result)
	}
}

func TestMarshalJSONEmpty(t *testing.T) {
	rate := Rate{}

	b, err := rate.MarshalJSON()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	if len(b) == 0 {
		t.Errorf("Expected empty JSON string, got: %s", b)
		return
	}

	if bytes.Compare(b, []byte(`""`)) != 0 {
		t.Errorf(`Expected "" but got: %s`, b)
		return
	}
}

func TestUnmarshalJSON(t *testing.T) {
	r, err := NewRate("")
	if err != nil {
		t.Errorf("Unable to generate NewRate: %s", err)
		return
	}
	jStr := []byte(`"100/1s"`)
	err = r.UnmarshalJSON(jStr)

	if err != nil {
		t.Errorf("Unable to unmarshal %s: %s", jStr, err)
		return
	}

	if r.IsEmpty() {
		t.Error("r is empty")
		return
	}

	if !(r.Count == 100 && r.Sleep.Duration == time.Second) {
		t.Errorf("Invalid content, expected 100/1s got: %+v", r)
	}
}
