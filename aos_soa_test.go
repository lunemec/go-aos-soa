package main

import (
	"testing"

	"github.com/bjwbell/gensimd/simd"
)

var lSum = 10_000

type EntityAOS struct {
	X, Y, Z                                 int
	E1, E2, E3, E4, E5, E6, E7, E8, E9, E10 string
}

var AccX, AccY, AccZ int

func sumEntitiesInt(data []EntityAOS) (int, int, int) {
	var x, y, z int
	for _, entity := range data {
		x += entity.X
		y += entity.Y
		z += entity.Z
	}
	return x, y, z
}

func TestSumEntityInt(t *testing.T) {
	wantX, wantY, wantZ := 3, 6, 9
	entities := []EntityAOS{{X: 1, Y: 2, Z: 3}, {X: 1, Y: 2, Z: 3}, {X: 1, Y: 2, Z: 3}}
	gotX, gotY, gotZ := sumEntitiesInt(entities)
	if wantX != gotX {
		t.Errorf("X not correct, want: %d, got: %d", wantX, gotX)
	}
	if wantY != gotY {
		t.Errorf("Y not correct, want: %d, got: %d", wantY, gotY)
	}
	if wantZ != gotZ {
		t.Errorf("Z not correct, want: %d, got: %d", wantZ, gotZ)
	}
}

func BenchmarkAoSInt(b *testing.B) {
	var slice = make([]EntityAOS, lSum, lSum)
	var i int
	for i = 0; int(i) < lSum; i++ {
		slice[i] = EntityAOS{X: i, Y: i + 1, Z: i + 2}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AccX, AccY, AccZ = sumEntitiesInt(slice)
	}
}

type EntityUint16 struct {
	X, Y, Z uint16
}

func sumEntitiesUint16(data []EntityUint16) (uint16, uint16, uint16) {
	var x, y, z uint16
	for _, entity := range data {
		x += entity.X
		y += entity.Y
		z += entity.Z
	}
	return x, y, z
}

func TestSumEntityUint16(t *testing.T) {
	var (
		wantX uint16 = 3
		wantY uint16 = 6
		wantZ uint16 = 9
	)
	entities := []EntityUint16{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}
	gotX, gotY, gotZ := sumEntitiesUint16(entities)
	if wantX != gotX {
		t.Errorf("X not correct, want: %d, got: %d", wantX, gotX)
	}
	if wantY != gotY {
		t.Errorf("Y not correct, want: %d, got: %d", wantY, gotY)
	}
	if wantZ != gotZ {
		t.Errorf("Z not correct, want: %d, got: %d", wantZ, gotZ)
	}
}

var AccXUint16, AccYUint16, AccZUint16 uint16

func enchmarkAoSUint16(b *testing.B) {
	var slice = make([]EntityUint16, lSum, lSum)
	var i uint16
	for i = 0; int(i) < lSum; i++ {
		slice[i] = EntityUint16{X: i, Y: i + 1, Z: i + 2}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AccXUint16, AccXUint16, AccXUint16 = sumEntitiesUint16(slice)
	}
}

type EntitiesSOA struct {
	X                                       []int
	Y                                       []int
	Z                                       []int
	E1, E2, E3, E4, E5, E6, E7, E8, E9, E10 []string
}

func sumEntitiesSoA(data EntitiesSOA) (int, int, int) {
	var x, y, z int
	for _, vx := range data.X {
		x += vx
	}
	for _, vy := range data.Y {
		y += vy
	}
	for _, vz := range data.Z {
		z += vz
	}
	return x, y, z
}

func TestSoAEntities(t *testing.T) {
	var (
		wantX int = 3
		wantY int = 6
		wantZ int = 9
	)
	entities := EntitiesSOA{X: []int{1, 1, 1}, Y: []int{2, 2, 2}, Z: []int{3, 3, 3}}
	gotX, gotY, gotZ := sumEntitiesSoA(entities)
	if wantX != gotX {
		t.Errorf("X not correct, want: %d, got: %d", wantX, gotX)
	}
	if wantY != gotY {
		t.Errorf("Y not correct, want: %d, got: %d", wantY, gotY)
	}
	if wantZ != gotZ {
		t.Errorf("Z not correct, want: %d, got: %d", wantZ, gotZ)
	}
}

func BenchmarkSoA(b *testing.B) {
	var entities = EntitiesSOA{
		X: make([]int, lSum, lSum),
		Y: make([]int, lSum, lSum),
		Z: make([]int, lSum, lSum),
	}
	for i := 0; i < lSum; i++ {
		entities.X[i] = i
		entities.Y[i] = i + 1
		entities.Z[i] = i + 2
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AccX, AccY, AccZ = sumEntitiesSoA(entities)
	}
}

type EntitiesU64 struct {
	X []uint64
	Y []uint64
	Z []uint64
}

func sumSimd(data []uint64) uint64 {
	var (
		idx int
		acc uint64
	)
	for idx = 0; idx+4 <= len(data); idx += 4 {
		x1 := simd.U64x2{data[idx], data[idx+1]}
		x2 := simd.U64x2{data[idx+2], data[idx+3]}
		out := simd.AddU64x2(x1, x2)
		acc += out[0] + out[1]
	}
	for _, v := range data[idx:len(data)] {
		acc += v
	}
	return acc
}

func TestSumSimd(t *testing.T) {
	var want uint64 = 36
	sum := sumSimd([]uint64{1, 2, 3, 4, 5, 6, 7, 8})
	if sum != want {
		t.Errorf("sumSimd is not what is expected, want: %d, got: %d", want, sum)
	}

	want = 21
	sum = sumSimd([]uint64{1, 2, 3, 4, 5, 6})
	if sum != want {
		t.Errorf("sumSimd is not what is expected, want: %d, got: %d", want, sum)
	}
}

var (
	AccXSimd uint64
	AccYSimd uint64
	AccZSimd uint64
)

func BenchmarkSimd(b *testing.B) {
	var entities = EntitiesU64{
		X: make([]uint64, lSum, lSum),
		Y: make([]uint64, lSum, lSum),
		Z: make([]uint64, lSum, lSum),
	}
	var i uint64
	for i = 0; i < uint64(lSum); i++ {
		entities.X[i] = i
		entities.Y[i] = i + 1
		entities.Z[i] = i + 2
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AccXSimd = sumSimd(entities.X)
		AccYSimd = sumSimd(entities.Y)
		AccZSimd = sumSimd(entities.Z)
	}
}
