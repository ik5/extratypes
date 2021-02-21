package extratypes

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type intForTestingConvertion struct {
	I Int `json:"i" toml:"i"`
}

var (
	validInt = Int{
		Val: 10,
		Nil: false,
	}

	validMinusInt = Int{
		Val: -10,
		Nil: false,
	}

	nilInt = Int{
		Val: 0,
		Nil: true,
	}
)

var (
	testIntJSONError      = []byte("10")
	testIntMinusJSONError = []byte("-10")
	testIntNilJSONError   = []byte("null")
	testIntErrJSONError   = []byte("a")
	testIntValidText      = []byte("10")
	testIntMinuxValidText = []byte("-10")
	testIntNilText        = []byte("")
	testIntErrText        = []byte("a")
)

func TestIntString(t *testing.T) {
	i := Int{
		Val: 10,
		Nil: false,
	}

	if i.String() != "10" {
		t.Errorf("i [%s] is not 10", i)
	}

	i = Int{
		Val: 0,
		Nil: true,
	}

	if i.String() != "nil" {
		t.Errorf("i [%s] is not nil", i)
	}
}

func TestIntScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("No error was expected but have: %s", err)
	}
	defer db.Close()

	t.Run("test scan valid", func(te *testing.T) {
		rows := mock.NewRows([]string{"i"}).
			AddRow(-1).
			AddRow(1).
			AddRow(1.1)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		rs, _ := db.Query("SELECT")
		defer rs.Close()
		for rs.Next() {
			var i Int
			err := rs.Scan(&i)
			if err != nil {
				t.Errorf("Unable to scan Int: %s", err)
			}
			if i.Nil {
				te.Errorf("i [%d] not expected to be nil", i.Val)
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
			var i Int
			err := rs.Scan(&i)
			if err != nil {
				t.Errorf("Unable to scan Int: %s", err)
			}
			if !i.Nil {
				te.Errorf("i [%d] expected to be nil", i.Val)
			}
		}

		if rs.Err() != nil {
			te.Errorf("got rows error: %s", rs.Err())
		}
	})
}

func TestIntValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database: %s", err)
		return
	}
	defer db.Close()

	mock.ExpectExec("^INSERT (.+)").WithArgs(validInt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = db.Exec("INSERT (i)", validInt)
	if err != nil {
		t.Errorf("Unable to insert record: %s", err)
	}

}

func TestIntJSONMarshal(t *testing.T) {
	t.Run("marshal int", func(te *testing.T) {
		result, err := validInt.MarshalJSON()
		if err != nil {
			te.Errorf("Error marshaling Int to JSON: %s", err)
		}

		cmp := bytes.Compare(result, testIntJSONError)
		if cmp != 0 {
			te.Errorf("Expected '%s', got '%s'", testIntJSONError, result)
		}

	})

	t.Run("marshal minus", func(te *testing.T) {
		result, err := validMinusInt.MarshalJSON()

		if err != nil {
			te.Errorf("Error marshaling Int to JSON: %s", err)
		}

		cmp := bytes.Compare(result, testIntMinusJSONError)
		if cmp != 0 {
			te.Errorf("Expected '%s', got '%s'", testIntNilJSONError, result)
		}

	})

	t.Run("marshal nil", func(te *testing.T) {
		result, err := nilInt.MarshalJSON()

		if err != nil {
			te.Errorf("Error marshaling Int to JSON: %s", err)
		}

		cmp := bytes.Compare(result, testIntNilJSONError)
		if cmp != 0 {
			te.Errorf("Expected '%s', got '%s'", testIntNilJSONError, result)
		}

	})
}

func TestIntJSONUnmarshal(t *testing.T) {
	t.Run("unmarshal int", func(te *testing.T) {
		var result Int
		err := result.UnmarshalJSON(testIntJSONError)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, validInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, validInt)
		}
	})

	t.Run("unmarshal minus int", func(te *testing.T) {
		var result Int
		err := result.UnmarshalJSON(testIntMinusJSONError)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, validMinusInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, validMinusInt)
		}
	})

	t.Run("unmarshal nil", func(te *testing.T) {
		var result Int
		err := result.UnmarshalJSON(testIntNilJSONError)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, nilInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, nilInt)
		}

	})

	t.Run("unmarshal err", func(te *testing.T) {
		var result Int
		err := result.UnmarshalJSON(testIntErrJSONError)

		if err == nil {
			te.Errorf("expected error, but none given")
		}

	})
}

func TestIntTextMarshal(t *testing.T) {
	t.Run("marshal valid", func(te *testing.T) {
		b, err := validInt.MarshalText()
		if err != nil {
			te.Errorf("Unexpected err: %s", err)
		}

		cmp := bytes.Compare(b, testIntValidText)
		if cmp != 0 {
			t.Errorf("b %s is not %s", b, testIntValidText)
		}
	})

	t.Run("text minus valid", func(te *testing.T) {
		b, err := validMinusInt.MarshalText()
		if err != nil {
			te.Errorf("Unexpected err: %s", err)
		}

		cmp := bytes.Compare(b, testIntMinuxValidText)
		if cmp != 0 {
			t.Errorf("b %s is not %s", b, testIntMinuxValidText)
		}

	})

	t.Run("test nil marshal", func(te *testing.T) {
		b, err := nilInt.MarshalText()
		if err != nil {
			te.Errorf("Unexpected err: %s", err)
		}

		cmp := bytes.Compare(b, testIntNilText)
		if cmp != 0 {
			t.Errorf("b %s is not %s", b, testIntNilText)
		}

	})
}
func TestIntTextUnmarshal(t *testing.T) {
	t.Run("unmarshal int", func(te *testing.T) {
		var result Int
		err := result.UnmarshalJSON(testIntValidText)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, validInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, validInt)
		}
	})

	t.Run("unmarshal minus int", func(te *testing.T) {
		var result Int
		err := result.UnmarshalText(testIntMinuxValidText)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, validMinusInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, validMinusInt)
		}
	})

	t.Run("unmarshal nil", func(te *testing.T) {
		var result Int
		err := result.UnmarshalText(testIntNilText)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, nilInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, nilInt)
		}

		err = result.UnmarshalText(nil)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp = reflect.DeepEqual(result, nilInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, nilInt)
		}

	})

}
