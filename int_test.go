package extratypes

import "testing"

type intForTestingConvertion struct {
	I Int `json:"i" toml:"i"`
}

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
