package gmp

import (
	"testing"
)

var setStringTests = []struct {
	in, out string
	ok      bool
}{
	{"0", "0", true},
	{"-0", "0", true},
	{"1", "1", true},
	{"-1", "-1", true},
	{"1.", "1", true},
	{"1e0", "1", true},
	{"1.e1", "10", true},
	{in: "1e", ok: false},
	{in: "1.e", ok: false},
	{in: "1e+14e-5", ok: false},
	{in: "1e4.5", ok: false},
	{in: "r", ok: false},
	{in: "a/b", ok: false},
	{in: "a.b", ok: false},
	{"-0.1", "-1/10", true},
	{"-.1", "-1/10", true},
	{"2/4", "1/2", true},
	{".25", "1/4", true},
	{"-1/5", "-1/5", true},
	{"8129567.7690E14", "812956776900000000000", true},
	{"78189e+4", "781890000", true},
	{"553019.8935e+8", "55301989350000", true},
	{"98765432109876543210987654321e-10", "98765432109876543210987654321/10000000000", true},
	{"9877861857500000E-7", "3951144743/4", true},
	{"2169378.417e-3", "2169378417/1000000", true},
	{"884243222337379604041632732738665534", "884243222337379604041632732738665534", true},
	{"53/70893980658822810696", "53/70893980658822810696", true},
	{"106/141787961317645621392", "53/70893980658822810696", true},
	{"204211327800791583.81095", "4084226556015831676219/20000", true},
}

func TestRatSetString(t *testing.T) {
	for i, test := range setStringTests {
		x, ok := new(Rat).SetString(test.in)

		if ok {
			if !test.ok {
				t.Errorf("#%d SetString(%q) expected failure", i, test.in)
			} else if x.RatString() != test.out {
				t.Errorf("#%d SetString(%q) got %s want %s", i, test.in, x.RatString(), test.out)
			}
		} else if x != nil {
			t.Errorf("#%d SetString(%q) got %p want nil", i, test.in, x)
		}
	}
}

var ratCmpTests = []struct {
	rat1, rat2 string
	out        int
}{
	{"0", "0/1", 0},
	{"1/1", "1", 0},
	{"-1", "-2/2", 0},
	{"1", "0", 1},
	{"0/1", "1/1", -1},
	{"-5/1434770811533343057144", "-5/1434770811533343057145", -1},
	{"49832350382626108453/8964749413", "49832350382626108454/8964749413", -1},
	{"-37414950961700930/7204075375675961", "37414950961700930/7204075375675961", -1},
	{"37414950961700930/7204075375675961", "74829901923401860/14408150751351922", 0},
}

func TestRatCmp(t *testing.T) {
	for i, test := range ratCmpTests {
		x, _ := new(Rat).SetString(test.rat1)
		y, _ := new(Rat).SetString(test.rat2)

		out := x.Cmp(y)
		if out != test.out {
			t.Errorf("#%d got out = %v; want %v", i, out, test.out)
		}
	}
}

var getStringTests = []struct {
	in, out string
	ok      bool
}{
	{"1/1", "1/1", true},
	{"-1/1", "-1/1", true},
	{"2/1", "2/1", true},
	{"4/2", "2/1", true},
}

func TestGetString(t *testing.T) {
	for i, test := range getStringTests {
		x, _ := new(Rat).SetString(test.in)
		if x.String() != test.out {
			t.Errorf("#%d String() got %s want %s", i, x.String(), test.out)
		}
	}
}

func TestNumDenomAreReferences(t *testing.T) {
	x := NewRat(1, 2)
	n := x.Num()
	d := x.Denom()

	x.Add(x, NewRat(1, 4))
	if n.Cmp(NewInt(3)) != 0 {
		t.Error("*Int returned by q.Num() is not a reference to the num. of q.")
	}
	if d.Cmp(NewInt(4)) != 0 {
		t.Error("*Int returned by q.Denom() is not a reference to the den. of q.")
	}
}

var setFrac64Tests = []struct {
	x, y int64
	out  string
}{
	{1, 2, "1/2"},
	{1, -2, "-1/2"},
	{-1, 2, "-1/2"},
	{-1, -2, "1/2"},
	{2, 4, "1/2"},
}

func TestSetFrac64(t *testing.T) {
	for i, test := range setFrac64Tests {
		q := new(Rat).SetFrac64(test.x, test.y)
		if q.String() != test.out {
			t.Errorf("#%d SetFrac64(%d, %d) got %s want %s", i, test.x, test.y, q.String(), test.out)
		}
	}
}

var setInt64Tests = []struct {
	x   int64
	out string
}{
	{1, "1"},
	{-1, "-1"},
}

func TestSetInt64(t *testing.T) {
	for i, test := range setInt64Tests {
		q := new(Rat).SetInt64(test.x)
		if q.RatString() != test.out {
			t.Errorf("#%d SetInt64(%d) got %s want %s", i, test.x, q.RatString(), test.out)
		}
	}
}
