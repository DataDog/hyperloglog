package hyperloglog

import (
	"testing"

	"github.com/OneOfOne/xxhash"
	humanize "github.com/dustin/go-humanize"
)

func int2bytes(i int) []byte {
	return []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
}

func uabs(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}

func TestStableHash(t *testing.T) {
	h := xxhash.Checksum64([]byte("hello world"))
	if h != 5020219685658847592 {
		t.Errorf("Expected 5020219685658847592, got %d", h)
	}
}

func TestHll64ComparisonSmall(t *testing.T) {
	hll32, _ := New(1 << 18)
	hll64, _ := New64(1 << 18)

	err32 := uint64(0)
	err64 := uint64(0)
	worse := uint64(0)
	same := uint64(0)
	better := uint64(0)
	for i := 0; i < 10000; i++ {
		hll32.Add(xxhash.Checksum32(int2bytes(i)))
		e32 := uabs(hll32.Count(), uint64(i+1))
		err32 += e32

		hll64.AddHash(xxhash.Checksum64(int2bytes(i)))
		e64 := uabs(hll64.Count(), uint64(i+1))
		err64 += e64

		if e32 < e64 {
			worse++
		} else if e32 == e64 {
			same++
		} else {
			better++
		}
	}

	if err32 < err64 {
		t.Errorf("32-bit error %s is less than 64-bit error %s", humanize.Comma(int64(err32)), humanize.Comma(int64(err64)))
	} else {
		t.Logf("32-bit error %s is greater than 64-bit error %s", humanize.Comma(int64(err32)), humanize.Comma(int64(err64)))
	}
	t.Logf("64-bit vs 32-bit:\nbetter:\t%d\nsame:\t%d\nworse:\t%d\n", better, same, worse)
}

func TestHll64ComparisonLarge(t *testing.T) {
	hll32, _ := New(1 << 18)
	hll64, _ := New64(1 << 18)

	for i := 0; i < 100000000; i++ {
		hll32.Add(xxhash.Checksum32(int2bytes(i)))
		hll64.AddHash(xxhash.Checksum64(int2bytes(i)))
	}

	err32 := uint64(0)
	err64 := uint64(0)
	worse := uint64(0)
	same := uint64(0)
	better := uint64(0)
	for i := 100000000; i < 100010000; i++ {
		hll32.Add(xxhash.Checksum32(int2bytes(i)))
		e32 := uabs(hll32.Count(), uint64(i+1))
		err32 += e32

		hll64.AddHash(xxhash.Checksum64(int2bytes(i)))
		e64 := uabs(hll64.Count(), uint64(i+1))
		err64 += e64

		if e32 < e64 {
			worse++
		} else if e32 == e64 {
			same++
		} else {
			better++
		}
	}

	if err32 < err64 {
		t.Errorf("32-bit error %s is less than 64-bit error %s", humanize.Comma(int64(err32)), humanize.Comma(int64(err64)))
	} else {
		t.Logf("32-bit error %s is greater than 64-bit error %s", humanize.Comma(int64(err32)), humanize.Comma(int64(err64)))
	}
	t.Logf("64-bit vs 32-bit:\nbetter:\t%d\nsame:\t%d\nworse:\t%d\n", better, same, worse)
}

func benchAdd32(b *testing.B, m uint) {
	hll32, _ := New(1 << m)
	for i := 0; i < b.N; i++ {
		hll32.Add(uint32(i))
	}
}

func benchAdd64(b *testing.B, m uint) {
	hll64, _ := New64(1 << m)
	for i := 0; i < b.N; i++ {
		hll64.AddHash(uint64(i))
	}
}

func BenchmarkHllAdd32_12(b *testing.B) {
	benchAdd32(b, 12)
}

func BenchmarkHllAdd64_12(b *testing.B) {
	benchAdd64(b, 12)
}

func BenchmarkHllAdd32_18(b *testing.B) {
	benchAdd32(b, 18)
}

func BenchmarkHllAdd64_18(b *testing.B) {
	benchAdd64(b, 18)
}

func benchCount32(b *testing.B, m uint) {
	hll32, _ := New(1 << m)
	for i := 0; i < 10000; i++ {
		hll32.Add(uint32(i))
	}
	for i := 0; i < b.N; i++ {
		hll32.Count()
	}
}

func benchCount64(b *testing.B, m uint) {
	hll64, _ := New64(1 << m)
	for i := 0; i < 10000; i++ {
		hll64.AddHash(uint64(i))
	}
	for i := 0; i < b.N; i++ {
		hll64.Count()
	}
}

func BenchmarkHllCount32_12(b *testing.B) {
	benchCount32(b, 12)
}

func BenchmarkHllCount64_12(b *testing.B) {
	benchCount64(b, 12)
}

func BenchmarkHllCount32_18(b *testing.B) {
	benchCount32(b, 18)
}

func BenchmarkHllCount64_18(b *testing.B) {
	benchCount64(b, 18)
}
