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

// An Int represents a signed multi-precision integer.
// The zero value for an Int represents the value 0.
type Int struct {
	i    C.mpz_t
	init bool
}

// NewInt returns a new Int initialized to x.
func NewInt(x int64) *Int { return new(Int).SetInt64(x) }

// Int promises that the zero value is a 0, but in gmp
// the zero value is a crash.  To bridge the gap, the
// init bool says whether this is a valid gmp value.
// doinit initializes z.i if it needs it.  This is not inherent
// to FFI, just a mismatch between Go's convention of
// making zero values useful and gmp's decision not to.
func (z *Int) doinit() {
	if z.init {
		return
	}
	z.init = true
	C.mpz_init(&z.i[0])
}

// Bytes returns z's representation as a big-endian byte array.
func (z *Int) Bytes() []byte {
	b := make([]byte, (z.Len()+7)/8)
	n := C.size_t(len(b))
	C.mpz_export(unsafe.Pointer(&b[0]), &n, 1, 1, 1, 0, &z.i[0])
	return b[0:n]
}

// Len returns the length of z in bits.  0 is considered to have length 1.
func (z *Int) Len() int {
	z.doinit()
	return int(C.mpz_sizeinbase(&z.i[0], 2))
}

// Set sets z = x and returns z.
func (z *Int) Set(x *Int) *Int {
	z.doinit()
	C.mpz_set(&z.i[0], &x.i[0])
	return z
}

// SetBytes interprets b as the bytes of a big-endian integer
// and sets z to that value.
func (z *Int) SetBytes(b []byte) *Int {
	z.doinit()
	if len(b) == 0 {
		z.SetInt64(0)
	} else {
		C.mpz_import(&z.i[0], C.size_t(len(b)), 1, 1, 1, 0,
			unsafe.Pointer(&b[0]))
	}
	return z
}

// SetInt64 sets z = x and returns z.
func (z *Int) SetInt64(x int64) *Int {
	z.doinit()
	// TODO(rsc): more work on 32-bit platforms
	C.mpz_set_si(&z.i[0], C.long(x))
	return z
}

// SetString interprets s as a number in the given base
// and sets z to that value.  The base must be in the range [2,36].
// SetString returns an error if s cannot be parsed or the base is invalid.
func (z *Int) SetString(s string, base int) error {
	z.doinit()
	if base < 2 || base > 36 {
		return os.ErrInvalid
	}
	p := C.CString(s)
	defer C.free(unsafe.Pointer(p))
	if C.mpz_set_str(&z.i[0], p, C.int(base)) < 0 {
		return os.ErrInvalid
	}
	return nil
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
	p := C.mpz_get_str(nil, C.int(base), &z.i[0])
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))
	return s, nil
}

func (z *Int) destroy() {
	if z.init {
		C.mpz_clear(&z.i[0])
	}
	z.init = false
}

func (z *Int) Clear() {
	z.destroy()
}

/*
 * arithmetic
 */

// Add sets z = x + y and returns z.
func (z *Int) Add(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_add(&z.i[0], &x.i[0], &y.i[0])
	return z
}

// Sub sets z = x - y and returns z.
func (z *Int) Sub(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_sub(&z.i[0], &x.i[0], &y.i[0])
	return z
}

// Mul sets z = x * y and returns z.
func (z *Int) Mul(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_mul(&z.i[0], &x.i[0], &y.i[0])
	return z
}

// Div sets z = x / y, rounding toward zero, and returns z.
func (z *Int) Div(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_tdiv_q(&z.i[0], &x.i[0], &y.i[0])
	return z
}

// Mod sets z = x % y and returns z.
// Like the result of the Go % operator, z has the same sign as x.
func (z *Int) Mod(x, y *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	C.mpz_tdiv_r(&z.i[0], &x.i[0], &y.i[0])
	return z
}

// ModInverse sets z to the multiplicative inverse of g in the group ℤ/pℤ
// (where p is a prime) and returns z.
func (z *Int) ModInverse(g, p *Int) *Int {
	g.doinit()
	p.doinit()
	z.doinit()
	C.mpz_invert(&z.i[0], &g.i[0], &p.i[0])
	return z
}

// Lsh sets z = x << s and returns z.
func (z *Int) Lsh(x *Int, s uint) *Int {
	x.doinit()
	z.doinit()
	C._mpz_mul_2exp(&z.i[0], &x.i[0], C.ulong(s))
	return z
}

// Rsh sets z = x >> s and returns z.
func (z *Int) Rsh(x *Int, s uint) *Int {
	x.doinit()
	z.doinit()
	C._mpz_div_2exp(&z.i[0], &x.i[0], C.ulong(s))
	return z
}

// Exp sets z = x^y % m and returns z.
// If m == nil, Exp sets z = x^y.
func (z *Int) Exp(x, y, m *Int) *Int {
	x.doinit()
	y.doinit()
	z.doinit()
	if m == nil {
		C.mpz_pow_ui(&z.i[0], &x.i[0], C.mpz_get_ui(&y.i[0]))
	} else {
		m.doinit()
		C.mpz_powm(&z.i[0], &x.i[0], &y.i[0], &m.i[0])
	}
	return z
}

// Sqrt sets z = floor(sqrt(x)) and returns z.
func (z *Int) Sqrt(x *Int) *Int {
	z.doinit()
	x.doinit()
	C.mpz_sqrt(&z.i[0], &x.i[0])
	return z
}

func (z *Int) Int64() int64 {
	if !z.init {
		return 0
	}
	return int64(C.mpz_get_si(&z.i[0]))
}

// Neg sets z = -x and returns z.
func (z *Int) Neg(x *Int) *Int {
	x.doinit()
	z.doinit()
	C.mpz_neg(&z.i[0], &x.i[0])
	return z
}

// Abs sets z to the absolute value of x and returns z.
func (z *Int) Abs(x *Int) *Int {
	x.doinit()
	z.doinit()
	C.mpz_abs(&z.i[0], &x.i[0])
	return z
}

// Sign returns:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0
//
func (z *Int) Sign() int {
	z.doinit()
	return int(C._mpz_sgn(&z.i[0]))
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
	switch cmp := C.mpz_cmp(&x.i[0], &y.i[0]); {
	case cmp < 0:
		return -1
	case cmp == 0:
		return 0
	}
	return +1
}

/*
 * functions without a clear receiver
 */

// DivModInt sets q = x / y and r = x % y.
func DivModInt(q, r, x, y *Int) {
	q.doinit()
	r.doinit()
	x.doinit()
	y.doinit()
	C.mpz_tdiv_qr(&q.i[0], &r.i[0], &x.i[0], &y.i[0])
}

// GcdInt sets d to the greatest common divisor of a and b,
// which must be positive numbers.
// If x and y are not nil, GcdInt sets x and y such that d = a*x + b*y.
// If either a or b is not positive, GcdInt sets d = x = y = 0.
func GcdInt(d, x, y, a, b *Int) {
	d.doinit()
	x.doinit()
	y.doinit()
	a.doinit()
	b.doinit()
	C.mpz_gcdext(&d.i[0], &x.i[0], &y.i[0], &a.i[0], &b.i[0])
}

// ProbablyPrime performs n Miller-Rabin tests to check whether z is prime.
// If it returns true, z is prime with probability 1 - 1/4^n.
// If it returns false, z is not prime.
func (z *Int) ProbablyPrime(n int) bool {
	z.doinit()
	return int(C.mpz_probab_prime_p(&z.i[0], C.int(n))) > 0
}
