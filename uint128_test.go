// Copyright 2014, David Minor. All rights reserved.
// Use of this source code is governed by the MIT
// license which can be found in the LICENSE file.

package uint128

import (
  "math/big"
  "math/rand"
  "testing"
  "encoding/binary"
)

func makeUint128(a []uint32) Uint128 {
  return Uint128{uint64(a[0]) << 32 | uint64(a[1]), 
                uint64(a[2]) << 32 | uint64(a[3])}
}

func makeBigInt(a []uint32) *big.Int {
  buf := make([]byte, 16)
  binary.BigEndian.PutUint32(buf, a[0])
  binary.BigEndian.PutUint32(buf[4:], a[1])
  binary.BigEndian.PutUint32(buf[8:], a[2])
  binary.BigEndian.PutUint32(buf[12:], a[3])
  return (&big.Int{}).SetBytes(buf)
}

func bigIntToUint128(a big.Int) Uint128 {
  bytes := a.Bytes()
  resbytes := make([]byte, 16)
  i := len(bytes) - 1;
  if i > 15 {
    i = 15;
  }
  for i := len(bytes) - 1; i >= 0; i-- {
    resbytes[i] = bytes[i]
  }
  result := Uint128{}
  result.H = binary.BigEndian.Uint64(bytes)
  result.L = binary.BigEndian.Uint64(bytes[8:])
  return result
}

func uint128ToBigInt(a Uint128) *big.Int {
  buf := make([]byte, 16)
  binary.BigEndian.PutUint64(buf, a.H)
  binary.BigEndian.PutUint64(buf[8:], a.L)
  return (&big.Int{}).SetBytes(buf)
}

func TestMult(t *testing.T) {
  rand.Seed(0)
  mod := big.NewInt(1)
  mod.Lsh(mod, 128)
  resbig, multbig := &big.Int{}, &big.Int{}
  for i := 0; i < 100000; i++ {
    a := []uint32{rand.Uint32(),rand.Uint32(),rand.Uint32(),rand.Uint32()}
    b := []uint32{rand.Uint32(),rand.Uint32(),rand.Uint32(),rand.Uint32()}
    au := makeUint128(a)
    bu := makeUint128(b)
    abig := makeBigInt(a)
    bbig := makeBigInt(b)
    resu := au.Mult(bu)
    resbig.Mod(multbig.Mul(abig, bbig), mod)
    //compare abig, resu 
    resubig := uint128ToBigInt(resu)
    if resbig.Cmp(resubig) != 0 {
      t.Errorf("Multiplied %x by %x (%x %x x %x %x), expected %x got %x (%x %x)",
        abig, bbig, au.H, au.L, bu.H, bu.L, resbig, resubig, resu.H, resu.L)
    }
  }
}
