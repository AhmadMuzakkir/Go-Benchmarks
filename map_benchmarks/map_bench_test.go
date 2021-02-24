package map_benchmarks

import (
	"encoding/binary"
	"runtime"
	"strconv"
	"testing"

	"github.com/cespare/xxhash"
)

const n = 100_000_000
const gcCount = 1

// https://github.com/golang/go/issues/9477
// https://github.com/dgraph-io/benchmarks/blob/master/cachebench/cache_bench_test.go

func BenchmarkMapGC_KeyString(b *testing.B) {
	m := make(map[string][]byte, n)

	for i := 0; i < n; i++ {
		is := strconv.Itoa(i)

		k := "key_" + is

		v := []byte("val_ " + is)

		m[k] = v
	}

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < gcCount; i++ {
			runtime.GC()
		}
	}

	for k, v := range m {
		_, _ = k, v
		break
	}
}

func BenchmarkMapGC_KeyArray(b *testing.B) {
	m := make(map[[16]byte][]byte, n)

	for i := 0; i < n; i++ {
		is := strconv.Itoa(i)

		var k [16]byte
		binary.LittleEndian.PutUint64(k[:], xxhash.Sum64String("key_"+is))

		v := []byte("val_ " + strconv.Itoa(i))

		m[k] = v
	}

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < gcCount; i++ {
			runtime.GC()
		}
	}

	for k, v := range m {
		_, _ = k, v
		break
	}
}

func BenchmarkMapGC_KeyUint64(b *testing.B) {
	m := make(map[uint64][]byte, n)

	for i := 0; i < n; i++ {
		is := strconv.Itoa(i)

		k := xxhash.Sum64String("key_" + is)

		v := []byte("val_ " + is)

		m[k] = v
	}

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < gcCount; i++ {
			runtime.GC()
		}
	}

	for k, v := range m {
		_, _ = k, v
		break
	}
}
