package extratypes

import "database/sql/driver"

// Int struct contains int data type that can be null, and also string
// on JSON and SQL, but value will be converted to int type
type Int struct {
	Val int
	Nil bool
}

// Value interface for db
func (i Int) Value() (driver.Value, error) {
	return i.Val, nil
}

// Scan implement the Scan function from db interface
func (i *Int) Scan(v interface{}) error {
	isNil, err := toType(v, &i.Val)
	if err != nil {
		return err
	}

	i.Nil = isNil
	return nil
}
