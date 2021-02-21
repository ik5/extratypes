package extratypes

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var (
	validBoolTrue = Bool{
		Val: true,
		Nil: false,
	}

	validBoolFalse = Bool{
		Val: false,
		Nil: false,
	}

	validBoolNil = Bool{
		Val: false,
		Nil: true,
	}
)

var (
	testValidBoolTrueJSON  = []byte("true")
	testValidBoolFalseJSON = []byte("false")
	testValidBoolNilJSON   = []byte("null")
	testInvalidBoolJSON    = []byte("10")
)

func TestBoolString(t *testing.T) {
	t.Run("normal", func(te *testing.T) {
		list := []Bool{
			Bool{Val: false, Nil: false},
			Bool{Val: true, Nil: false},
		}

		valid := []string{
			"false", "true",
		}

		for i, item := range list {
			if item.String() != valid[i] {
				te.Errorf("item [%+v] '%s' is not '%s'", item, item, valid[i])
			}
		}
	})

	t.Run("nil", func(te *testing.T) {
		list := []Bool{
			Bool{Val: false, Nil: true},
			Bool{Val: true, Nil: true},
		}

		valid := []string{
			"nil", "nil",
		}

		for i, item := range list {
			if item.String() != valid[i] {
				te.Errorf("item [%+v] '%s' is not '%s'", item, item, valid[i])
			}
		}

	})
}

func TestBoolScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("No error was expected but have: %s", err)
	}
	defer db.Close()

	t.Run("test scan valid", func(te *testing.T) {
		rows := mock.NewRows([]string{"b"}).
			AddRow(true).
			AddRow(false)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		rs, _ := db.Query("SELECT")
		defer rs.Close()
		for rs.Next() {
			var b Bool
			err := rs.Scan(&b)
			if err != nil {
				t.Errorf("Unable to scan Bool: %s", err)
			}
			if b.Nil {
				te.Errorf("b [%t] not expected to be nil", b.Val)
			}
		}

		if rs.Err() != nil {
			te.Errorf("got rows error: %s", rs.Err())
		}
	})

	t.Run("test nil", func(te *testing.T) {
		rows := mock.NewRows([]string{"i"}).
			AddRow(nil)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		rs, _ := db.Query("SELECT")
		defer rs.Close()
		for rs.Next() {
			var b Bool
			err := rs.Scan(&b)
			if err != nil {
				t.Errorf("Unable to scan Bool: %s", err)
			}
			if !b.Nil {
				te.Errorf("b [%t] expected to be nil", b.Val)
			}
		}

		if rs.Err() != nil {
			te.Errorf("got rows error: %s", rs.Err())
		}
	})

}

func TestBoolValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database: %s", err)
		return
	}
	defer db.Close()

	mock.ExpectExec("^INSERT (.+)").WithArgs(validBoolTrue, validBoolFalse, validBoolNil).
		WillReturnResult(sqlmock.NewResult(int64(3), 1))

	_, err = db.Exec("INSERT (b)", validBoolTrue, validBoolFalse, validBoolNil)
	if err != nil {
		t.Errorf("Unable to insert record: %s", err)
	}

}

func TestBoolJSONMarshal(t *testing.T) {
	t.Run("marshal true", func(te *testing.T) {
		result, err := validBoolTrue.MarshalJSON()
		if err != nil {
			te.Errorf("Error marshaling Int to JSON: %s", err)
		}

		cmp := bytes.Compare(result, testValidBoolTrueJSON)
		if cmp != 0 {
			te.Errorf("Expected '%s', got '%s'", testValidBoolTrueJSON, result)
		}

	})

	t.Run("marshal false", func(te *testing.T) {
		result, err := validBoolFalse.MarshalJSON()

		if err != nil {
			te.Errorf("Error marshaling Int to JSON: %s", err)
		}

		cmp := bytes.Compare(result, testValidBoolFalseJSON)
		if cmp != 0 {
			te.Errorf("Expected '%s', got '%s'", testValidBoolFalseJSON, result)
		}

	})

	t.Run("marshal nil", func(te *testing.T) {
		result, err := validBoolNil.MarshalJSON()

		if err != nil {
			te.Errorf("Error marshaling Int to JSON: %s", err)
		}

		cmp := bytes.Compare(result, testValidBoolNilJSON)
		if cmp != 0 {
			te.Errorf("Expected '%s', got '%s'", testValidBoolNilJSON, result)
		}

	})
}

func TestBoolJSONUnmarshal(t *testing.T) {
	t.Run("unmarshal true", func(te *testing.T) {
		var result Bool
		err := result.UnmarshalJSON(testValidBoolTrueJSON)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, validBoolTrue)
		if !cmp {
			te.Errorf("result %+v not equal to %+v", result, validBoolTrue)
		}
	})

	t.Run("unmarshal false", func(te *testing.T) {
		var result Bool
		err := result.UnmarshalJSON(testValidBoolFalseJSON)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, validBoolFalse)
		if !cmp {
			te.Errorf("result %+v not equal to %+v", result, validBoolFalse)
		}
	})

	t.Run("unmarshal nil", func(te *testing.T) {
		var result Bool
		err := result.UnmarshalJSON(testValidBoolNilJSON)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, validBoolNil)
		if !cmp {
			te.Errorf("result %+v not equal to %+v", result, validBoolNil)
		}

	})

	t.Run("unmarshal err", func(te *testing.T) {
		var result Bool
		err := result.UnmarshalJSON(nil)

		if err == nil {
			te.Errorf("Error cannot be nil")
		}

	})
}
