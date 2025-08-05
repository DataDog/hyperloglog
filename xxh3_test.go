package hyperloglog

import (
	"encoding/binary"
	"math/rand/v2"
	"testing"
	"unsafe"

	"github.com/dustin/randbo"
	"github.com/zeebo/xxh3"
)

// Benchmarks
func benchmarkXxh3_Bytes(b *testing.B, input [][]byte) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, x := range input {
			_ = uint32(xxh3.Hash(x))
		}
	}
}

func benchmarkXxh3_64(b *testing.B, input []uint64) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, x := range input {
			buf := make([]byte, 8, 8)
			binary.LittleEndian.PutUint64(buf, x)
			_ = uint32(xxh3.Hash(buf))
		}
	}
}

func benchmarkXxh3_String(b *testing.B, input []string) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, x := range input {
			_ = uint32(xxh3.HashString(x))
		}
	}
}

func benchmarkXxh3_Hash32(b *testing.B, input []string) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, x := range input {
			b := *(*[]byte)(unsafe.Pointer(&x))
			_ = uint32(xxh3.Hash(b))
		}
	}
}

func Benchmark100MXxh3_Bytes(b *testing.B) {
	input := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		x := make([]byte, 1000)
		RandFill(x)
		input[i] = x
	}
	benchmarkXxh3_Bytes(b, input)
}

func Benchmark100Xxh3_64(b *testing.B) {
	input := make([]uint64, 100)
	for i := 0; i < 100; i++ {
		input[i] = rand.Uint64N(1 << 63)
	}
	benchmarkXxh3_64(b, input)
}

func Benchmark100Xxh3_String(b *testing.B) {
	input := make([]string, 100)
	for i := 0; i < 100; i++ {
		input[i] = randString((i % 15) + 5)
	}
	benchmarkXxh3_String(b, input)
}

func Benchmark100Xxh3_Hash32(b *testing.B) {
	input := make([]string, 100)
	for i := 0; i < 100; i++ {
		input[i] = randString((i % 15) + 5)
	}
	benchmarkXxh3_Hash32(b, input)
}

func BenchmarkXxh3_StringBig(b *testing.B) {
	// Make a 100Mb string and use that as a benchmark
	r := randbo.New()
	slice := make([]byte, 100*1024*1024)
	_, err := r.Read(slice)
	if err != nil {
		b.Fatalf("Failed to create benchmark data: %s", err)
	}
	s := string(slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uint32(xxh3.HashString(s))
	}
}
