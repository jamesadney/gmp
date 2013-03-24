package gmp

import (
	"testing"
)

func TestZeroRat(t *testing.T) {
	var x, y, z Rat
	y.SetFrac64(0, 42)

	if x.Cmp(&y) != 0 {
		t.Errorf("x and y should be both equal and zero")
	}

	if s := x.String(); s != "0/1" {
		t.Errorf("got x = %s, want 0/1", s)
	}

	if s := x.RatString(); s != "0" {
		t.Errorf("got x = %s, want 0", s)
	}

	z.Add(&x, &y)
	if s := z.RatString(); s != "0" {
		t.Errorf("got x+y = %s, want 0", s)
	}

	z.Sub(&x, &y)
	if s := z.RatString(); s != "0" {
		t.Errorf("got x-y = %s, want 0", s)
	}

	z.Mul(&x, &y)
	if s := z.RatString(); s != "0" {
		t.Errorf("got x*y = %s, want 0", s)
	}

	// FIXME: check for division by zero
	// defer func() {
	// 	if s := recover(); s == nil || s.(string) != "division by zero" {
	// 		panic(s)
	// 	}
	// }()
	// z.Quo(&x, &y)
}

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

type ratBinFun func(z, x, y *Rat) *Rat
type ratBinArg struct {
	x, y, z string
}

func testRatBin(t *testing.T, i int, name string, f ratBinFun, a ratBinArg) {
	x, _ := new(Rat).SetString(a.x)
	y, _ := new(Rat).SetString(a.y)
	z, _ := new(Rat).SetString(a.z)
	out := f(new(Rat), x, y)

	if out.Cmp(z) != 0 {
		t.Errorf("%s #%d got %s want %s", name, i, out, z)
	}
}

var ratBinTests = []struct {
	x, y      string
	sum, prod string
}{
	{"0", "0", "0", "0"},
	{"0", "1", "1", "0"},
	{"-1", "0", "-1", "0"},
	{"-1", "1", "0", "-1"},
	{"1", "1", "2", "1"},
	{"1/2", "1/2", "1", "1/4"},
	{"1/4", "1/3", "7/12", "1/12"},
	{"2/5", "-14/3", "-64/15", "-28/15"},
	{"4707/49292519774798173060", "-3367/70976135186689855734", "84058377121001851123459/1749296273614329067191168098769082663020", "-1760941/388732505247628681598037355282018369560"},
	{"-61204110018146728334/3", "-31052192278051565633/2", "-215564796870448153567/6", "950260896245257153059642991192710872711/3"},
	{"-854857841473707320655/4237645934602118692642972629634714039", "-18/31750379913563777419", "-27/133467566250814981", "15387441146526731771790/134546868362786310073779084329032722548987800600710485341"},
	{"618575745270541348005638912139/19198433543745179392300736", "-19948846211000086/637313996471", "27674141753240653/30123979153216", "-6169936206128396568797607742807090270137721977/6117715203873571641674006593837351328"},
	{"-3/26206484091896184128", "5/2848423294177090248", "15310893822118706237/9330894968229805033368778458685147968", "-5/24882386581946146755650075889827061248"},
	{"26946729/330400702820", "41563965/225583428284", "1238218672302860271/4658307703098666660055", "224002580204097/14906584649915733312176"},
	{"-8259900599013409474/7", "-84829337473700364773/56707961321161574960", "-468402123685491748914621885145127724451/396955729248131024720", "350340947706464153265156004876107029701/198477864624065512360"},
	{"575775209696864/1320203974639986246357", "29/712593081308", "410331716733912717985762465/940768218243776489278275419794956", "808/45524274987585732633"},
	{"1786597389946320496771/2066653520653241", "6269770/1992362624741777", "3559549865190272133656109052308126637/4117523232840525481453983149257", "8967230/3296219033"},
	{"-36459180403360509753/32150500941194292113930", "9381566963714/9633539", "301622077145533298008420642898530153/309723104686531919656937098270", "-3784609207827/3426986245"},
}

func TestRatBin(t *testing.T) {
	for i, test := range ratBinTests {
		arg := ratBinArg{test.x, test.y, test.sum}
		testRatBin(t, i, "Add", (*Rat).Add, arg)

		arg = ratBinArg{test.y, test.x, test.sum}
		testRatBin(t, i, "Add symmetric", (*Rat).Add, arg)

		arg = ratBinArg{test.sum, test.x, test.y}
		testRatBin(t, i, "Sub", (*Rat).Sub, arg)

		arg = ratBinArg{test.sum, test.y, test.x}
		testRatBin(t, i, "Sub symmetric", (*Rat).Sub, arg)

		arg = ratBinArg{test.x, test.y, test.prod}
		testRatBin(t, i, "Mul", (*Rat).Mul, arg)

		arg = ratBinArg{test.y, test.x, test.prod}
		testRatBin(t, i, "Mul symmetric", (*Rat).Mul, arg)

		if test.x != "0" {
			arg = ratBinArg{test.prod, test.x, test.y}
			testRatBin(t, i, "Quo", (*Rat).Quo, arg)
		}

		if test.y != "0" {
			arg = ratBinArg{test.prod, test.y, test.x}
			testRatBin(t, i, "Quo symmetric", (*Rat).Quo, arg)
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
