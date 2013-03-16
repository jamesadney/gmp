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
	"unsafe"
)

type Rat struct {
	i    C.mpq_t
	init bool
}

// NewRat returns a new Rat initialized to x/y.
func NewRat(x int64, y uint) *Rat { return new(Rat).SetInt64(x, y) }

// Int promises that the zero value is a 0, but in gmp
// the zero value is a crash.  To bridge the gap, the
// init bool says whether this is a valid gmp value.
// doinit initializes z.i if it needs it.  This is not inherent
// to FFI, just a mismatch between Go's convention of
// making zero values useful and gmp's decision not to.
func (q *Rat) doinit() {
	if q.init {
		return
	}
	q.init = true
	C.mpq_init(&q.i[0])
}

// Set sets z = x and returns z.
func (q *Rat) Set(x *Rat) *Rat {
	q.doinit()
	C.mpq_set(&q.i[0], &x.i[0])
	return q
}

// SetInt64 sets q to x/y and returns q.
func (q *Rat) SetInt64(x int64, y uint) *Rat {
	q.doinit()
	C.mpq_set_si(&q.i[0], C.long(x), C.ulong(y))
	return q
}

// SetUint sets q to x/y and returns q.
func (q *Rat) SetUint(x, y uint) *Rat {
	q.doinit()
	C.mpq_set_ui(&q.i[0], C.ulong(x), C.ulong(y))
	return q
}

// SetInt sets q to x and returns q.
func (q *Rat) SetInt(x *Int) *Rat {
	q.doinit()
	x.doinit()
	C.mpq_set_z(&q.i[0], &x.i[0])
	return q
}

// SetString interprets s as a number in the given base
// and sets z to that value.  The base must be in the range [2,36].
// SetString returns an error if s cannot be parsed or the base is invalid.
func (q *Rat) SetString(s string, base int) error {
	q.doinit()
	if base < 2 || base > 36 {
		return os.ErrInvalid
	}
	p := C.CString(s)
	defer C.free(unsafe.Pointer(p))
	if C.mpq_set_str(&q.i[0], p, C.int(base)) < 0 {
		return os.ErrInvalid
	}
	return nil
}

func SwapRat(x, y *Rat) {
	x.doinit()
	y.doinit()
	C.mpq_swap(&x.i[0], &y.i[0])
}

// String returns the decimal representation of z.
func (q *Rat) StringBase(base int) (string, error) {
	if q == nil {
		return "nil", nil
	}
	if base < 2 || base > 36 {
		return "", os.ErrInvalid
	}
	q.doinit()
	p := C.mpq_get_str(nil, C.int(base), &q.i[0])
	s := C.GoString(p)
	C.free(unsafe.Pointer(p))
	return s, nil
}

func (q *Rat) String() string {
	s, _ := q.StringBase(10)
	return s
}

func (q *Rat) Float64() float64 {
	q.doinit()
	return float64(C.mpq_get_d(&q.i[0]))
}

// SetFloat64 sets f = x and returns q.
func (q *Rat) SetFloat64(x float64) *Rat {
	q.doinit()
	C.mpq_set_d(&q.i[0], C.double(x))
	return q
}

// SetFloat sets f = x and returns f.
func (q *Rat) SetFloat(x *Float) *Rat {
	q.doinit()
	C.mpq_set_f(&q.i[0], &x.i[0])
	return q
}

func (q *Rat) destroy() {
	if q.init {
		C.mpq_clear(&q.i[0])
	}
	q.init = false
}

func (q *Rat) Clear() {
	q.destroy()
}

// Add sets q = x + y and returns f.
func (q *Rat) Add(x, y *Rat) *Rat {
	x.doinit()
	y.doinit()
	q.doinit()
	C.mpq_add(&q.i[0], &x.i[0], &y.i[0])
	return q
}

func (q *Rat) Sub(x, y *Rat) *Rat {
	x.doinit()
	y.doinit()
	q.doinit()
	C.mpq_sub(&q.i[0], &x.i[0], &y.i[0])
	return q
}

func (q *Rat) Mul(x, y *Rat) *Rat {
	x.doinit()
	y.doinit()
	q.doinit()
	C.mpq_mul(&q.i[0], &x.i[0], &y.i[0])
	return q
}

func (q *Rat) Div(x, y *Rat) *Rat {
	x.doinit()
	y.doinit()
	q.doinit()
	C.mpq_div(&q.i[0], &x.i[0], &y.i[0])
	return q
}

func (q *Rat) Abs(x *Rat) *Rat {
	x.doinit()
	q.doinit()
	C.mpq_abs(&q.i[0], &x.i[0])
	return q
}

func (q *Rat) Inv(x *Rat) *Rat {
	x.doinit()
	q.doinit()
	C.mpq_inv(&q.i[0], &x.i[0])
	return q
}

// Mul2Exp sets z = x * 2^s and returns z.
func (q *Rat) Mul2Exp(x *Rat, s uint) *Rat {
	x.doinit()
	q.doinit()
	C.mpq_mul_2exp(&q.i[0], &x.i[0], C.mp_bitcnt_t(s))
	return q
}

// Div2Exp sets z = x / 2^s and returns z.
func (q *Rat) Div2Exp(x *Rat, s uint) *Rat {
	x.doinit()
	q.doinit()
	C.mpq_div_2exp(&q.i[0], &x.i[0], C.mp_bitcnt_t(s))
	return q
}

func CmpRat(x, y *Rat) int {
	x.doinit()
	y.doinit()
	return int(C.mpq_cmp(&x.i[0], &y.i[0]))
}

func CmpRatUint(q *Rat, x, y uint) int {
	q.doinit()
	return 0 // FIXME(ug): Macro...
	//return int(C.mpq_cmp_ui(&x.i[0], C.ulong(x), C.ulong(y)))
}

func CmpRatInt64(q *Rat, x int64, y uint) int {
	q.doinit()
	return 0 // FIXME(ug): Macro...
	//return int(C.mpq_cmp_ui(&x.i[0], C.long(x), C.ulong(y)))
}

func (q *Rat) Sgn() int {
	q.doinit()
	//TODO(ug): mpf_sgn seems to be implemented as a macro.
	// We need to watch out for changes in the data structure :(

	//return int(C.mpq_sgn(&f.i[0]))
	switch size := int(q.i[0]._mp_num._mp_size); {
	case size < 0:
		return -1
	case size == 0:
		return 0
	}
	return 1
}

func EqRat(x, y *Rat) bool {
	x.doinit()
	y.doinit()
	return C.mpq_equal(&x.i[0], &y.i[0]) != 0
}

func (q *Rat) Num(n *Int) *Int {
	q.doinit()
	n.doinit()
	C.mpq_get_num(&n.i[0], &q.i[0])
	return n
}

func (q *Rat) Den(n *Int) *Int {
	q.doinit()
	n.doinit()
	C.mpq_get_den(&n.i[0], &q.i[0])
	return n
}

func (q *Rat) SetNum(n *Int) *Rat {
	q.doinit()
	n.doinit()
	C.mpq_set_num(&q.i[0], &n.i[0])
	return q
}

func (q *Rat) SetDen(n *Int) *Rat {
	q.doinit()
	n.doinit()
	C.mpq_set_den(&q.i[0], &n.i[0])
	return q
}

// TODO(ug): mpq_numref, mpq_denref
