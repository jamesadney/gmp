package gmp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
	"testing/quick"
)

type funZZ func(z, x, y *Int) *Int
type argZZ struct {
	z, x, y *Int
}

var sumZZ = []argZZ{
	{NewInt(0), NewInt(0), NewInt(0)},
	{NewInt(1), NewInt(1), NewInt(0)},
	{NewInt(1111111110), NewInt(123456789), NewInt(987654321)},
	{NewInt(-1), NewInt(-1), NewInt(0)},
	{NewInt(864197532), NewInt(-123456789), NewInt(987654321)},
	{NewInt(-1111111110), NewInt(-123456789), NewInt(-987654321)},
}

var prodZZ = []argZZ{
	{NewInt(0), NewInt(0), NewInt(0)},
	{NewInt(0), NewInt(1), NewInt(0)},
	{NewInt(1), NewInt(1), NewInt(1)},
	{NewInt(-991 * 991), NewInt(991), NewInt(-991)},
	// TODO(gri) add larger products
}

func TestSignZ(t *testing.T) {
	var zero Int
	for _, a := range sumZZ {
		s := a.z.Sign()
		e := a.z.Cmp(&zero)
		if s != e {
			t.Errorf("got %d; want %d for z = %v", s, e, a.z)
		}
	}
}

func TestSetZ(t *testing.T) {
	for _, a := range sumZZ {
		var z Int
		z.Set(a.z)
		if (&z).Cmp(a.z) != 0 {
			t.Errorf("got z = %v; want %v", z, a.z)
		}
	}
}

func TestAbsZ(t *testing.T) {
	var zero Int
	for _, a := range sumZZ {
		var z Int
		z.Abs(a.z)
		var e Int
		e.Set(a.z)
		if e.Cmp(&zero) < 0 {
			e.Sub(&zero, &e)
		}
		if z.Cmp(&e) != 0 {
			t.Errorf("got z = %v; want %v", z, e)
		}
	}
}

func testFunZZ(t *testing.T, msg string, f funZZ, a argZZ) {
	var z Int
	f(&z, a.x, a.y)
	if (&z).Cmp(a.z) != 0 {
		t.Errorf("%s%+v\n\tgot z = %v; want %v", msg, a, &z, a.z)
	}
}

func TestSumZZ(t *testing.T) {
	AddZZ := func(z, x, y *Int) *Int { return z.Add(x, y) }
	SubZZ := func(z, x, y *Int) *Int { return z.Sub(x, y) }
	for _, a := range sumZZ {
		arg := a
		testFunZZ(t, "AddZZ", AddZZ, arg)

		arg = argZZ{a.z, a.y, a.x}
		testFunZZ(t, "AddZZ symmetric", AddZZ, arg)

		arg = argZZ{a.x, a.z, a.y}
		testFunZZ(t, "SubZZ", SubZZ, arg)

		arg = argZZ{a.y, a.z, a.x}
		testFunZZ(t, "SubZZ symmetric", SubZZ, arg)
	}
}

func TestProdZZ(t *testing.T) {
	MulZZ := func(z, x, y *Int) *Int { return z.Mul(x, y) }
	for _, a := range prodZZ {
		arg := a
		testFunZZ(t, "MulZZ", MulZZ, arg)

		arg = argZZ{a.z, a.y, a.x}
		testFunZZ(t, "MulZZ symmetric", MulZZ, arg)
	}
}

// mulBytes returns x*y via grade school multiplication. Both inputs
// and the result are assumed to be in big-endian representation (to
// match the semantics of Int.Bytes and Int.SetBytes).
func mulBytes(x, y []byte) []byte {
	z := make([]byte, len(x)+len(y))

	// multiply
	k0 := len(z) - 1
	for j := len(y) - 1; j >= 0; j-- {
		d := int(y[j])
		if d != 0 {
			k := k0
			carry := 0
			for i := len(x) - 1; i >= 0; i-- {
				t := int(z[k]) + int(x[i])*d + carry
				z[k], carry = byte(t), t>>8
				k--
			}
			z[k] = byte(carry)
		}
		k0--
	}

	// normalize (remove leading 0's)
	i := 0
	for i < len(z) && z[i] == 0 {
		i++
	}

	return z[i:]
}

func checkMul(a, b []byte) bool {
	var x, y, z1 Int
	x.SetBytes(a)
	y.SetBytes(b)
	z1.Mul(&x, &y)

	var z2 Int
	z2.SetBytes(mulBytes(a, b))

	return z1.Cmp(&z2) == 0
}

func TestMul(t *testing.T) {
	if err := quick.Check(checkMul, nil); err != nil {
		t.Error(err)
	}
}

var mulRangesN = []struct {
	a, b uint64
	prod string
}{
	{0, 0, "0"},
	{1, 1, "1"},
	{1, 2, "2"},
	{1, 3, "6"},
	{10, 10, "10"},
	{0, 100, "0"},
	{0, 1e9, "0"},
	{1, 0, "1"},                    // empty range
	{100, 1, "1"},                  // empty range
	{1, 10, "3628800"},             // 10!
	{1, 20, "2432902008176640000"}, // 20!
	{1, 100,
		"933262154439441526816992388562667004907159682643816214685929" +
			"638952175999932299156089414639761565182862536979208272237582" +
			"51185210916864000000000000000000000000", // 100!
	},
}

var mulRangesZ = []struct {
	a, b int64
	prod string
}{
	// entirely positive ranges are covered by mulRangesN
	{-1, 1, "0"},
	{-2, -1, "2"},
	{-3, -2, "6"},
	{-3, -1, "-6"},
	{1, 3, "6"},
	{-10, -10, "-10"},
	{0, -1, "1"},                      // empty range
	{-1, -100, "1"},                   // empty range
	{-1, 1, "0"},                      // range includes 0
	{-1e9, 0, "0"},                    // range includes 0
	{-1e9, 1e9, "0"},                  // range includes 0
	{-10, -1, "3628800"},              // 10!
	{-20, -2, "-2432902008176640000"}, // -20!
	{-99, -1,
		"-933262154439441526816992388562667004907159682643816214685929" +
			"638952175999932299156089414639761565182862536979208272237582" +
			"511852109168640000000000000000000000", // -99!
	},
}

func TestMulRangeZ(t *testing.T) {
	var tmp Int
	// test entirely positive ranges
	for i, r := range mulRangesN {
		prod := tmp.MulRange(int64(r.a), int64(r.b)).String()
		if prod != r.prod {
			t.Errorf("#%da: got %s; want %s", i, prod, r.prod)
		}
	}
	// test other ranges
	for i, r := range mulRangesZ {
		prod := tmp.MulRange(r.a, r.b).String()
		if prod != r.prod {
			t.Errorf("#%db: got %s; want %s", i, prod, r.prod)
		}
	}
}

// Examples from the Go Language Spec, section "Arithmetic operators"
var divisionSignsTests = []struct {
	x, y int64
	q, r int64 // T-division
	d, m int64 // Euclidian division
}{
	{5, 3, 1, 2, 1, 2},
	{-5, 3, -1, -2, -2, 1},
	{5, -3, -1, 2, -1, 2},
	{-5, -3, 1, -2, 2, 1},
	{1, 2, 0, 1, 0, 1},
	{8, 4, 2, 0, 2, 0},
}

func TestDivisionSigns(t *testing.T) {
	for i, test := range divisionSignsTests {
		x := NewInt(test.x)
		y := NewInt(test.y)
		q := NewInt(test.q)
		r := NewInt(test.r)
		d := NewInt(test.d)
		m := NewInt(test.m)

		q1 := new(Int).Quo(x, y)
		r1 := new(Int).Rem(x, y)
		if q1.Cmp(q) != 0 || r1.Cmp(r) != 0 {
			t.Errorf("#%d QuoRem: got (%s, %s), want (%s, %s)", i, q1, r1, q, r)
		}

		q2, r2 := new(Int).QuoRem(x, y, new(Int))
		if q2.Cmp(q) != 0 || r2.Cmp(r) != 0 {
			t.Errorf("#%d QuoRem: got (%s, %s), want (%s, %s)", i, q2, r2, q, r)
		}

		d1 := new(Int).Div(x, y)
		m1 := new(Int).Mod(x, y)
		if d1.Cmp(d) != 0 || m1.Cmp(m) != 0 {
			t.Errorf("#%d DivMod: got (%s, %s), want (%s, %s)", i, d1, m1, d, m)
		}

		d2, m2 := new(Int).DivMod(x, y, new(Int))
		if d2.Cmp(d) != 0 || m2.Cmp(m) != 0 {
			t.Errorf("#%d DivMod: got (%s, %s), want (%s, %s)", i, d2, m2, d, m)
		}
	}
}

func checkSetBytes(b []byte) bool {
	hex1 := hex.EncodeToString(new(Int).SetBytes(b).Bytes())
	hex2 := hex.EncodeToString(b)

	for len(hex1) < len(hex2) {
		hex1 = "0" + hex1
	}

	for len(hex1) > len(hex2) {
		hex2 = "0" + hex2
	}

	return hex1 == hex2
}

func TestSetBytes(t *testing.T) {
	if err := quick.Check(checkSetBytes, nil); err != nil {
		t.Error(err)
	}
}

func checkBytes(b []byte) bool {
	b2 := new(Int).SetBytes(b).Bytes()
	return bytes.Equal(b, b2)
}

func TestBytes(t *testing.T) {
	if err := quick.Check(checkSetBytes, nil); err != nil {
		t.Error(err)
	}
}

func checkQuo(x, y []byte) bool {
	u := new(Int).SetBytes(x)
	v := new(Int).SetBytes(y)

	if v.Int64() == 0 {
		return true
	}

	r := new(Int)
	q, r := new(Int).QuoRem(u, v, r)

	if r.Cmp(v) >= 0 {
		return false
	}

	uprime := new(Int).Set(q)
	uprime.Mul(uprime, v)
	uprime.Add(uprime, r)

	return uprime.Cmp(u) == 0
}

var quoTests = []struct {
	x, y string
	q, r string
}{
	{
		"476217953993950760840509444250624797097991362735329973741718102894495832294430498335824897858659711275234906400899559094370964723884706254265559534144986498357",
		"9353930466774385905609975137998169297361893554149986716853295022578535724979483772383667534691121982974895531435241089241440253066816724367338287092081996",
		"50911",
		"1",
	},
	{
		"11510768301994997771168",
		"1328165573307167369775",
		"8",
		"885443715537658812968",
	},
}

func TestQuo(t *testing.T) {
	if err := quick.Check(checkQuo, nil); err != nil {
		t.Error(err)
	}

	for i, test := range quoTests {
		x, _ := new(Int).SetString(test.x, 10)
		y, _ := new(Int).SetString(test.y, 10)
		expectedQ, _ := new(Int).SetString(test.q, 10)
		expectedR, _ := new(Int).SetString(test.r, 10)

		r := new(Int)
		q, r := new(Int).QuoRem(x, y, r)

		if q.Cmp(expectedQ) != 0 || r.Cmp(expectedR) != 0 {
			t.Errorf("#%d got (%s, %s) want (%s, %s)", i, q, r, expectedQ, expectedR)
		}
	}
}

// func TestQuoStepD6(t *testing.T) {
// 	// See Knuth, Volume 2, section 4.3.1, exercise 21. This code exercises
// 	// a code path which only triggers 1 in 10^{-19} cases.

// 	u := &Int{false, nat{0, 0, 1 + 1<<(_W-1), _M ^ (1 << (_W - 1))}}
// 	v := &Int{false, nat{5, 2 + 1<<(_W-1), 1 << (_W - 1)}}

// 	r := new(Int)
// 	q, r := new(Int).QuoRem(u, v, r)
// 	const expectedQ64 = "18446744073709551613"
// 	const expectedR64 = "3138550867693340382088035895064302439801311770021610913807"
// 	const expectedQ32 = "4294967293"
// 	const expectedR32 = "39614081266355540837921718287"
// 	if q.String() != expectedQ64 && q.String() != expectedQ32 ||
// 		r.String() != expectedR64 && r.String() != expectedR32 {
// 		t.Errorf("got (%s, %s) want (%s, %s) or (%s, %s)", q, r, expectedQ64, expectedR64, expectedQ32, expectedR32)
// 	}
// }

var expTests = []struct {
	x, y, m string
	out     string
}{
	{"5", "-7", "", "1"},
	{"-5", "-7", "", "1"},
	{"5", "0", "", "1"},
	{"-5", "0", "", "1"},
	{"5", "1", "", "5"},
	{"-5", "1", "", "-5"},
	{"-2", "3", "2", "0"},
	{"5", "2", "", "25"},
	{"1", "65537", "2", "1"},
	{"0x8000000000000000", "2", "", "0x40000000000000000000000000000000"},
	{"0x8000000000000000", "2", "6719", "4944"},
	{"0x8000000000000000", "3", "6719", "5447"},
	{"0x8000000000000000", "1000", "6719", "1603"},
	{"0x8000000000000000", "1000000", "6719", "3199"},

	// FIXME: What to do about difference between gmp and math/big handling of negative exponents
	// {"0x8000000000000000", "-1000000", "6719", "1"},
	{
		"2938462938472983472983659726349017249287491026512746239764525612965293865296239471239874193284792387498274256129746192347",
		"298472983472983471903246121093472394872319615612417471234712061",
		"29834729834729834729347290846729561262544958723956495615629569234729836259263598127342374289365912465901365498236492183464",
		"23537740700184054162508175125554701713153216681790245129157191391322321508055833908509185839069455749219131480588829346291",
	},
}

func TestExp(t *testing.T) {
	for i, test := range expTests {
		x, ok1 := new(Int).SetString(test.x, 0)
		y, ok2 := new(Int).SetString(test.y, 0)
		out, ok3 := new(Int).SetString(test.out, 0)

		var ok4 bool
		var m *Int

		if len(test.m) == 0 {
			m, ok4 = nil, true
		} else {
			m, ok4 = new(Int).SetString(test.m, 0)
		}

		if !ok1 || !ok2 || !ok3 || !ok4 {
			t.Errorf("#%d: error in input", i)
			continue
		}

		z1 := new(Int).Exp(x, y, m)
		if z1.Cmp(out) != 0 {
			t.Errorf("#%d: got %s want %s", i, z1, out)
		}

		if m == nil {
			// the result should be the same as for m == 0;
			// specifically, there should be no div-zero panic
			m = NewInt(0)

			z2 := new(Int).Exp(x, y, m)
			if z2.Cmp(z1) != 0 {
				t.Errorf("#%d: got %s want %s", i, z1, z2)
			}
		}
	}
}

func checkGcd(aBytes, bBytes []byte) bool {
	x := new(Int)
	y := new(Int)
	a := new(Int).SetBytes(aBytes)
	b := new(Int).SetBytes(bBytes)

	d := new(Int).GCD(x, y, a, b)
	x.Mul(x, a)
	y.Mul(y, b)
	x.Add(x, y)

	return x.Cmp(d) == 0
}

var gcdTests = []struct {
	d, x, y, a, b string
}{
	// a <= 0 || b <= 0
	{"0", "0", "0", "0", "0"},
	{"0", "0", "0", "0", "7"},
	{"0", "0", "0", "11", "0"},
	{"0", "0", "0", "-77", "35"},
	{"0", "0", "0", "64515", "-24310"},
	{"0", "0", "0", "-64515", "-24310"},

	{"1", "-9", "47", "120", "23"},
	{"7", "1", "-2", "77", "35"},
	{"935", "-3", "8", "64515", "24310"},
	{"935000000000000000", "-3", "8", "64515000000000000000", "24310000000000000000"},
	{"1", "-221", "22059940471369027483332068679400581064239780177629666810348940098015901108344", "98920366548084643601728869055592650835572950932266967461790948584315647051443", "991"},

	// test early exit (after one Euclidean iteration) in binaryGCD
	{"1", "", "", "1", "98920366548084643601728869055592650835572950932266967461790948584315647051443"},
}

func testGcd(t *testing.T, d, x, y, a, b *Int) {
	var X *Int
	if x != nil {
		X = new(Int)
	}
	var Y *Int
	if y != nil {
		Y = new(Int)
	}
	D := new(Int).GCD(X, Y, a, b)
	if D.Cmp(d) != 0 {
		t.Errorf("GCD(%s, %s): got d = %s, want %s", a, b, D, d)
	}
	if x != nil && X.Cmp(x) != 0 {
		t.Errorf("GCD(%s, %s): got x = %s, want %s", a, b, X, x)
	}
	if y != nil && Y.Cmp(y) != 0 {
		t.Errorf("GCD(%s, %s): got y = %s, want %s", a, b, Y, y)
	}
}

func TestGcd(t *testing.T) {
	for _, test := range gcdTests {
		d, _ := new(Int).SetString(test.d, 0)
		x, _ := new(Int).SetString(test.x, 0)
		y, _ := new(Int).SetString(test.y, 0)
		a, _ := new(Int).SetString(test.a, 0)
		b, _ := new(Int).SetString(test.b, 0)

		testGcd(t, d, nil, nil, a, b)
		testGcd(t, d, x, nil, a, b)
		testGcd(t, d, nil, y, a, b)
		testGcd(t, d, x, y, a, b)
	}

	quick.Check(checkGcd, nil)
}

var stringTests = []struct {
	in   string
	out  string
	base int
	val  int64
	ok   bool
}{
	{in: "", ok: false},
	{in: "a", ok: false},
	{in: "z", ok: false},
	{in: "+", ok: false},
	{in: "-", ok: false},
	{in: "0b", ok: false},
	{in: "0x", ok: false},
	{in: "2", base: 2, ok: false},
	{in: "0b2", base: 0, ok: false},
	{in: "08", ok: false},
	{in: "8", base: 8, ok: false},
	{in: "0xg", base: 0, ok: false},
	{in: "g", base: 16, ok: false},
	{"0", "0", 0, 0, true},
	{"0", "0", 10, 0, true},
	{"0", "0", 16, 0, true},
	{"+0", "0", 0, 0, true},
	{"-0", "0", 0, 0, true},
	{"10", "10", 0, 10, true},
	{"10", "10", 10, 10, true},
	{"10", "10", 16, 16, true},
	{"-10", "-10", 16, -16, true},
	{"+10", "10", 16, 16, true},
	{"0x10", "16", 0, 16, true},
	{in: "0x10", base: 16, ok: false},
	{"-0x10", "-16", 0, -16, true},
	{"+0x10", "16", 0, 16, true},
	{"00", "0", 0, 0, true},
	{"0", "0", 8, 0, true},
	{"07", "7", 0, 7, true},
	{"7", "7", 8, 7, true},
	{"023", "19", 0, 19, true},
	{"23", "23", 8, 19, true},
	{"cafebabe", "cafebabe", 16, 0xcafebabe, true},
	{"0b0", "0", 0, 0, true},
	{"-111", "-111", 2, -7, true},
	{"-0b111", "-7", 0, -7, true},
	{"0b1001010111", "599", 0, 0x257, true},
	{"1001010111", "1001010111", 2, 0x257, true},
}

func format(base int) string {
	switch base {
	case 2:
		return "%b"
	case 8:
		return "%o"
	case 16:
		return "%x"
	}
	return "%d"
}

func TestSetString(t *testing.T) {
	tmp := new(Int)
	for i, test := range stringTests {
		// initialize to a non-zero value so that issues with parsing
		// 0 are detected
		tmp.SetInt64(1234567890)
		n1, ok1 := new(Int).SetString(test.in, test.base)
		n2, ok2 := tmp.SetString(test.in, test.base)
		expected := NewInt(test.val)
		if ok1 != test.ok || ok2 != test.ok {
			t.Errorf("#%d (input '%s') ok incorrect (should be %t)", i, test.in, test.ok)
			continue
		}
		if !ok1 {
			if n1 != nil {
				t.Errorf("#%d (input '%s') n1 != nil", i, test.in)
			}
			continue
		}
		if !ok2 {
			if n2 != nil {
				t.Errorf("#%d (input '%s') n2 != nil", i, test.in)
			}
			continue
		}

		if n1.Cmp(expected) != 0 {
			t.Errorf("#%d (input '%s') got: %s want: %d", i, test.in, n1, test.val)
		}
		if n2.Cmp(expected) != 0 {
			t.Errorf("#%d (input '%s') got: %s want: %d", i, test.in, n2, test.val)
		}
	}
}

var primes = []string{
	"2",
	"3",
	"5",
	"7",
	"11",

	"13756265695458089029",
	"13496181268022124907",
	"10953742525620032441",
	"17908251027575790097",

	// http://code.google.com/p/go/issues/detail?id=638
	"18699199384836356663",

	"98920366548084643601728869055592650835572950932266967461790948584315647051443",
	"94560208308847015747498523884063394671606671904944666360068158221458669711639",

	// http://primes.utm.edu/lists/small/small3.html
	"449417999055441493994709297093108513015373787049558499205492347871729927573118262811508386655998299074566974373711472560655026288668094291699357843464363003144674940345912431129144354948751003607115263071543163",
	"230975859993204150666423538988557839555560243929065415434980904258310530753006723857139742334640122533598517597674807096648905501653461687601339782814316124971547968912893214002992086353183070342498989426570593",
	"5521712099665906221540423207019333379125265462121169655563495403888449493493629943498064604536961775110765377745550377067893607246020694972959780839151452457728855382113555867743022746090187341871655890805971735385789993",
	"203956878356401977405765866929034577280193993314348263094772646453283062722701277632936616063144088173312372882677123879538709400158306567338328279154499698366071906766440037074217117805690872792848149112022286332144876183376326512083574821647933992961249917319836219304274280243803104015000563790123",
}

var composites = []string{
	"21284175091214687912771199898307297748211672914763848041968395774954376176754",
	"6084766654921918907427900243509372380954290099172559290432744450051395395951",
	"84594350493221918389213352992032324280367711247940675652888030554255915464401",
	"82793403787388584738507275144194252681",
}

func TestProbablyPrime(t *testing.T) {
	nreps := 20
	if testing.Short() {
		nreps = 1
	}
	for i, s := range primes {
		p, _ := new(Int).SetString(s, 10)
		if !p.ProbablyPrime(nreps) {
			t.Errorf("#%d prime found to be non-prime (%s)", i, s)
		}
	}

	for i, s := range composites {
		c, _ := new(Int).SetString(s, 10)
		if c.ProbablyPrime(nreps) {
			t.Errorf("#%d composite found to be prime (%s)", i, s)
		}
		if testing.Short() {
			break
		}
	}
}

type intShiftTest struct {
	in    string
	shift uint
	out   string
}

var rshTests = []intShiftTest{
	{"0", 0, "0"},
	{"-0", 0, "0"},
	{"0", 1, "0"},
	{"0", 2, "0"},
	{"1", 0, "1"},
	{"1", 1, "0"},
	{"1", 2, "0"},
	{"2", 0, "2"},
	{"2", 1, "1"},
	{"-1", 0, "-1"},
	{"-1", 1, "-1"},
	{"-1", 10, "-1"},
	{"-100", 2, "-25"},
	{"-100", 3, "-13"},
	{"-100", 100, "-1"},
	{"4294967296", 0, "4294967296"},
	{"4294967296", 1, "2147483648"},
	{"4294967296", 2, "1073741824"},
	{"18446744073709551616", 0, "18446744073709551616"},
	{"18446744073709551616", 1, "9223372036854775808"},
	{"18446744073709551616", 2, "4611686018427387904"},
	{"18446744073709551616", 64, "1"},
	{"340282366920938463463374607431768211456", 64, "18446744073709551616"},
	{"340282366920938463463374607431768211456", 128, "1"},
}

func TestRsh(t *testing.T) {
	for i, test := range rshTests {
		in, _ := new(Int).SetString(test.in, 10)
		expected, _ := new(Int).SetString(test.out, 10)
		out := new(Int).Rsh(in, test.shift)

		if out.Cmp(expected) != 0 {
			t.Errorf("#%d: got %s want %s", i, out, expected)
		}
	}
}

func TestRshSelf(t *testing.T) {
	for i, test := range rshTests {
		z, _ := new(Int).SetString(test.in, 10)
		expected, _ := new(Int).SetString(test.out, 10)
		z.Rsh(z, test.shift)

		if z.Cmp(expected) != 0 {
			t.Errorf("#%d: got %s want %s", i, z, expected)
		}
	}
}

var lshTests = []intShiftTest{
	{"0", 0, "0"},
	{"0", 1, "0"},
	{"0", 2, "0"},
	{"1", 0, "1"},
	{"1", 1, "2"},
	{"1", 2, "4"},
	{"2", 0, "2"},
	{"2", 1, "4"},
	{"2", 2, "8"},
	{"-87", 1, "-174"},
	{"4294967296", 0, "4294967296"},
	{"4294967296", 1, "8589934592"},
	{"4294967296", 2, "17179869184"},
	{"18446744073709551616", 0, "18446744073709551616"},
	{"9223372036854775808", 1, "18446744073709551616"},
	{"4611686018427387904", 2, "18446744073709551616"},
	{"1", 64, "18446744073709551616"},
	{"18446744073709551616", 64, "340282366920938463463374607431768211456"},
	{"1", 128, "340282366920938463463374607431768211456"},
}

func TestLsh(t *testing.T) {
	for i, test := range lshTests {
		in, _ := new(Int).SetString(test.in, 10)
		expected, _ := new(Int).SetString(test.out, 10)
		out := new(Int).Lsh(in, test.shift)

		if out.Cmp(expected) != 0 {
			t.Errorf("#%d: got %s want %s", i, out, expected)
		}
	}
}

func TestLshSelf(t *testing.T) {
	for i, test := range lshTests {
		z, _ := new(Int).SetString(test.in, 10)
		expected, _ := new(Int).SetString(test.out, 10)
		z.Lsh(z, test.shift)

		if z.Cmp(expected) != 0 {
			t.Errorf("#%d: got %s want %s", i, z, expected)
		}
	}
}

func TestLshRsh(t *testing.T) {
	for i, test := range rshTests {
		in, _ := new(Int).SetString(test.in, 10)
		out := new(Int).Lsh(in, test.shift)
		out = out.Rsh(out, test.shift)

		if in.Cmp(out) != 0 {
			t.Errorf("#%d: got %s want %s", i, out, in)
		}
	}
	for i, test := range lshTests {
		in, _ := new(Int).SetString(test.in, 10)
		out := new(Int).Lsh(in, test.shift)
		out.Rsh(out, test.shift)

		if in.Cmp(out) != 0 {
			t.Errorf("#%d: got %s want %s", i, out, in)
		}
	}
}

var int64Tests = []int64{
	0,
	1,
	-1,
	4294967295,
	-4294967295,
	4294967296,
	-4294967296,
	9223372036854775807,
	-9223372036854775807,
	-9223372036854775808,
}

func TestInt64(t *testing.T) {
	for i, testVal := range int64Tests {
		in := NewInt(testVal)
		out := in.Int64()

		if out != testVal {
			t.Errorf("#%d got %d want %d", i, out, testVal)
		}
	}
}

var uint64Tests = []uint64{
	0,
	1,
	4294967295,
	4294967296,
	8589934591,
	8589934592,
	9223372036854775807,
	9223372036854775808,
	18446744073709551615, // 1<<64 - 1
}

func TestUint64(t *testing.T) {
	in := new(Int)
	for i, testVal := range uint64Tests {
		in.SetUint64(testVal)
		out := in.Uint64()

		if out != testVal {
			t.Errorf("#%d got %d want %d", i, out, testVal)
		}

		str := fmt.Sprint(testVal)
		strOut := in.String()
		if strOut != str {
			t.Errorf("#%d.String got %s want %s", i, strOut, str)
		}
	}
}

var modInverseTests = []struct {
	element string
	prime   string
}{
	{"1", "7"},
	{"1", "13"},
	{"239487239847", "2410312426921032588552076022197566074856950548502459942654116941958108831682612228890093858261341614673227141477904012196503648957050582631942730706805009223062734745341073406696246014589361659774041027169249453200378729434170325843778659198143763193776859869524088940195577346119843545301547043747207749969763750084308926339295559968882457872412993810129130294592999947926365264059284647209730384947211681434464714438488520940127459844288859336526896320919633919"},
}

func TestModInverse(t *testing.T) {
	var element, prime Int
	one := NewInt(1)
	for i, test := range modInverseTests {
		(&element).SetString(test.element, 10)
		(&prime).SetString(test.prime, 10)
		inverse := new(Int).ModInverse(&element, &prime)
		inverse.Mul(inverse, &element)
		inverse.Mod(inverse, &prime)
		if inverse.Cmp(one) != 0 {
			t.Errorf("#%d: failed (e·e^(-1)=%s)", i, inverse)
		}
	}
}
