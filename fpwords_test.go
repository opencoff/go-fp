package fp

import (
	"bytes"
	crand "crypto/rand"
	"encoding/hex"
	mrand "math/rand/v2"
	"testing"
)

var proQuintTests = []struct {
	b string
	s string
}{
	{"7f000001", "lusab-babad"},
	{"3f54dcc1", "gutih-tugad"},
	{"3f760723", "gutuk-bisog"},
	{"8c62c18d", "mudof-sakat"},
}

// newAsserter returns a func that calls t.Fatalf on failure.
func newAsserter(t *testing.T) func(bool, string, ...any) {
	return func(ok bool, msg string, args ...any) {
		t.Helper()
		if !ok {
			t.Fatalf(msg, args...)
		}
	}
}

func TestEncodeWords(t *testing.T) {
	assert := newAsserter(t)

	s2b := func(s string) []byte {
		b, err := hex.DecodeString(s)
		assert(err == nil, "invalid hex '%s': %s", s, err)
		return b
	}

	for i := range proQuintTests {
		z := &proQuintTests[i]
		s := ToWords(s2b(z.b))
		assert(s == z.s, "exp %s: saw %s", z.s, s)
	}

	// now lets decode em
	for i := range proQuintTests {
		z := &proQuintTests[i]
		b, err := FromWords(z.s)
		assert(err == nil, "decode %s: %s", z.s, err)

		bx := hex.EncodeToString(b)
		assert(z.b == bx, "exp %s: saw %s", z.b, bx)
	}
}

func TestEncodeDecode(t *testing.T) {
	assert := newAsserter(t)

	buf := make([]byte, 256)
	for range 100 {
		sz := mrand.IntN(8) + mrand.IntN(128)
		if (sz & 1) > 0 {
			sz++
		}
		b := buf[:sz]
		randbuf(b)

		s := ToWords(b)
		x, err := FromWords(s)
		assert(err == nil, "%d: %s", sz, err)
		assert(bytes.Equal(b, x), "%d: bytes uneq", sz)
	}
}

func randbuf(b []byte) []byte {
	_, err := crand.Read(b)
	if err != nil {
		panic(err.Error())
	}
	return b
}
