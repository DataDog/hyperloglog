package hyperloglog

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"math"
	"os"
	"testing"
)

// Return a dictionary up to n words. If n is zero, return the entire
// dictionary.
func dictionary(n int) []string {
	path := "/usr/share/dict/words"
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("can't open dictionary file '%s': %v\n", path, err)
		os.Exit(1)
	}
	buf := bufio.NewReader(f)

	count := 0
	words := make([]string, 0, n)
	for {
		if n != 0 && count < n {
			break
		}
		word, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		words = append(words, word)
		count++
	}

	f.Close()
	return words
}

func geterror(actual uint64, estimate uint64) (result float64) {
	return (float64(estimate) - float64(actual)) / float64(actual)
}

func testHyperLogLog(t *testing.T, n int, lowB, highB uint) {
	words := dictionary(n)
	bad := 0
	nWords := uint64(len(words))
	for i := lowB; i < highB; i++ {
		m := 1 << i

		h := New(m)
		hash := fnv.New32()
		for _, word := range words {
			hash.Write([]byte(word))
			h.Add(hash.Sum32())
			hash.Reset()
		}

		expectedError := 1.04 / math.Sqrt(float64(m))
		actualError := math.Abs(geterror(nWords, h.Count()))

		if actualError > expectedError {
			bad++
			t.Logf("m=%d: error=%.5f, expected <%.5f; actual=%d, estimated=%d\n",
				m, actualError, expectedError, nWords, h.Count())
		}

	}
	t.Logf("%d of %d tests exceeded estimated error", bad, highB-lowB)
}

func TestHyperLogLogSmall(t *testing.T) {
	testHyperLogLog(t, 5, 4, 17)
}

func TestHyperLogLogBig(t *testing.T) {
	testHyperLogLog(t, 0, 4, 17)
}

func benchmarkCount(b *testing.B, m int) {
	words := dictionary(0)

	h := New(m)
	hash := fnv.New32()
	for _, word := range words {
		hash.Write([]byte(word))
		h.Add(hash.Sum32())
		hash.Reset()
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		h.Count()
	}
}

func BenchmarkCount4(b *testing.B) {
	benchmarkCount(b, 16)
}

func BenchmarkCount5(b *testing.B) {
	benchmarkCount(b, 32)
}

func BenchmarkCount6(b *testing.B) {
	benchmarkCount(b, 64)
}

func BenchmarkCount7(b *testing.B) {
	benchmarkCount(b, 128)
}

func BenchmarkCount8(b *testing.B) {
	benchmarkCount(b, 256)
}

func BenchmarkCount9(b *testing.B) {
	benchmarkCount(b, 512)
}

func BenchmarkCount10(b *testing.B) {
	benchmarkCount(b, 1024)
}
