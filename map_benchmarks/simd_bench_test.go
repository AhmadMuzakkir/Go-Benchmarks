package map_benchmarks

import (
	"testing"
)

var list []int

func init() {
	size := 10_000_000
	list = make([]int, size)
	for i := 0; i < size; i++ {
		list[i] = i + 1
	}
}

func BenchmarkAdd_Simple(b *testing.B) {
	var count int
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(list); i++ {
			count += list[i]
		}
	}

	_ = count
}

func BenchmarkAdd_Multiple2(b *testing.B) {
	var count int
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(list); i += 2 {
			count += list[i] + list[i+1]
		}
	}

	_ = count
}

func BenchmarkAdd_Multiple4(b *testing.B) {
	var count int
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(list); i += 4 {
			count += list[i] + list[i+1] + list[2] + list[i+3]
		}
	}

	_ = count
}


func BenchmarkAdd_Multiple8(b *testing.B) {
	var count int
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(list); i += 8 {
			count += list[i] + list[i+1] + list[2] + list[i+3] + list[4] + list[i+5] + list[6] + list[i+7]
		}
	}

	_ = count
}
