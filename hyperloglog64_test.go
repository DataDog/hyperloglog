package hyperloglog

import (
	"testing"

	humanize "github.com/dustin/go-humanize"
	"github.com/spaolacci/murmur3"
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

func TestHll64ComparisonSmall64(t *testing.T) {
	hll32, _ := New(1 << 12)
	hll64, _ := New64(1 << 12)

	err32 := uint64(0)
	err64 := uint64(0)
	worse := uint64(0)
	same := uint64(0)
	better := uint64(0)
	for i := 0; i < 1000000; i++ {
		hll32.Add(murmur3.Sum32(int2bytes(i)))
		e32 := uabs(hll32.Count(), uint64(i+1))
		err32 += e32

		hll64.AddHash(murmur3.Sum64(int2bytes(i)))
		e64 := uabs(hll64.Count64(), uint64(i+1))
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

func TestHll64ComparisonSmall(t *testing.T) {
	hll32, _ := New(1 << 12)
	hll64, _ := New64(1 << 12)

	err32 := uint64(0)
	err64 := uint64(0)
	worse := uint64(0)
	same := uint64(0)
	better := uint64(0)
	for i := 0; i < 1000000; i++ {
		hll32.Add(murmur3.Sum32(int2bytes(i)))
		e32 := uabs(hll32.Count(), uint64(i+1))
		err32 += e32

		hll64.AddHash(murmur3.Sum64(int2bytes(i)))
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
		hll32.Add(murmur3.Sum32(int2bytes(i)))
		hll64.AddHash(murmur3.Sum64(int2bytes(i)))
	}

	err32 := uint64(0)
	err64 := uint64(0)
	worse := uint64(0)
	same := uint64(0)
	better := uint64(0)
	for i := 100000000; i < 100010000; i++ {
		hll32.Add(murmur3.Sum32(int2bytes(i)))
		e32 := uabs(hll32.Count(), uint64(i+1))
		err32 += e32

		hll64.AddHash(murmur3.Sum64(int2bytes(i)))
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

func TestHll64ComparisonLarge64(t *testing.T) {
	hll32, _ := New(1 << 18)
	hll64, _ := New64(1 << 18)

	for i := 0; i < 100000000; i++ {
		hll32.Add(murmur3.Sum32(int2bytes(i)))
		hll64.AddHash(murmur3.Sum64(int2bytes(i)))
	}

	err32 := uint64(0)
	err64 := uint64(0)
	worse := uint64(0)
	same := uint64(0)
	better := uint64(0)
	for i := 100000000; i < 100010000; i++ {
		hll32.Add(murmur3.Sum32(int2bytes(i)))
		e32 := uabs(hll32.Count(), uint64(i+1))
		err32 += e32

		hll64.AddHash(murmur3.Sum64(int2bytes(i)))
		e64 := uabs(hll64.Count64(), uint64(i+1))
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
