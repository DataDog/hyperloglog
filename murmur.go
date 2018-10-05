package hyperloglog

import (
	"encoding/binary"
	"reflect"
	"runtime"
	"unsafe"
)

const uint32Size = unsafe.Sizeof(uint32(0))

func advance(sh *reflect.StringHeader) uint32 {
	b := *(*[]byte)(unsafe.Pointer(sh))
	k := binary.LittleEndian.Uint32(b)
	sh.Len -= int(uint32Size)
	sh.Data += uint32Size
	return k
}

// This file implements the murmur3 32-bit hash on 32bit and 64bit integers
// for little endian machines only with no heap allocation.  If you are using
// HLL to count integer IDs on intel machines, this is your huckleberry.

// MurmurString implements a fast version of the murmur hash function for strings
// for little endian machines.  Suitable for adding strings to HLL counter.
func MurmurString(key string) uint32 {
	var c1, c2 uint32 = 0xcc9e2d51, 0x1b873593
	var h, k uint32

	// Reinterpret the string as a `StringHeader`. This comes with three important caveats:
	// 1. We must never write through the pointer derived. Golang strings are immutable and we cannot
	//    break that assumption.
	// 2. Golang continues to have a non-moving GC. This only works because the Golang GC is
	//    (currently) non-moving. There are no plans to break this yet, but it remains a caveat.
	// 3. `key` is used after the `StringHeader` is no longer needed. Currently, `runtime.KeepAlive`
	//    is used as a no-op use.
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&key))
	blen := strHeader.Len

	if strHeader.Len >= int(uint32Size) {
		// for each 4 byte chunk of `key'

		for {
			k := advance(strHeader)
			// encode next 4 byte chunk of `key'
			k *= c1
			k = (k << 15) | (k >> (32 - 15))
			k *= c2
			h ^= k
			h = (h << 13) | (h >> (32 - 13))
			h = (h * 5) + 0xe6546b64
			if strHeader.Len < int(uint32Size) {
				break
			}
		}
	}

	k = 0
	// remainder
	tail := *(*[]byte)(unsafe.Pointer(strHeader))
	switch len(tail) {
	case 3:
		k ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k ^= uint32(tail[0])
		k *= c1
		k = (k << 15) | (k >> (32 - 15))
		k *= c2
		h ^= k
	}

	h ^= uint32(blen)
	h ^= (h >> 16)
	h *= 0x85ebca6b
	h ^= (h >> 13)
	h *= 0xc2b2ae35
	h ^= (h >> 16)

	runtime.KeepAlive(&key)

	return h
}

// Murmur32 implements a fast version of the murmur hash function for uint32 for
// little endian machines.  Suitable for adding 32bit integers to a HLL counter.
func Murmur32(i uint32) uint32 {
	var c1, c2 uint32 = 0xcc9e2d51, 0x1b873593
	var h, k uint32
	k = i
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// second part
	h ^= 4
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

// Murmur64 implements a fast version of the murmur hash function for uint64 for
// little endian machines.  Suitable for adding 64bit integers to a HLL counter.
func Murmur64(i uint64) uint32 {
	var c1, c2 uint32 = 0xcc9e2d51, 0x1b873593
	var h, k uint32
	//first 4-byte chunk
	k = uint32(i)
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// second 4-byte chunk
	k = uint32(i >> 32)
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// second part
	h ^= 8
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

// Murmur128 implements a fast version of the murmur hash function for two uint64s
// for little endian machines.  Suitable for adding a 128bit value to an HLL counter.
func Murmur128(i, j uint64) uint32 {
	var c1, c2 uint32 = 0xcc9e2d51, 0x1b873593
	var h, k uint32
	//first 4-byte chunk
	k = uint32(i)
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// second 4-byte chunk
	k = uint32(i >> 32)
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// third 4-byte chunk
	k = uint32(j)
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// fourth 4-byte chunk
	k = uint32(j >> 32)
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// second part
	h ^= 16
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h

}
