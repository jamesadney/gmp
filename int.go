// Copyright 2010 Utkan Güngördü.
// Based on $(GOROOT)/misc/cgo/gmp/gmp.go
// Released under the BSD-style license that can
// be found in Go's LICENSE file.

// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gmp

/*
#cgo LDFLAGS: -lgmp
#include <gmp.h>
#include <stdlib.h>

// gmp 5.0.0+ changed the type of the 3rd argument to mp_bitcnt_t,
// so, to support older versions, we wrap these two functions.
void _mpz_mul_2exp(mpz_ptr a, mpz_ptr b, unsigned long n) {
	mpz_mul_2exp(a, b, n);
}
void _mpz_div_2exp(mpz_ptr a, mpz_ptr b, unsigned long n) {
	mpz_div_2exp(a, b, n);
}

// since mpz_sgn is a macro we have to wrap it in a function.
int _mpz_sgn(mpz_ptr z) {
	return mpz_sgn(z);
}
*/
import "C"

import (
	"os"
	"unsafe"
)

var (
	intZero = NewInt(0)
	intOne  = NewInt(1)
)

// An Int represents a signed multi-precision integer.
// The zero value for an Int represents the value 0.
type Int struct {
	ptr  C.mpz_ptr // pointer to underlying mpz so Int can be a reference
	t    C.mpz_t   // mpz_t that needs to be initialized
	init bool
}

// Int promises that the zero value is a 0, but in gmp
// the zero value is a crash.  To bridge the gap, the
// init bool says whether this is a valid gmp value.
// doinit initializes z.t if it needs it.  This is not inherent
// to FFI, just a mismatch between Go's convention of
// making zero values useful and gmp's decision not to.
//
// Make z a reference to another mpz_t (not z.t) by directly assigning
// z.ptr and leaving z.t uninitialized. See Rat.Denom() for an example.
func (z *Int) doinit() {
	if z.init {
		return
	}
	z.init = true
	z.ptr = &z.t[0]
	C.mpz_init(z.ptr)
}

// Free the space occupied by the underlying gmp object.
func (z *Int) Clear() {
	if z.init {
		C.mpz_clear(z.ptr)
	}
	z.init = false
}

/*
 * assigning and converting
 */

// NewInt returns a new Int initialized to x.
func NewInt(x int64) *Int { return new(Int).SetInt64(x) }

// Set sets z = x and returns z.
func (z *Int) Set(x *Int) *Int {
	z.doinit()
	C.mpz_set(z.ptr, x.ptr)
	return z
}

// Bytes returns z's representation as a big-endian byte array.
func (z *Int) Bytes() []byte {
	b := make([]byte, (z.Len()+7)/8)
	n := C.size_t(len(b))
	C.mpz_export(unsafe.Pointer(&b[0]), &n, 1, 1, 1, 0, z.ptr)
	return b[0:n]
}

// SetBytes interprets b as the bytes of a big-endian integer
// and sets z to that value.
func (z *Int) SetBytes(b []byte) *Int {
	z.doinit()
	if len(b) == 0 {
		z.SetInt64(0)
	} else {
		C.mpz_import(z.ptr, C.size_t(len(b)), 1, 1, 1, 0,
			unsafe.Pointer(&b[0]))
	}
	return z
}

// BitLen returns the length of the absolute value of z in bits.
// The bit length of 0 is 0.
func (x *Int) BitLen() int {
	x.doinit()
	if x.Cmp(intZero) == 0 {
		return 0
	}
	return x.Len()
}

// Len returns the length of z in bits.  0 is considered to have length 1.
func (z *Int) Len() int {
	z.doinit()
	return int(C.mpz_sizeinbase(z.ptr, 2))
}

// Int64 returns the int64 representation of x. If x cannot be represented
// in an int64, the result is undefined.
func (z *Int) Int64() int64 {
	if !z.init {
		return 0
	}
	return int64(C.mpz_get_si(z.ptr))
}

// SetInt64 sets z = x and returns z.
func (z *Int) SetInt64(x int64) *Int {
	z.doinit()
	// TODO(rsc): more work on 32-bit platforms
	C.mpz_set_si(z.ptr, C.long(x))
	return z
}

// Uint64 returns the uint64 representation of x. If x cannot be
// represented in an uint64, the result is undefined.
func (z *Int) Uint64() uint64 {
	if !z.init {
		return 0
	}
	return uint64(C.mpz_get_ui(z.ptr))
}

// SetUint64 sets z to x and returns z.
func (z *Int) SetUint64(x uint64) *Int {
	z.doinit()
	C.mpz_set_ui(z.ptr, C.ulong(x))
	return z
}

// String returns the decimal representation of z.
func (z *Int) String() string {
	s, _ := z.StringBase(10)
	return s
}

func (z *Int) StringBase(base int) (string, error) {
	if z == nil {
		return "nil", nil
	}
	if base < 2 || base > 36 {
		return "", os.ErrInvalid
	}
	z.doinit()
	p := C.mpz_get_str(nil, C.int(base), z.ptr)
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))
	return s, nil
}

// SetString sets z to the value of s, interpreted in the given base,
// and returns z and a boolean indicating success. If SetString fails, the
// value of z is undefined but the returned value is nil.

// The base argument must be 0 or a value from 2 through 36. If the base is 0,
// the string prefix determines the actual conversion base. A prefix of “0x” or
// “0X” selects base 16; the “0” prefix selects base 8, and a “0b” or “0B”
// prefix selects base 2. Otherwise the selected base is 10.
func (z *Int) SetString(s string, base int) (*Int, bool) {
	z.doinit()
	if base < 0 || base == 1 || base > 36 {
		return nil, false
	}

	// no need to call mpz_set_str here.
	if len(s) == 0 {
		return nil, false
	}

	// positive signs should be ignored
	if s[0] == '+' {
		s = s[1:]
	}

	// attempting to set "0x" and "0b" should return nil like math/big
	if len(s) == 2 {
		switch s {
		case "0x", "0X", "0b", "0B":
			return nil, false
		}
	}

	p := C.CString(s)
	defer C.free(unsafe.Pointer(p))
	if C.mpz_set_str(z.ptr, p, C.int(base)) < 0 {
		return nil, false
	}
	return z, true
}

/*
 * arithmetic
 */

// Add sets z = x + y and returns z.
func (z *Int) Add(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_add(z.ptr, x.ptr, y.ptr)
	return z
}

// Sub sets z = x - y and returns z.
func (z *Int) Sub(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_sub(z.ptr, x.ptr, y.ptr)
	return z
}

// Mul sets z = x * y and returns z.
func (z *Int) Mul(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_mul(z.ptr, x.ptr, y.ptr)
	return z
}

// mulRange computes the product of all the unsigned integers in the
// range [a, b] inclusively. If a > b (empty range), the result is 1.
func (z *Int) mulRange(a, b uint64) *Int {
	switch {
	case a == 0:
		// cut long ranges short (optimization)
		return z.SetUint64(0)
	case a > b:
		return z.SetUint64(1)
	case a == b:
		return z.SetUint64(a)
	case a+1 == b:
		A, B := new(Int).SetUint64(a), new(Int).SetUint64(b)
		z.Mul(A, B)
		A.Clear()
		B.Clear()
		return z
	}
	m := (a + b) / 2
	temp_a := new(Int).mulRange(a, m)
	temp_b := new(Int).mulRange(m+1, b)

	z.Mul(temp_a, temp_b)

	temp_a.Clear()
	temp_b.Clear()
	return z
}

// MulRange sets z to the product of all integers
// in the range [a, b] inclusively and returns z.
// If a > b (empty range), the result is 1.
func (z *Int) MulRange(a, b int64) *Int {
	switch {
	case a > b:
		return z.SetInt64(1) // empty range
	case a <= 0 && b >= 0:
		return z.SetInt64(0) // range includes 0
	}
	// a <= b && (b < 0 || a > 0)

	neg := false
	if a < 0 {
		neg = (b-a)&1 == 0
		a, b = -b, -a
	}

	z = z.mulRange(uint64(a), uint64(b))
	if neg {
		negativeOne := NewInt(-1)
		z.Mul(z, negativeOne)
		negativeOne.Clear()
	}
	return z
}

// Quo sets z to the quotient x/y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Quo implements truncated division (like Go).
func (z *Int) Quo(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_tdiv_q(z.ptr, x.ptr, y.ptr)
	return z
}

// Rem sets z to the remainder x%y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Rem implements truncated modulus (like Go); see QuoRem for more details.
func (z *Int) Rem(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_tdiv_r(z.ptr, x.ptr, y.ptr)
	return z
}

// QuoRem sets z to the quotient x/y and r to the remainder x%y
// and returns the pair (z, r) for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
//
// QuoRem implements T-division and modulus (like Go):
//
//	q = x/y      with the result truncated to zero
//	r = x - y*q
//
// (See Daan Leijen, ``Division and Modulus for Computer Scientists''.)
// See DivMod for Euclidean division and modulus (unlike Go).
//
func (z *Int) QuoRem(x, y, r *Int) (*Int, *Int) {
	x.doinit()
	y.doinit()
	r.doinit()
	z.doinit()
	C.mpz_tdiv_qr(z.ptr, r.ptr, x.ptr, y.ptr)
	return z, r
}

// Div sets z to the quotient x/y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Div implements Euclidean division (unlike Go); see DivMod for more details.
func (z *Int) Div(x, y *Int) *Int {
	y_neg := y.Sign() == -1 // z may be an alias for y
	var r Int
	z.QuoRem(x, y, &r)
	if r.Sign() == -1 {
		if y_neg {
			z.Add(z, intOne)
		} else {
			z.Sub(z, intOne)
		}
	}
	return z
}

// Mod sets z to the modulus x%y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Mod implements Euclidean modulus (unlike Go); see DivMod for more details.
func (z *Int) Mod(x, y *Int) *Int {
	y0 := y // save y
	if z == y {
		y0 = new(Int).Set(y)
		defer y0.Clear()
	}
	var q Int
	q.QuoRem(x, y, z)
	if z.Sign() == -1 {
		if y0.Sign() == -1 {
			z.Sub(z, y0)
		} else {
			z.Add(z, y0)
		}
	}
	return z
}

// DivMod sets z to the quotient x div y and m to the modulus x mod y
// and returns the pair (z, m) for y != 0.
// If y == 0, a division-by-zero run-time panic occurs.
//
// DivMod implements Euclidean division and modulus (unlike Go):
//
//	q = x div y  such that
//	m = x - y*q  with 0 <= m < |q|
//
// (See Raymond T. Boute, ``The Euclidean definition of the functions
// div and mod''. ACM Transactions on Programming Languages and
// Systems (TOPLAS), 14(2):127-144, New York, NY, USA, 4/1992.
// ACM press.)
// See QuoRem for T-division and modulus (like Go).
//
func (z *Int) DivMod(x, y, m *Int) (*Int, *Int) {
	y0 := y // save y
	if z == y {
		y0 = new(Int).Set(y)
		defer y0.Clear()
	}
	z.QuoRem(x, y, m)
	if m.Sign() == -1 {
		if y0.Sign() == -1 {
			z.Add(z, intOne)
			m.Sub(m, y0)
		} else {
			z.Sub(z, intOne)
			m.Add(m, y0)
		}
	}
	return z, m
}

// Exp sets z = x^y % m and returns z. If m != nil, negative exponents are
// allowed if x^-1 mod m exists. If the inverse doesn't exist then a
// division-by-zero run-time panic occurs.
//
// If m == nil, Exp sets z = x^y for positive y and 1 for negative y.
func (z *Int) Exp(x, y, m *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	if m == nil || m.Cmp(intZero) == 0 {
		if y.Sign() == -1 {
			z := NewInt(1)
			return z
		}
		C.mpz_pow_ui(z.ptr, x.ptr, C.mpz_get_ui(y.ptr))
	} else {
		m.doinit()
		C.mpz_powm(z.ptr, x.ptr, y.ptr, m.ptr)
	}
	return z
}

// Sqrt sets z = floor(sqrt(x)) and returns z.
func (z *Int) Sqrt(x *Int) *Int {
	z.doinit()
	x.doinit()
	C.mpz_sqrt(z.ptr, x.ptr)
	return z
}

// Neg sets z = -x and returns z.
func (z *Int) Neg(x *Int) *Int {
	x.doinit()
	z.doinit()
	C.mpz_neg(z.ptr, x.ptr)
	return z
}

// Abs sets z to the absolute value of x and returns z.
func (z *Int) Abs(x *Int) *Int {
	x.doinit()
	z.doinit()
	C.mpz_abs(z.ptr, x.ptr)
	return z
}

/*
 * logic and bit fiddling
 */

// Lsh sets z = x << s and returns z.
func (z *Int) Lsh(x *Int, s uint) *Int {
	x.doinit()
	z.doinit()
	C._mpz_mul_2exp(z.ptr, x.ptr, C.ulong(s))
	return z
}

// Rsh sets z = x >> s and returns z.
func (z *Int) Rsh(x *Int, s uint) *Int {
	x.doinit()
	z.doinit()
	C._mpz_div_2exp(z.ptr, x.ptr, C.ulong(s))
	return z
}

// And sets z = x & y and returns z.
func (z *Int) And(x, y *Int) *Int {
	z.doinit()
	x.doinit()
	y.doinit()
	C.mpz_and(z.ptr, x.ptr, y.ptr)
	return z
}

// AndNot sets z = x &^ y and returns z.
func (z *Int) AndNot(x, y *Int) *Int {
	z.doinit()
	x.doinit()
	y.doinit()

	temp := new(Int).Not(y)
	defer temp.Clear()
	return z.And(x, temp)
}

// Or sets z = x | y and returns z.
func (z *Int) Or(x, y *Int) *Int {
	z.doinit()
	x.doinit()
	y.doinit()
	C.mpz_ior(z.ptr, x.ptr, y.ptr)
	return z
}

// Xor sets z = x ^ y and returns z.
func (z *Int) Xor(x, y *Int) *Int {
	z.doinit()
	x.doinit()
	y.doinit()
	C.mpz_xor(z.ptr, x.ptr, y.ptr)
	return z
}

// Not sets z = ^x and returns z.
func (z *Int) Not(x *Int) *Int {
	z.doinit()
	x.doinit()
	C.mpz_com(z.ptr, x.ptr)
	return z
}

// Bit returns the value of the i'th bit of x. That is, it
// returns (x>>i)&1. The bit index i must be >= 0.
func (x *Int) Bit(i int) uint {
	x.doinit()
	return uint(C.mpz_tstbit(x.ptr, C.mp_bitcnt_t(i)))
}

// SetBit sets z to x, with x's i'th bit set to b (0 or 1).
// That is, if bit is 1 SetBit sets z = x | (1 << i);
// if bit is 0 it sets z = x &^ (1 << i). If bit is not 0 or 1,
// SetBit will panic.
func (z *Int) SetBit(x *Int, i int, b uint) *Int {
	z.doinit()

	if i < 0 {
		panic("negative bit index")
	}

	z.Set(x)
	switch b {
	case 0:
		C.mpz_clrbit(z.ptr, C.mp_bitcnt_t(i))
	case 1:
		C.mpz_setbit(z.ptr, C.mp_bitcnt_t(i))
	default:
		panic("set bit is not 0 or 1")
	}
	return z
}

/*
 * number theory
 */

// ModInverse sets z to the multiplicative inverse of g in the group ℤ/pℤ
// (where p is a prime) and returns z.
func (z *Int) ModInverse(g, p *Int) *Int {
	g.doinit()
	p.doinit()
	z.doinit()
	C.mpz_invert(z.ptr, g.ptr, p.ptr)
	return z
}

// GCD sets z to the greatest common divisor of a and b, which must be positive
// numbers, and returns z. If x and y are not nil, GCD sets x and y such that
// z = a*x + b*y. If either a or b is not positive, GCD sets z = x = y = 0.
func (z *Int) GCD(x, y, a, b *Int) *Int {
	z.doinit()

	// Compatibility with math/big
	if a.Cmp(intZero) <= 0 || b.Cmp(intZero) <= 0 {
		z.Set(intZero)
		return z
	}

	// allow for nil x and y
	var x_ptr C.mpz_ptr = nil
	if x != nil {
		x.doinit()
		x_ptr = x.ptr
	}
	var y_ptr C.mpz_ptr = nil
	if y != nil {
		y.doinit()
		y_ptr = y.ptr
	}

	a.doinit()
	b.doinit()
	C.mpz_gcdext(z.ptr, x_ptr, y_ptr, a.ptr, b.ptr)
	return z
}

// ProbablyPrime performs n Miller-Rabin tests to check whether z is prime.
// If it returns true, z is prime with probability 1 - 1/4^n.
// If it returns false, z is not prime.
func (z *Int) ProbablyPrime(n int) bool {
	z.doinit()
	return int(C.mpz_probab_prime_p(z.ptr, C.int(n))) > 0
}

/*
 * comparisons
 */

// Sign returns:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0
//
func (z *Int) Sign() int {
	z.doinit()
	return int(C._mpz_sgn(z.ptr))
}

// Cmp compares x and y. The result is
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
func (x *Int) Cmp(y *Int) int {
	x.doinit()
	y.doinit()
	switch cmp := int(C.mpz_cmp(x.ptr, y.ptr)); {
	case cmp < 0:
		return -1
	case cmp == 0:
		return 0
	}
	return 1
}
