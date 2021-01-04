// Copyright (c) 2016-2019 Isis Agora Lovecruft, Henry de Valence. All rights reserved.
// Copyright (c) 2020-2021 Oasis Labs Inc.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// 1. Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright
// notice, this list of conditions and the following disclaimer in the
// documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
// IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED
// TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package scalar

import "testing"

func TestUnpackedScalar(t *testing.T) {
	t.Run("Add", testUnpackedAdd)
	t.Run("Sub", testUnpackedSub)
	t.Run("Mul", testUnpackedMul)
	t.Run("Mul/Max", testUnpackedMulMax)
	t.Run("Square/Max", testUnpackedSquareMax)
	t.Run("MontgomeryMul", testUnpackedMontgomeryMul)
	t.Run("MontgomeryMul/Max", testUnpackedMontgomeryMulMax)
	t.Run("MontgomerySquare/Max", testUnpackedMontgomerySquareMax)
	t.Run("FromBytesWide", testUnpackedFromBytesWide)
	t.Run("ToFromMontgomery", testUnpackedToFromMontgomery)
	t.Run("ToFromBytes", testUnpackedToFromBytes)
}

func testUnpackedAdd(t *testing.T) {
	a, b := unpackedTestConstants["A"], unpackedTestConstants["B"]

	var res, zero unpackedScalar
	res.add(a, b)
	if res != zero {
		t.Fatalf("A+B != 0 (Got %v)", res)
	}
}

func testUnpackedSub(t *testing.T) {
	a, b := unpackedTestConstants["A"], unpackedTestConstants["B"]
	ab := *unpackedTestConstants["AB"]

	var res unpackedScalar
	res.sub(a, b)
	if res != ab {
		t.Fatalf("A-B != AB (Got %v)", res)
	}
}

func testUnpackedMul(t *testing.T) {
	x, y := unpackedTestConstants["X"], unpackedTestConstants["Y"]
	xy := *unpackedTestConstants["XY"]

	var res unpackedScalar
	res.mul(x, y)
	if res != xy {
		t.Fatalf("X*Y != XY (Got %v)", res)
	}
}

func testUnpackedMulMax(t *testing.T) {
	x := unpackedTestConstants["X"]
	xx := *unpackedTestConstants["XX"]

	var res unpackedScalar
	res.mul(x, x)
	if res != xx {
		t.Fatalf("X*X != XX (Got %v)", res)
	}
}

func testUnpackedSquareMax(t *testing.T) {
	x := unpackedTestConstants["X"]
	xx := *unpackedTestConstants["XX"]

	res := *x
	res.square()
	if res != xx {
		t.Fatalf("X*X != XX (Got %v)", res)
	}
}

func testUnpackedMontgomeryMul(t *testing.T) {
	x, y := unpackedTestConstants["X"], unpackedTestConstants["Y"]
	xyMont := *unpackedTestConstants["XY_MONT"]

	var res unpackedScalar
	res.montgomeryMul(x, y)
	if res != xyMont {
		t.Fatalf("X*Y / R != XY_MONT (Got %v)", res)
	}
}

func testUnpackedMontgomeryMulMax(t *testing.T) {
	x := unpackedTestConstants["X"]
	xxMont := *unpackedTestConstants["XX_MONT"]

	var res unpackedScalar
	res.montgomeryMul(x, x)
	if res != xxMont {
		t.Fatalf("X*Y / R != XX_MONT (Got %v)", res)
	}
}

func testUnpackedMontgomerySquareMax(t *testing.T) {
	x := unpackedTestConstants["X"]
	xxMont := *unpackedTestConstants["XX_MONT"]

	res := *x
	res.montgomerySquare()
	if res != xxMont {
		t.Fatalf("X*X / R != XX_MONT (Got %v)", res)
	}
}

func testUnpackedFromBytesWide(t *testing.T) {
	c := *unpackedTestConstants["C"]

	var bignum [ScalarWideSize]byte // 2^512 - 1
	for i := range bignum {
		bignum[i] = 255
	}

	var reduced unpackedScalar
	if err := reduced.fromBytesWide(bignum[:]); err != nil {
		t.Fatalf("FromBytesWide(bignum): %v", err)
	}
	if reduced != c {
		t.Fatalf("FromBytesWide(bignum) != C (Got %v)", reduced)
	}
}

func testUnpackedToFromMontgomery(t *testing.T) {
	y := *unpackedTestConstants["Y"]

	// At least test if this round-trips.
	tmp := y

	tmp.toMontgomery()
	if tmp == y {
		t.Fatalf("Y.ToMontgomery() = Y (Got %v)", tmp)
	}
	tmp.fromMontgomery()
	if tmp != y {
		t.Fatalf("tmp.FromMontgomery() != Y (Got %v)", tmp)
	}
}

func testUnpackedToFromBytes(t *testing.T) {
	x, y := unpackedTestConstants["X"], unpackedTestConstants["Y"]

	var (
		tmp unpackedScalar
		out [ScalarSize]byte
	)

	// X is not in canonical form, but unpackedScalar's s11n routines
	// do not reduce.
	x.toBytes(out[:])
	tmp.fromBytes(out[:])
	if tmp != *x {
		t.Fatalf("tmp.FromBytes(X.ToBytes) != X (Got %v)", tmp)
	}

	y.toBytes(out[:])
	tmp.fromBytes(out[:])
	if tmp != *y {
		t.Fatalf("tmp.FromBytes(Y.ToBytes) != Y (Got %v)", tmp)
	}
}
