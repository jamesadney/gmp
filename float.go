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
*/
import "C"

import (
	"os"
	"strconv"
	"unsafe"
)

type Float struct {
	i    C.mpf_t
	init bool
	prec uint // 0 = use the default precision
}

// NewInt returns a new Int initialized to x.
func NewFloat(x float64) *Float { return new(Float).SetFloat64(x) }

// NewInt returns a new Int initialized to x, with precision prec.
func NewFloat2(x float64, prec uint) *Float {
	f := new(Float)
	f.prec = prec
	f.SetFloat64(x)
	return f
}

// Int promises that the zero value is a 0, but in gmp
// the zero value is a crash.  To bridge the gap, the
// init bool says whether this is a valid gmp value.
// doinit initializes f.i if it needs it.  This is not inherent
// to FFI, just a mismatch between Go's convention of
// making zero values useful and gmp's decision not to.
func (f *Float) doinit() {
	if f.init {
		return
	}
	if f.prec != 0 {
		C.mpf_init2(&f.i[0], C.mp_bitcnt_t(f.prec))
	} else {
		C.mpf_init(&f.i[0])
	}
	f.init = true
}

// Set sets f = x and returns f.
func (f *Float) Set(x *Float) *Float {
	f.doinit()
	C.mpf_set(&f.i[0], &x.i[0])
	return f
}

// SetInt sets f = x and returns f.
func (f *Float) SetInt64(x int64) *Float {
	f.doinit()
	C.mpf_set_si(&f.i[0], C.long(x))
	return f
}

// SetFloat64 sets f = x and returns f.
func (f *Float) SetFloat64(x float64) *Float {
	f.doinit()
	C.mpf_set_d(&f.i[0], C.double(x))
	return f
}

// SetString interprets s as a number in the given base
// and sets f to that value.  The base must be in the range [2,36].
// SetString returns an error if s cannot be parsed or the base is invalid.
func (f *Float) SetString(s string, base int) error {
	f.doinit()
	if base < 2 || base > 36 {
		return os.ErrInvalid
	}
	p := C.CString(s)
	defer C.free(unsafe.Pointer(p))
	if C.mpf_set_str(&f.i[0], p, C.int(base)) < 0 {
		return os.ErrInvalid
	}
	return nil
}

func (f *Float) StringBase(base int, ndigits uint) (string, error) {
	if f == nil {
		return "nil", nil
	}
	if base < 2 || base > 36 {
		return "", os.ErrInvalid
	}
	f.doinit()
	var exp_ C.mp_exp_t
	p := C.mpf_get_str(nil, &exp_, C.int(base), C.size_t(ndigits), &f.i[0])
	exp := int(exp_)
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))

	if len(s) == 0 {
		return "0", nil
	}

	if exp > 0 && exp < len(s) {
		return s[:exp] + "." + s[exp:], nil
	}
	return s + "e" + strconv.Itoa(exp), nil
}

// String returns the decimal representation of z.
func (f *Float) String() string {
	s, _ := f.StringBase(10, 0)
	return s
}

func (f *Float) Float64() float64 {
	f.doinit()
	return float64(C.mpf_get_d(&f.i[0]))
}

func (f *Float) Int64() int64 {
	f.doinit()
	return int64(C.mpf_get_si(&f.i[0]))
}

// FIXME: Float2Exp is inconsistent, Float642Exp is silly.

// Convert f to a float64, truncating if necessary (ie. rounding
// towards zero), and with an exponent returned separately.
func (f *Float) Float2Exp() (d float64, exp int) {
	var exp_ C.long
	d = float64(C.mpf_get_d_2exp(&exp_, &f.i[0]))
	exp = int(exp_)
	return
}

func (f *Float) destroy() {
	if f.init {
		C.mpf_clear(&f.i[0])
	}
	f.init = false
}

func (f *Float) Clear() {
	f.destroy()
}

func (f *Float) GetPrec() uint {
	f.doinit()
	return uint(C.mpf_get_prec(&f.i[0]))
}

func (f *Float) SetPrec(prec uint) {
	f.doinit()
	C.mpf_set_prec(&f.i[0], C.mp_bitcnt_t(prec))
	f.prec = prec
}

func (f *Float) SetPrecRaw(prec uint) {
	f.doinit()
	C.mpf_set_prec_raw(&f.i[0], C.mp_bitcnt_t(prec))
}

func SetDefaultPrec(prec uint) {
	C.mpf_set_default_prec(C.mp_bitcnt_t(prec))
}

func GetDefaultPrec() uint {
	return uint(C.mpf_get_default_prec())
}

/*
 * arithmetic
 */

// Add sets f = x + y and returns f.
func (f *Float) Add(x, y *Float) *Float {
	x.doinit()
	y.doinit()
	f.doinit()
	C.mpf_add(&f.i[0], &x.i[0], &y.i[0])
	return f
}

func (f *Float) AddUint(x *Float, y uint) *Float {
	x.doinit()
	f.doinit()
	C.mpf_add_ui(&f.i[0], &x.i[0], C.ulong(y))
	return f
}

// Sub sets f = x - y and returns f.
func (f *Float) Sub(x, y *Float) *Float {
	x.doinit()
	y.doinit()
	f.doinit()
	C.mpf_sub(&f.i[0], &x.i[0], &y.i[0])
	return f
}

func (f *Float) SubUint(x *Float, y uint) *Float {
	x.doinit()
	f.doinit()
	C.mpf_sub_ui(&f.i[0], &x.i[0], C.ulong(y))
	return f
}

// Mul sets f = x * y and returns f.
func (f *Float) Mul(x, y *Float) *Float {
	x.doinit()
	y.doinit()
	f.doinit()
	C.mpf_mul(&f.i[0], &x.i[0], &y.i[0])
	return f
}

func (f *Float) MulUint(x *Float, y uint) *Float {
	x.doinit()
	f.doinit()
	C.mpf_mul_ui(&f.i[0], &x.i[0], C.ulong(y))
	return f
}

// Div sets f = x / y and returns f.
func (f *Float) Div(x, y *Float) *Float {
	x.doinit()
	y.doinit()
	f.doinit()
	C.mpf_div(&f.i[0], &x.i[0], &y.i[0])
	return f
}

func (f *Float) DivUint(x *Float, y uint) *Float {
	x.doinit()
	f.doinit()
	C.mpf_div_ui(&f.i[0], &x.i[0], C.ulong(y))
	return f
}

func (f *Float) UintDiv(x uint, y *Float) *Float {
	y.doinit()
	f.doinit()
	C.mpf_ui_div(&f.i[0], C.ulong(x), &y.i[0])
	return f
}

// Sqrt sets f = Sqrt(x) and returns f.
func (f *Float) Sqrt(x *Float) *Float {
	x.doinit()
	f.doinit()
	C.mpf_sqrt(&f.i[0], &x.i[0])
	return f
}

// Sqrt sets f = Sqrt(x) and returns f.
func (f *Float) SqrtUint(x uint) *Float {
	f.doinit()
	C.mpf_sqrt_ui(&f.i[0], C.ulong(x))
	return f
}

// PowUint sets f = x^y and returns f
func (f *Float) PowUint(x *Float, y uint) *Float {
	x.doinit()
	f.doinit()
	C.mpf_pow_ui(&f.i[0], &x.i[0], C.ulong(y))
	return f
}

// Neg sets z = -x and returns z.
func (f *Float) Neg(x *Float) *Float {
	x.doinit()
	f.doinit()
	C.mpf_neg(&f.i[0], &x.i[0])
	return f
}

// Abs sets z to the absolute value of x and returns z.
func (f *Float) Abs(x *Float) *Float {
	x.doinit()
	f.doinit()
	C.mpf_abs(&f.i[0], &x.i[0])
	return f
}

// Mul2Exp sets z = x * 2^s and returns z.
func (f *Float) Mul2Exp(x *Float, s uint) *Float {
	x.doinit()
	f.doinit()
	C.mpf_mul_2exp(&f.i[0], &x.i[0], C.mp_bitcnt_t(s))
	return f
}

// Div2Exp sets z = x / 2^s and returns z.
func (f *Float) Div2Exp(x *Float, s uint) *Float {
	x.doinit()
	f.doinit()
	C.mpf_div_2exp(&f.i[0], &x.i[0], C.mp_bitcnt_t(s))
	return f
}

/*
 * Comparison
 */

// Compute the relative difference between x and y and store the result in f.
// This is abs(x-y)/x.
func (f *Float) RelDiff(x, y *Float) *Float {
	x.doinit()
	y.doinit()
	f.doinit()
	C.mpf_reldiff(&f.i[0], &x.i[0], &y.i[0])
	return f
}

// Return +1 if f > 0, 0 if f = 0, and -1 if f < 0.
func (f *Float) Sgn() int {
	f.doinit()
	//TODO(ug): mpf_sgn seems to be implemented as a macro.
	// We need to watch out for changes in the data structure :(

	//return int(C.mpf_sgn(&f.i[0]))
	switch size := int(f.i[0]._mp_size); {
	case size < 0:
		return -1
	case size == 0:
		return 0
	}
	return 1
}

/*
 * functions without a clear receiver
 */

// CmpInt compares x and y. The result is
//
//   neg if x <  y
//    0 if x == y
//   pos if x >  y
//
func CmpFloat(x, y *Float) int {
	x.doinit()
	y.doinit()
	return int(C.mpf_cmp(&x.i[0], &y.i[0]))
}

func CmpFloatFloat64(x *Float, y float64) int {
	x.doinit()
	return int(C.mpf_cmp_d(&x.i[0], C.double(y)))
}

func CmpFloatUint(x *Float, y uint) int {
	x.doinit()
	return int(C.mpf_cmp_ui(&x.i[0], C.ulong(y)))
}

func CmpFloatInt64(x *Float, y int64) int {
	x.doinit()
	return int(C.mpf_cmp_si(&x.i[0], C.long(y)))
}

// Return non-zero if the first n bits of x and y are equal,
// zero otherwise.  I.e., test if x and y are approximately equal.
func EqFloat(x, y *Float, n uint) int {
	x.doinit()
	y.doinit()
	return int(C.mpf_eq(&x.i[0], &y.i[0], C.mp_bitcnt_t(n)))
}

func SwapFloat(x, y *Float) {
	x.doinit()
	y.doinit()
	C.mpf_swap(&x.i[0], &y.i[0])
}

// Sets f = Ceil(x) and returns f.
func (f *Float) Ceil(x *Float) *Float {
	x.doinit()
	f.doinit()
	C.mpf_ceil(&f.i[0], &x.i[0])
	return f
}

// Sets f = Floor(x) and returns f.
func (f *Float) Floor(x *Float) *Float {
	x.doinit()
	f.doinit()
	C.mpf_floor(&f.i[0], &x.i[0])
	return f
}

// Sets f = Trunc(x) (=round towards zero) and returns f.
func (f *Float) Trunc(x *Float) *Float {
	x.doinit()
	f.doinit()
	C.mpf_trunc(&f.i[0], &x.i[0])
	return f
}

func (f *Float) IsInteger() bool {
	f.doinit()
	return int(C.mpf_integer_p(&f.i[0])) != 0
}

//TODO(ug) mpf_fits_* and random functions
