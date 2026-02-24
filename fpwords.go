// fpwords.go -- procedurally convert a byte array into a readable string
//
// Uses the Proquint encoding method

package fp

import (
	"errors"
	"strings"
	"sync"
)

var (
	_Consonants = []byte("bdfghjklmnprstvz")
	_Vowels     = []byte("aiou")

	// Simple lookup table for character indices
	_CharMap  [256]byte
	_InitOnce sync.Once
)

func initLookupTable() {
	// Set consonant indices - all other chars default to zero
	m := _CharMap[:]
	for i, c := range _Consonants {
		m[c] = byte(i)
	}

	// Set vowel indices; we can reuse the same map because Vowels and Consonants
	// are distinct - but can share the same index
	for i, v := range _Vowels {
		m[v] = byte(i)
	}
	// All other characters default to zero, which is a no-op in bit operations
}

// Encode converts binary data to a proquint string
// Note that ProQuint encoding always works in 16-bit chunks; so, odd length
// inputs are extended with a zero-byte.
func ToWords(b []byte) string {
	// Pad b to 16-bit boundary if needed
	if len(b)%2 != 0 {
		b = append(b, 0)
	}

	// Calculate exact number of words
	nw := len(b) / 2
	words := make([]string, 0, nw)

	c := _Consonants
	v := _Vowels
	for len(b) > 0 {
		w := (uint16(b[0]) << 8) | uint16(b[1])

		x := [5]byte{
			c[(w>>12)&0xF], // First consonant (4 bits)
			v[(w>>10)&0x3], // First vowel (2 bits)
			c[(w>>6)&0xF],  // Second consonant (4 bits)
			v[(w>>4)&0x3],  // Second vowel (2 bits)
			c[w&0xF],       // Third consonant (4 bits)
		}

		words = append(words, string(x[:]))
		b = b[2:]
	}

	return strings.Join(words, "-")
}

// Decode converts a proquint string back to binary data
// Note that ProQuint encoding always returns even length output.
func FromWords(s string) ([]byte, error) {
	// Ensure lookup table is initialized
	_InitOnce.Do(initLookupTable)

	b := []byte(s)
	out := make([]byte, 0, 2*len(b)/5+1)

	m := _CharMap
	for len(b) >= 5 {
		if b[0] == '-' {
			b = b[1:]
			continue
		}

		c1 := m[b[0]]
		v1 := m[b[1]]
		c2 := m[b[2]]
		v2 := m[b[3]]
		c3 := m[b[4]]

		// Assemble the 16-bit word
		w := uint16(c1)<<12 | uint16(v1)<<10 | uint16(c2)<<6 | uint16(v2)<<4 | uint16(c3)

		out = append(out, byte(w>>8))
		out = append(out, byte(w&0xff))
		b = b[5:]
	}

	if len(b) > 0 {
		return nil, ErrIncomplete
	}

	return out, nil
}

var (
	ErrIncomplete = errors.New("invalid proquint format: incomplete quintet")
)
