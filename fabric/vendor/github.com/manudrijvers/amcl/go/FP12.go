/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

/* MiotCL Fp^12 functions */
/* FP12 elements are of the form a+i.b+i^2.c */

package amcl

//import "fmt"

type FP12 struct {
	a *FP4
	b *FP4
	c *FP4
}

/* Constructors */
func NewFP12fp4(d *FP4) *FP12 {
	F := new(FP12)
	F.a = NewFP4copy(d)
	F.b = NewFP4int(0)
	F.c = NewFP4int(0)
	return F
}

func NewFP12int(d int) *FP12 {
	F := new(FP12)
	F.a = NewFP4int(d)
	F.b = NewFP4int(0)
	F.c = NewFP4int(0)
	return F
}

func NewFP12fp4s(d *FP4, e *FP4, f *FP4) *FP12 {
	F := new(FP12)
	F.a = NewFP4copy(d)
	F.b = NewFP4copy(e)
	F.c = NewFP4copy(f)
	return F
}

func NewFP12copy(x *FP12) *FP12 {
	F := new(FP12)
	F.a = NewFP4copy(x.a)
	F.b = NewFP4copy(x.b)
	F.c = NewFP4copy(x.c)
	return F
}

/* reduce all components of this mod Modulus */
func (F *FP12) reduce() {
	F.a.reduce()
	F.b.reduce()
	F.c.reduce()
}

/* normalise all components of this */
func (F *FP12) norm() {
	F.a.norm()
	F.b.norm()
	F.c.norm()
}

/* test x==0 ? */
func (F *FP12) iszilch() bool {
	F.reduce()
	return (F.a.iszilch() && F.b.iszilch() && F.c.iszilch())
}

/* test x==1 ? */
func (F *FP12) isunity() bool {
	one := NewFP4int(1)
	return (F.a.equals(one) && F.b.iszilch() && F.c.iszilch())
}

/* return 1 if x==y, else 0 */
func (F *FP12) equals(x *FP12) bool {
	return (F.a.equals(x.a) && F.b.equals(x.b) && F.c.equals(x.c))
}

/* extract a from this */
func (F *FP12) geta() *FP4 {
	return F.a
}

/* extract b */
func (F *FP12) getb() *FP4 {
	return F.b
}

/* extract c */
func (F *FP12) getc() *FP4 {
	return F.c
}

/* copy this=x */
func (F *FP12) copy(x *FP12) {
	F.a.copy(x.a)
	F.b.copy(x.b)
	F.c.copy(x.c)
}

/* set this=1 */
func (F *FP12) one() {
	F.a.one()
	F.b.zero()
	F.c.zero()
}

/* this=conj(this) */
func (F *FP12) conj() {
	F.a.conj()
	F.b.nconj()
	F.c.conj()
}

/* Granger-Scott Unitary Squaring */
func (F *FP12) usqr() {
	A := NewFP4copy(F.a)
	B := NewFP4copy(F.c)
	C := NewFP4copy(F.b)
	D := NewFP4int(0)

	F.a.sqr()
	D.copy(F.a)
	D.add(F.a)
	F.a.add(D)

	F.a.norm()
	A.nconj()

	A.add(A)
	F.a.add(A)
	B.sqr()
	B.times_i()

	D.copy(B)
	D.add(B)
	B.add(D)
	B.norm()

	C.sqr()
	D.copy(C)
	D.add(C)
	C.add(D)
	C.norm()

	F.b.conj()
	F.b.add(F.b)
	F.c.nconj()

	F.c.add(F.c)
	F.b.add(B)
	F.c.add(C)
	F.reduce()

}

/* Chung-Hasan SQR2 method from http://cacr.uwaterloo.ca/techreports/2006/cacr2006-24.pdf */
func (F *FP12) sqr() {
	A := NewFP4copy(F.a)
	B := NewFP4copy(F.b)
	C := NewFP4copy(F.c)
	D := NewFP4copy(F.a)

	A.sqr()
	B.mul(F.c)
	B.add(B)
	C.sqr()
	D.mul(F.b)
	D.add(D)

	F.c.add(F.a)
	F.c.add(F.b)
	F.c.sqr()

	F.a.copy(A)

	A.add(B)
	A.norm()
	A.add(C)
	A.add(D)
	A.norm()

	A.neg()
	B.times_i()
	C.times_i()

	F.a.add(B)

	F.b.copy(C)
	F.b.add(D)
	F.c.add(A)
	F.norm()
}

/* FP12 full multiplication this=this*y */
func (F *FP12) mul(y *FP12) {
	z0 := NewFP4copy(F.a)
	z1 := NewFP4int(0)
	z2 := NewFP4copy(F.b)
	z3 := NewFP4int(0)
	t0 := NewFP4copy(F.a)
	t1 := NewFP4copy(y.a)

	z0.mul(y.a)
	z2.mul(y.b)

	t0.add(F.b)
	t1.add(y.b)

	z1.copy(t0)
	z1.mul(t1)
	t0.copy(F.b)
	t0.add(F.c)

	t1.copy(y.b)
	t1.add(y.c)
	z3.copy(t0)
	z3.mul(t1)

	t0.copy(z0)
	t0.neg()
	t1.copy(z2)
	t1.neg()

	z1.add(t0)
	z1.norm()
	F.b.copy(z1)
	F.b.add(t1)

	z3.add(t1)
	z2.add(t0)

	t0.copy(F.a)
	t0.add(F.c)
	t1.copy(y.a)
	t1.add(y.c)
	t0.mul(t1)
	z2.add(t0)

	t0.copy(F.c)
	t0.mul(y.c)
	t1.copy(t0)
	t1.neg()

	z2.norm()
	z3.norm()
	F.b.norm()

	F.c.copy(z2)
	F.c.add(t1)
	z3.add(t1)
	t0.times_i()
	F.b.add(t0)

	z3.times_i()
	F.a.copy(z0)
	F.a.add(z3)
	F.norm()
}

/* Special case of multiplication arises from special form of ATE pairing line function */
func (F *FP12) smul(y *FP12) {
	z0 := NewFP4copy(F.a)
	z2 := NewFP4copy(F.b)
	z3 := NewFP4copy(F.b)
	t0 := NewFP4int(0)
	t1 := NewFP4copy(y.a)

	z0.mul(y.a)
	z2.pmul(y.b.real())
	F.b.add(F.a)
	t1.real().add(y.b.real())

	F.b.mul(t1)
	z3.add(F.c)
	z3.pmul(y.b.real())

	t0.copy(z0)
	t0.neg()
	t1.copy(z2)
	t1.neg()

	F.b.add(t0)
	F.b.norm()

	F.b.add(t1)
	z3.add(t1)
	z2.add(t0)

	t0.copy(F.a)
	t0.add(F.c)
	t0.mul(y.a)
	F.c.copy(z2)
	F.c.add(t0)

	z3.times_i()
	F.a.copy(z0)
	F.a.add(z3)

	F.norm()
}

/* this=1/this */
func (F *FP12) inverse() {
	f0 := NewFP4copy(F.a)
	f1 := NewFP4copy(F.b)
	f2 := NewFP4copy(F.a)
	f3 := NewFP4int(0)

	F.norm()
	f0.sqr()
	f1.mul(F.c)
	f1.times_i()
	f0.sub(f1)

	f1.copy(F.c)
	f1.sqr()
	f1.times_i()
	f2.mul(F.b)
	f1.sub(f2)

	f2.copy(F.b)
	f2.sqr()
	f3.copy(F.a)
	f3.mul(F.c)
	f2.sub(f3)

	f3.copy(F.b)
	f3.mul(f2)
	f3.times_i()
	F.a.mul(f0)
	f3.add(F.a)
	F.c.mul(f1)
	F.c.times_i()

	f3.add(F.c)
	f3.inverse()
	F.a.copy(f0)
	F.a.mul(f3)
	F.b.copy(f1)
	F.b.mul(f3)
	F.c.copy(f2)
	F.c.mul(f3)
}

/* this=this^p using Frobenius */
func (F *FP12) frob(f *FP2) {
	f2 := NewFP2copy(f)
	f3 := NewFP2copy(f)

	f2.sqr()
	f3.mul(f2)

	F.a.frob(f3)
	F.b.frob(f3)
	F.c.frob(f3)

	F.b.pmul(f)
	F.c.pmul(f2)
}

/* trace function */
func (F *FP12) trace() *FP4 {
	t := NewFP4int(0)
	t.copy(F.a)
	t.imul(3)
	t.reduce()
	return t
}

/* convert from byte array to FP12 */
func FP12_fromBytes(w []byte) *FP12 {
	var t [int(MODBYTES)]byte
	MB := int(MODBYTES)

	for i := 0; i < MB; i++ {
		t[i] = w[i]
	}
	a := fromBytes(t[:])
	for i := 0; i < MB; i++ {
		t[i] = w[i+MB]
	}
	b := fromBytes(t[:])
	c := NewFP2bigs(a, b)

	for i := 0; i < MB; i++ {
		t[i] = w[i+2*MB]
	}
	a = fromBytes(t[:])
	for i := 0; i < MB; i++ {
		t[i] = w[i+3*MB]
	}
	b = fromBytes(t[:])
	d := NewFP2bigs(a, b)

	e := NewFP4fp2s(c, d)

	for i := 0; i < MB; i++ {
		t[i] = w[i+4*MB]
	}
	a = fromBytes(t[:])
	for i := 0; i < MB; i++ {
		t[i] = w[i+5*MB]
	}
	b = fromBytes(t[:])
	c = NewFP2bigs(a, b)

	for i := 0; i < MB; i++ {
		t[i] = w[i+6*MB]
	}
	a = fromBytes(t[:])
	for i := 0; i < MB; i++ {
		t[i] = w[i+7*MB]
	}
	b = fromBytes(t[:])
	d = NewFP2bigs(a, b)

	f := NewFP4fp2s(c, d)

	for i := 0; i < MB; i++ {
		t[i] = w[i+8*MB]
	}
	a = fromBytes(t[:])
	for i := 0; i < MB; i++ {
		t[i] = w[i+9*MB]
	}
	b = fromBytes(t[:])

	c = NewFP2bigs(a, b)

	for i := 0; i < MB; i++ {
		t[i] = w[i+10*MB]
	}
	a = fromBytes(t[:])
	for i := 0; i < MB; i++ {
		t[i] = w[i+11*MB]
	}
	b = fromBytes(t[:])
	d = NewFP2bigs(a, b)

	g := NewFP4fp2s(c, d)

	return NewFP12fp4s(e, f, g)
}

/* convert this to byte array */
func (F *FP12) toBytes(w []byte) {
	var t [int(MODBYTES)]byte
	MB := int(MODBYTES)
	F.a.geta().getA().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i] = t[i]
	}
	F.a.geta().getB().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+MB] = t[i]
	}
	F.a.getb().getA().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+2*MB] = t[i]
	}
	F.a.getb().getB().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+3*MB] = t[i]
	}

	F.b.geta().getA().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+4*MB] = t[i]
	}
	F.b.geta().getB().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+5*MB] = t[i]
	}
	F.b.getb().getA().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+6*MB] = t[i]
	}
	F.b.getb().getB().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+7*MB] = t[i]
	}

	F.c.geta().getA().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+8*MB] = t[i]
	}
	F.c.geta().getB().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+9*MB] = t[i]
	}
	F.c.getb().getA().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+10*MB] = t[i]
	}
	F.c.getb().getB().toBytes(t[:])
	for i := 0; i < MB; i++ {
		w[i+11*MB] = t[i]
	}
}

/* convert to hex string */
func (F *FP12) toString() string {
	return ("[" + F.a.toString() + "," + F.b.toString() + "," + F.c.toString() + "]")
}

/* this=this^e */
func (F *FP12) pow(e *BIG) *FP12 {
	F.norm()
	e.norm()
	w := NewFP12copy(F)
	z := NewBIGcopy(e)
	r := NewFP12int(1)

	for true {
		bt := z.parity()
		z.fshr(1)
		if bt == 1 {
			r.mul(w)
		}
		if z.iszilch() {
			break
		}
		w.usqr()
	}
	r.reduce()
	return r
}

/* constant time powering by small integer of max length bts */
func (F *FP12) pinpow(e int, bts int) {
	var R []*FP12
	R = append(R, NewFP12int(1))
	R = append(R, NewFP12copy(F))

	for i := bts - 1; i >= 0; i-- {
		b := (e >> uint(i)) & 1
		R[1-b].mul(R[b])
		R[b].usqr()
	}
	F.copy(R[0])
}

/* p=q0^u0.q1^u1.q2^u2.q3^u3 */
/* Timing attack secure, but not cache attack secure */

func pow4(q []*FP12, u []*BIG) *FP12 {
	var a [4]int8
	var g []*FP12
	var s []*FP12
	c := NewFP12int(1)
	p := NewFP12int(0)
	var w [NLEN*int(BASEBITS) + 1]int8
	var t []*BIG
	mt := NewBIGint(0)

	for i := 0; i < 4; i++ {
		t = append(t, NewBIGcopy(u[i]))
	}

	s = append(s, NewFP12int(0))
	s = append(s, NewFP12int(0))

	g = append(g, NewFP12copy(q[0]))
	s[0].copy(q[1])
	s[0].conj()
	g[0].mul(s[0])
	g = append(g, NewFP12copy(g[0]))
	g = append(g, NewFP12copy(g[0]))
	g = append(g, NewFP12copy(g[0]))
	g = append(g, NewFP12copy(q[0]))
	g[4].mul(q[1])
	g = append(g, NewFP12copy(g[4]))
	g = append(g, NewFP12copy(g[4]))
	g = append(g, NewFP12copy(g[4]))

	s[1].copy(q[2])
	s[0].copy(q[3])
	s[0].conj()
	s[1].mul(s[0])
	s[0].copy(s[1])
	s[0].conj()
	g[1].mul(s[0])
	g[2].mul(s[1])
	g[5].mul(s[0])
	g[6].mul(s[1])
	s[1].copy(q[2])
	s[1].mul(q[3])
	s[0].copy(s[1])
	s[0].conj()
	g[0].mul(s[0])
	g[3].mul(s[1])
	g[4].mul(s[0])
	g[7].mul(s[1])

	/* if power is even add 1 to power, and add q to correction */

	for i := 0; i < 4; i++ {
		if t[i].parity() == 0 {
			t[i].inc(1)
			t[i].norm()
			c.mul(q[i])
		}
		mt.add(t[i])
		mt.norm()
	}
	c.conj()
	nb := 1 + mt.nbits()

	/* convert exponent to signed 1-bit window */
	for j := 0; j < nb; j++ {
		for i := 0; i < 4; i++ {
			a[i] = int8(t[i].lastbits(2) - 2)
			t[i].dec(int(a[i]))
			t[i].norm()
			t[i].fshr(1)
		}
		w[j] = (8*a[0] + 4*a[1] + 2*a[2] + a[3])
	}
	w[nb] = int8(8*t[0].lastbits(2) + 4*t[1].lastbits(2) + 2*t[2].lastbits(2) + t[3].lastbits(2))
	p.copy(g[(w[nb]-1)/2])

	for i := nb - 1; i >= 0; i-- {
		m := w[i] >> 7
		j := (w[i] ^ m) - m /* j=abs(w[i]) */
		j = (j - 1) / 2
		s[0].copy(g[j])
		s[1].copy(g[j])
		s[1].conj()
		p.usqr()
		p.mul(s[m&1])
	}
	p.mul(c) /* apply correction */
	p.reduce()
	return p
}
