PACKAGE

package gogmp
    import "github.com/i-neda/gogmp"


FUNCTIONS

func CmpFloat(x, y *Float) int

func CmpFloatDouble(x *Float, y float64) int

func CmpFloatSint(x *Float, y int) int

func CmpFloatUint(x *Float, y uint) int

func CmpInt(x, y *Int) int
    CmpInt compares x and y. The result is

	-1 if x <  y
	 0 if x == y
	+1 if x >  y

func CmpRat(x, y *Rat) int

func CmpRatSint(q *Rat, x int, y uint) int

func CmpRatUint(q *Rat, x, y uint) int

func DivModInt(q, r, x, y *Int)
    DivModInt sets q = x / y and r = x % y.

func EqFloat(x, y *Float, n uint) int
    Return non-zero if the first n bits of x and y are equal,

	zero otherwise.  I.e., test if x and y are approximately equal.

func EqRat(x, y *Rat) bool

func GcdInt(d, x, y, a, b *Int)
    GcdInt sets d to the greatest common divisor of a and b, which must be
    positive numbers. If x and y are not nil, GcdInt sets x and y such that
    d = a*x + b*y. If either a or b is not positive, GcdInt sets d = x = y =
    0.

func GetDefaultPrec() uint

func SetDefaultPrec(prec uint)

func SwapFloat(x, y *Float)

func SwapRat(x, y *Rat)


TYPES

type Float struct {
    // contains filtered or unexported fields
}

func NewFloat(x float64) *Float
    NewInt returns a new Int initialized to x.

func NewFloat2(x float64, prec uint) *Float
    NewInt returns a new Int initialized to x, with precision prec.

func (f *Float) Abs(x *Float) *Float
    Abs sets z to the absolute value of x and returns z.

func (f *Float) Add(x, y *Float) *Float
    Add sets f = x + y and returns f.

func (f *Float) AddUint(x *Float, y uint) *Float

func (f *Float) Ceil(x *Float) *Float
    Sets f = Ceil(x) and returns f.

func (f *Float) Clear()

func (f *Float) Div(x, y *Float) *Float
    Div sets f = x / y and returns f.

func (f *Float) Div2Exp(x *Float, s uint) *Float
    Div2Exp sets z = x / 2^s and returns z.

func (f *Float) DivUint(x *Float, y uint) *Float

func (f *Float) Double() float64

func (f *Float) Double2Exp() (d float64, exp int)
    * Convert f to a `double', truncating if necessary (ie. rounding *
    towards zero), and with an exponent returned separately.

func (f *Float) Floor(x *Float) *Float
    Sets f = Floor(x) and returns f.

func (f *Float) GetPrec() uint

func (f *Float) IsInteger() bool

func (f *Float) Mul(x, y *Float) *Float
    Mul sets f = x * y and returns f.

func (f *Float) Mul2Exp(x *Float, s uint) *Float
    Mul2Exp sets z = x * 2^s and returns z.

func (f *Float) MulUint(x *Float, y uint) *Float

func (f *Float) Neg(x *Float) *Float
    Neg sets z = -x and returns z.

func (f *Float) PowUint(x *Float, y uint) *Float
    PowUint sets f = x^y and returns f

func (f *Float) RelDiff(x, y *Float) *Float

func (f *Float) Set(x *Float) *Float
    Set sets f = x and returns f.

func (f *Float) SetDouble(x float64) *Float
    SetDouble sets f = x and returns f.

func (f *Float) SetPrec(prec uint)

func (f *Float) SetPrecRaw(prec uint)

func (f *Float) SetSint(x int) *Float
    SetInt sets f = x and returns f.

func (f *Float) SetString(s string, base int) error
    SetString interprets s as a number in the given base and sets f to that
    value. The base must be in the range [2,36]. SetString returns an error
    if s cannot be parsed or the base is invalid.

func (f *Float) Sgn() int
    Return +1 if f > 0, 0 if f = 0, and -1 if f < 0.

func (f *Float) Sint() int

func (f *Float) Sqrt(x *Float) *Float
    Sqrt sets f = Sqrt(x) and returns f.

func (f *Float) SqrtUint(x uint) *Float
    Sqrt sets f = Sqrt(x) and returns f.

func (f *Float) String() string
    String returns the decimal representation of z.

func (f *Float) StringBase(base int, ndigits uint) (string, error)

func (f *Float) Sub(x, y *Float) *Float
    Sub sets f = x - y and returns f.

func (f *Float) SubUint(x *Float, y uint) *Float

func (f *Float) Trunc(x *Float) *Float
    Sets f = Trunc(x) (=round towards zero) and returns f.

func (f *Float) UintDiv(x uint, y *Float) *Float

type Int struct {
    // contains filtered or unexported fields
}
    An Int represents a signed multi-precision integer. The zero value for
    an Int represents the value 0.

func NewInt(x int) *Int
    NewInt returns a new Int initialized to x.

func (z *Int) Abs(x *Int) *Int
    Abs sets z to the absolute value of x and returns z.

func (z *Int) Add(x, y *Int) *Int
    Add sets z = x + y and returns z.

func (z *Int) Bytes() []byte
    Bytes returns z's representation as a big-endian byte array.

func (z *Int) Clear()

func (z *Int) Div(x, y *Int) *Int
    Div sets z = x / y, rounding toward zero, and returns z.

func (z *Int) Exp(x, y, m *Int) *Int
    Exp sets z = x^y % m and returns z. If m == nil, Exp sets z = x^y.

func (z *Int) Int64() int64
    Provided for compatibility with big package

func (z *Int) Len() int
    Len returns the length of z in bits. 0 is considered to have length 1.

func (z *Int) Lsh(x *Int, s uint) *Int
    Lsh sets z = x << s and returns z.

func (z *Int) Mod(x, y *Int) *Int
    Mod sets z = x % y and returns z. Like the result of the Go % operator,
    z has the same sign as x.

func (z *Int) Mul(x, y *Int) *Int
    Mul sets z = x * y and returns z.

func (z *Int) Neg(x *Int) *Int
    Neg sets z = -x and returns z.

func (z *Int) ProbablyPrime(n int) bool
    ProbablyPrime performs n Miller-Rabin tests to check whether z is prime.
    If it returns true, z is prime with probability 1 - 1/4^n. If it returns
    false, z is not prime.

func (z *Int) Rsh(x *Int, s uint) *Int
    Rsh sets z = x >> s and returns z.

func (z *Int) Set(x *Int) *Int
    Set sets z = x and returns z.

func (z *Int) SetBytes(b []byte) *Int
    SetBytes interprets b as the bytes of a big-endian integer and sets z to
    that value.

func (z *Int) SetInt64(x int64) *Int
    Provided for compatibility with big package

func (z *Int) SetSint(x int) *Int
    SetSint sets z = x and returns z.

func (z *Int) SetString(s string, base int) error
    SetString interprets s as a number in the given base and sets z to that
    value. The base must be in the range [2,36]. SetString returns an error
    if s cannot be parsed or the base is invalid.

func (z *Int) Sint() int

func (z *Int) String() string
    String returns the decimal representation of z.

func (z *Int) StringBase(base int) (string, error)

func (z *Int) Sub(x, y *Int) *Int
    Sub sets z = x - y and returns z.

type Rat struct {
    // contains filtered or unexported fields
}

func NewRat(x int, y uint) *Rat
    NewRat returns a new Rat initialized to x/y.

func (q *Rat) Abs(x *Rat) *Rat

func (q *Rat) Add(x, y *Rat) *Rat
    Add sets q = x + y and returns f.

func (q *Rat) Clear()

func (q *Rat) Den(n *Int) *Int

func (q *Rat) Div(x, y *Rat) *Rat

func (q *Rat) Div2Exp(x *Rat, s uint) *Rat
    Div2Exp sets z = x / 2^s and returns z.

func (q *Rat) Double() float64

func (q *Rat) Inv(x *Rat) *Rat

func (q *Rat) Mul(x, y *Rat) *Rat

func (q *Rat) Mul2Exp(x *Rat, s uint) *Rat
    Mul2Exp sets z = x * 2^s and returns z.

func (q *Rat) Num(n *Int) *Int

func (q *Rat) Set(x *Rat) *Rat
    Set sets z = x and returns z.

func (q *Rat) SetDen(n *Int) *Rat

func (q *Rat) SetDouble(x float64) *Rat
    SetDouble sets f = x and returns q.

func (q *Rat) SetFloat(x *Float) *Rat
    SetFloat sets f = x and returns f.

func (q *Rat) SetInt(x *Int) *Rat
    SetInt sets q to x and returns q.

func (q *Rat) SetNum(n *Int) *Rat

func (q *Rat) SetSint(x int, y uint) *Rat
    SetSint sets q to x/y and returns q.

func (q *Rat) SetString(s string, base int) error
    SetString interprets s as a number in the given base and sets z to that
    value. The base must be in the range [2,36]. SetString returns an error
    if s cannot be parsed or the base is invalid.

func (q *Rat) SetUint(x, y uint) *Rat
    SetUint sets q to x/y and returns q.

func (q *Rat) Sgn() int

func (q *Rat) String() string

func (q *Rat) StringBase(base int) (string, error)
    String returns the decimal representation of z.

func (q *Rat) Sub(x, y *Rat) *Rat


