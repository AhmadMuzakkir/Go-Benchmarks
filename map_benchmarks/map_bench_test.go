package map_benchmarks

import (
	"encoding/binary"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/cespare/xxhash"
)

// https://github.com/golang/go/issues/9477
// https://github.com/dgraph-io/benchmarks/blob/master/cachebench/cache_bench_test.go

//func BenchmarkMap_KeyString_ValString(b *testing.B) {
//	n := int(30e6)
//	m := make(map[string]string, n)
//
//	for i := 0; i < n; i++ {
//		is := strconv.Itoa(i)
//		k := "key_ " + is
//		v := "val_ " + is
//
//		m[k] = v
//	}
//
//	runtime.GC()
//
//	b.ReportAllocs()
//	b.ResetTimer()
//
//	for n := 0; n < b.N; n++ {
//		for i := 0; i < 10; i++ {
//			runtime.GC()
//		}
//	}
//
//	_ = m["key_0"]
//}

var ki interface{}

func TestKeyUint64_ValString_Interface(t *testing.T) {
	n := int(10e6)
	m := make(map[uint64]interface{}, n)

	for i := 0; i < n; i++ {
		is := strconv.Itoa(i)
		k := xxhash.Sum64([]byte("k_" + is))
		v := "val_ " + is

		m[k] = v
	}

	runtime.GC()

	t.Logf("GC took: %s", timeGC())

	ki = m[0]
}

var ks interface{}
func TestKeyUint64_ValString_Slice(t *testing.T) {
	n := int(10e6)
	m := make(map[uint64][]byte, n)

	for i := 0; i < n; i++ {
		is := strconv.Itoa(i)
		k := xxhash.Sum64([]byte("k_" + is))
		v := "val_ " + is

		m[k] = []byte(v)
	}

	runtime.GC()

	t.Logf("GC took: %s", timeGC())

	ks = m[0]
}

func TestKeyUint64_ValString_Shard(t *testing.T) {
	n := int(10e6)
	m := make([]map[uint64]interface{}, 100)

	for i := range m {
		m[i] = make(map[uint64]interface{})
	}

	for i := 0; i < n; i++ {
		is := strconv.Itoa(i)
		k := xxhash.Sum64([]byte("k_" + is))
		v := "val_ " + is

		shard := k % 100
		m[shard][k] = v
	}

	runtime.GC()

	t.Logf("GC took: %s", timeGC())

	_ = m[0][0]

	//for i := range m {
	//	t.Log(len(m[i]))
	//}
}

func timeGC() time.Duration {
	start := time.Now()
	runtime.GC()
	return time.Since(start)
}

func BenchmarkMap_KeyArray_ValString(b *testing.B) {
	n := int(30e6)
	m := make(map[[16]byte]string, n)

	var k [16]byte
	for i := 0; i < n; i++ {
		is := strconv.Itoa(i)
		v := "val_ " + is

		binary.LittleEndian.PutUint16(k[:], uint16(i))

		m[k] = v
	}

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	var l []time.Duration

	for n := 0; n < b.N; n++ {
		for i := 0; i < 10; i++ {
			s := time.Now()
			runtime.GC()
			l = append(l, time.Since(s))
		}
	}

	_ = m[k]

	b.Log(l)
}

func BenchmarkMap_KeyArray_ValSlice(b *testing.B) {
	n := int(10e6)
	m := make(map[[16]byte][]byte, n)

	var k [16]byte
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint16(k[:], uint16(i))

		m[k] = []byte("val_ " + strconv.Itoa(i))
	}

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < 10; i++ {
			runtime.GC()
		}
	}

	_ = m[k]
}


func BenchmarkMap_KeyUint64_ValString(b *testing.B) {
	n := int(30e6)
	m := make(map[uint64]string, n)

	for i := 0; i < n; i++ {
		is := strconv.Itoa(i)
		v := "val_ " + is

		k := uint64(i)

		m[k] = v
	}

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	var l []time.Duration

	for n := 0; n < b.N; n++ {
		for i := 0; i < 10; i++ {
			s := time.Now()
			runtime.GC()
			l = append(l, time.Since(s))
		}
	}

	_ = m[0]

	b.Log(l)
}

func BenchmarkMap_KeyUint64_ValSlice(b *testing.B) {
	n := int(30e6)
	m := make(map[uint64][]byte, n)

	for i := 0; i < n; i++ {
		k := uint64(i)

		m[k] = []byte("val_ " + strconv.Itoa(i))
	}

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	var l []time.Duration

	for n := 0; n < b.N; n++ {
		for i := 0; i < 10; i++ {
			s := time.Now()
			runtime.GC()
			l = append(l, time.Since(s))
		}
	}

	_ = m[0]

	b.Log(l)
}