package extratypes

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type durationForTestingMaps struct {
	Duration Duration `json:"d" toml:"d"`
}

var (
	durationStruct = Duration{
		Duration: time.Second,
		Nil:      false,
	}

	durationStructAsMap = durationForTestingMaps{
		Duration: Duration{
			time.Second,
			false,
		},
	}
)

var (
	testDurationJSONStr                     = []byte(`"1s"`)
	testDurationJSONInt                     = []byte(`1000000000`)
	testDurationJSONInvalidDuration         = []byte(`true`)
	testDurationJSONInvalidJSON             = []byte(`{"d":}`)
	testDurationJSONParsingDurationError    = []byte(`"1x"`)
	testDurationJSONParsingDurationErrorNaN = []byte(`"NaN"`)
	testDurationMapJSONStr                  = []byte(`{"d":"1s"}`)
	testDurationMapJSONInt                  = []byte(`{"d":1000000000}`)
	testDurationJSONInvalidDurationParse    = []byte(`{"d": "1x"}`)
	testDurationJSONNoContentFound          = []byte(`{}`)
	testDurationJSONInvalidDataType         = []byte(`{"d": true}`)
	testDurationJSONLengthTooBig            = []byte(`{"d": true, "x": "1s"}`)
	testDurationNil                         = []byte(`null`)
	testDurationNilMap                      = []byte(`{"d": null}`)

	testDurationTextBasic = []byte(`1s`)
	testDurationTextErr   = []byte(`true`)
	testDurationTextNil   = []byte("")
)

func TestScan(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("expected no error, but got: %s", err)
		return
	}
	defer db.Close()

	row := mock.NewRows([]string{"duration"}).
		AddRow("1s").
		AddRow(1000000000).
		AddRow(1000000000.0)

	mock.ExpectQuery("SELECT").WillReturnRows(row)
	rs, _ := db.Query("SELECT")
	defer rs.Close()
	for rs.Next() {
		var d Duration
		rs.Scan(&d)
		if d.Duration != time.Second {
			t.Errorf("Duration is not 1 second: %s", d)
		}

	}

	if rs.Err() != nil {
		t.Errorf("got rows error: %s", rs.Err())
	}
}

func TestScanNil(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("expected no error, but got: %s", err)
		return
	}
	defer db.Close()

	row := mock.NewRows([]string{"duration"}).
		AddRow(0).
		AddRow("").
		AddRow(nil)

	mock.ExpectQuery("SELECT").WillReturnRows(row)
	rs, _ := db.Query("SELECT")
	defer rs.Close()
	for rs.Next() {
		var d Duration
		rs.Scan(&d)
		if d.Duration != 0 {
			t.Errorf("Duration is not 1 second: %s", d)
		}

	}

	if rs.Err() != nil {
		t.Errorf("got rows error: %s", rs.Err())
	}
}

func TestScanInvalidValue(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("expected no error, but got: %s", err)
		return
	}
	defer db.Close()

	row := mock.NewRows([]string{"duration"}).
		AddRow(true).
		AddRow(false)

	mock.ExpectQuery("SELECT").WillReturnRows(row)
	rs, _ := db.Query("SELECT")
	defer rs.Close()
	for rs.Next() {
		var d Duration
		err = rs.Scan(&d)
		if err == nil {
			t.Error("There was no error on Scan")
		}

	}

	if rs.Err() != nil {
		t.Errorf("got rows error: %s", rs.Err())
	}
}

func TestValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database: %s", err)
		return
	}
	defer db.Close()

	mock.ExpectExec("^INSERT (.+)").WithArgs(durationStruct).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = db.Exec("INSERT (d)", durationStruct)
	if err != nil {
		t.Errorf("Unable to insert record: %s", err)
	}

}

func TestJSONMarshal(t *testing.T) {
	result, err := durationStruct.MarshalJSON()

	if err != nil {
		t.Errorf("Error marshaling duration to JSON: %s", err)
		return
	}

	cmp := bytes.Compare(result, testDurationJSONStr)
	if cmp != 0 {
		t.Errorf("Expected '%s' to equal '%s', but found %d is different",
			result, testDurationJSONStr, cmp)
	}
}

func TestJSONMarshalAsMap(t *testing.T) {
	result, err := json.Marshal(durationStructAsMap)

	if err != nil {
		t.Errorf("Error marshaling duration to JSON: %s", err)
		return
	}

	cmp := bytes.Compare(result, testDurationMapJSONStr)
	if cmp != 0 {
		t.Errorf("Expected '%s' to equal '%s', but found %d is different",
			result, testDurationMapJSONStr, cmp)
	}
}

func TestJSONUnMarshalAsMap(t *testing.T) {
	result := Duration{}
	err := result.UnmarshalJSON(testDurationMapJSONStr)

	if err != nil {
		t.Errorf("Unable to marshal %s: %s", testDurationJSONStr, err)
		return
	}

	cmp := reflect.DeepEqual(result, durationStruct)
	if !cmp {
		t.Errorf("Expected %+v to equal %+v", result, durationStruct)
	}
}

func TestJSONUnMarshalAsMapInt(t *testing.T) {
	result := Duration{}
	err := result.UnmarshalJSON(testDurationMapJSONInt)

	if err != nil {
		t.Errorf("Unable to marshal %s: %s", testDurationJSONInt, err)
		return
	}

	cmp := reflect.DeepEqual(result, durationStruct)
	if !cmp {
		t.Errorf("Expected %+v to equal %+v", result, durationStruct)
	}
}

func TestJSONUnMarshalStr(t *testing.T) {
	result := Duration{}
	err := result.UnmarshalJSON(testDurationJSONStr)

	if err != nil {
		t.Errorf("Unable to marshal %s: %s", testDurationJSONStr, err)
		return
	}

	cmp := reflect.DeepEqual(result, durationStruct)
	if !cmp {
		t.Errorf("Expected %+v to equal %+v", result, durationStruct)
	}

}

func TestJSONUnMarshalInt(t *testing.T) {
	var result Duration
	err := result.UnmarshalJSON(testDurationJSONInt)

	if err != nil {
		t.Errorf("Unable to marshal %s: %s", testDurationJSONInt, err)
		return
	}

	cmp := reflect.DeepEqual(result, durationStruct)
	if !cmp {
		t.Errorf("Expected %+v to equal %+v", result, durationStruct)
	}

}

func TestJSONUnMarshalInvalidDuration(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalJSON(testDurationJSONInvalidDuration)
	if err == nil {
		t.Errorf("Expected err, but got nil")
		return
	}

	if !reflect.DeepEqual(d, Duration{}) {
		t.Errorf("Expected empty duration, got %+v instead", d)
		return
	}
}

func TestJSONUnMarshalInvalidJSON(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalJSON(testDurationJSONInvalidJSON)
	if err == nil {
		t.Errorf("Expected err, but got nil")
		return
	}

	if !reflect.DeepEqual(d, Duration{}) {
		t.Errorf("Expected empty duration, got %+v instead", d)
		return
	}
}
func TestJSONUnMarshalParsingDurationError(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalJSON(testDurationJSONParsingDurationError)
	if err == nil {
		t.Errorf("Expected err, but got nil")
		return
	}

	if !reflect.DeepEqual(d, Duration{}) {
		t.Errorf("Expected empty duration, got %+v instead", d)
		return
	}
}
func TestJSONUnMarshalParsingDurationErrorNaN(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalJSON(testDurationJSONParsingDurationErrorNaN)
	if err == nil {
		t.Errorf("Expected err, but got nil")
		return
	}

	if !reflect.DeepEqual(d, Duration{}) {
		t.Errorf("Expected empty duration, got %+v instead", d)
		return
	}
}

func TestJSONUnMarshalInvalidDurationParse(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalJSON(testDurationJSONInvalidDurationParse)
	if err == nil {
		t.Errorf("Expected err, but got nil")
		return
	}

	if !reflect.DeepEqual(d, Duration{}) {
		t.Errorf("Expected empty duration, got %+v instead", d)
		return
	}
}

func TestJSONUnMarshalNoContentFound(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalJSON(testDurationJSONNoContentFound)
	if err == nil {
		t.Errorf("Expected err, but got nil")
		return
	}

	if !reflect.DeepEqual(d, Duration{}) {
		t.Errorf("Expected empty duration, got %+v instead", d)
		return
	}
}

func TestJSONUnMarshalInvalidDataType(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalJSON(testDurationJSONInvalidDataType)
	if err == nil {
		t.Errorf("Expected err, but got nil")
		return
	}

	if !reflect.DeepEqual(d, Duration{}) {
		t.Errorf("Expected empty duration, got %+v instead", d)
		return
	}
}

func TestJSONUnMarshalLengthTooBig(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalJSON(testDurationJSONLengthTooBig)
	if err == nil {
		t.Errorf("Expected err, but got nil")
		return
	}

	if !reflect.DeepEqual(d, Duration{}) {
		t.Errorf("Expected empty duration, got %+v instead", d)
		return
	}
}

func TestMarshalJSONNil(t *testing.T) {
	d := Duration{
		Duration: 0,
		Nil:      true,
	}

	result, err := d.MarshalJSON()
	if err != nil {
		t.Errorf("Error marshal nil json: %s", err)
		return
	}

	cmp := bytes.Compare(result, testDurationNil)
	if cmp != 0 {
		t.Errorf("Invalid result: %d %s", cmp, result)
	}
}

func TestUnmarshalJSONBasicNil(t *testing.T) {
	result := Duration{}
	err := result.UnmarshalJSON(testDurationNil)
	if err != nil {
		t.Errorf("Unable to unmarshal nil duration in JSON: %s", err)
		return
	}

	if result.Duration != 0 {
		t.Errorf("Expected duration to be 0, found %d", result.Duration)
	}

	if !result.Nil {
		t.Errorf("Expected nil")
	}
}

func TestUnmarshalJSONMapNil(t *testing.T) {
	result := Duration{}
	err := result.UnmarshalJSON(testDurationNilMap)

	if err != nil {
		t.Errorf("Unable to unmarshal nil duration in JSON: %s", err)
		return
	}

	if result.Duration != 0 {
		t.Errorf("Expected duration to be 0, found %d", result.Duration)
	}

	if !result.Nil {
		t.Errorf("Expected nil")
	}

}

func TestNilString(t *testing.T) {
	d := Duration{
		Duration: 0,
		Nil:      true,
	}

	if d.String() != "nil" {
		t.Errorf("Expected nil, found %s", d.String())
	}
}

func TestMarshalText(t *testing.T) {
	result, err := durationStruct.MarshalText()
	if err != nil {
		t.Errorf("Unable to marshal text: %s", err)
		return
	}

	cmp := bytes.Compare(result, testDurationTextBasic)
	if cmp != 0 {
		t.Errorf("cmp: %d result: %s", cmp, result)
	}
}

func TestMarshalTextNil(t *testing.T) {
	d := Duration{
		Duration: 0,
		Nil:      true,
	}

	result, err := d.MarshalText()
	if err != nil {
		t.Errorf("Unable to marshal text: %s", err)
		return
	}

	cmp := bytes.Compare(result, testDurationTextNil)
	if cmp != 0 {
		t.Errorf("cmp: %d result: %s", cmp, result)
	}
}

func TestUnmarshalBasic(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalText(testDurationTextBasic)
	if err != nil {
		t.Errorf("Unable to unmarshal text: %s", err)
		return
	}

	if d.Duration != time.Second {
		t.Errorf("Expected one second, got: %d (%s)", d.Duration, d.Duration)
	}

	if d.Nil {
		t.Errorf("Unexpected nil")
	}
}

func TestUnmarshalTextError(t *testing.T) {
	d := Duration{}

	err := d.UnmarshalText(testDurationTextErr)
	if err == nil {
		t.Error("Error must not be nil")
	}

	expectedErr := "time: invalid duration " + string(testDurationTextErr)

	if err.Error() != expectedErr {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	if d.Duration != 0 {
		t.Errorf("Expected duration 0, got %d (%s)", d.Duration, d.Duration)
	}

	if !d.Nil {
		t.Errorf("Expected nil")
	}
}

func TestUnmarshalTextNil(t *testing.T) {
	d := Duration{}
	err := d.UnmarshalText(testDurationTextNil)
	if err != nil {
		t.Errorf("Error unmarshal empty duration: %s", err)
		return
	}

	if !d.Nil {
		t.Errorf("Expected nil")
	}

	if d.Duration != 0 {
		t.Errorf("Expected 0, got %d (%s)", d.Duration, d.Duration)
	}
}
