# Golang AOS vs SOA vs SIMD

## Iteration
First of, I wanted to see just by how much is iterating over hashmap slower than slice,
and it turns out its quite a lot:
```
BenchmarkHashMapIter-16       	    9862	    120403 ns/op
BenchmarkSliceRangeIter-16    	  196664	      6241 ns/op
BenchmarkSliceIndexIter-16    	  163333	      7561 ns/op
```
I was also surprised that iterating by index access is slower, so I tried to change
the number of items in the slice from `10_000` to `1_000_000`:
```
BenchmarkHashMapIter-16       	      88	  13626294 ns/op
BenchmarkSliceRangeIter-16    	    1506	    817917 ns/op
BenchmarkSliceIndexIter-16    	    1321	    877711 ns/op
```
It did not help. Ideas?

Another set for iteration is to test if re-slicing N elements is faster than just iterating N times and assigning by index. I thought that reslicing is slower, but it is way, way faster:
```
BenchmarkSliceReslicing-4   	1000000000	         0.989 ns/op
BenchmarkSliceTake3-4       	234024609	         5.09 ns/op
```

## AOS vs SOA vs SIMD
Then I wanted to try different data layout in memory, the theory is, that `SOA` (structure of arrays) should be faster than `AOS` (array of structs), because the data are loaded from main memory to L3/L2/L1 cache in contiguous blocks.
So having all the items together in one array should make it more friendly for the CPU cache.

At first, I thought that this is correct, but it turned out I got incorrect results by having `-race` on when doing benchmarks! Thanks VSCode!

After using regular: `go test -bench ".*"`, I got these (quite surprising) results:
```go
// Array of structs.
type EntityAOS struct {
	X, Y, Z int
}
type EntitiesAOS []EntityAOS

// Struct of arrays.
type EntitiesSOA struct {
	X []int
	Y []int
	Z []int
}
```
```
BenchmarkAoSInt-16            	  229623	      5321 ns/op
BenchmarkAoSUint16-16         	  243499	      5038 ns/op
BenchmarkSoA-16               	  128618	      9456 ns/op
BenchmarkSimd-16              	    7244	    166975 ns/op
```

So I thought that `AoS` is faster, but then I tried adding more "junk" data, that is not
relevant to the computation, and will cause more fragmentation in memory. Voilà!
```go
// Array of structs.
type EntityAOS struct {
    X, Y, Z int
    E1, E2, E3, E4, E5, E6, E7, E8, E9, E10 string
}
type EntitiesAOS []EntityAOS

// Struct of arrays.
type EntitiesSOA struct {
	X []int
	Y []int
    Z []int
    E1, E2, E3, E4, E5, E6, E7, E8, E9, E10 []string
}
```
```
BenchmarkAoSInt-16            	   12324	    100404 ns/op
BenchmarkSoA-16               	  150512	      7913 ns/op
BenchmarkSimd-16              	    7964	    152958 ns/op
```
Now we can see how `SoA` is supposed to work.

I also wanted to see how easy it is to use `SIMD` to process the `SoA` data, but its even slower than
`AoS`. I'm not yet sure why.
