package hash

import (
	"fmt"
)

// New will return the Adler32 hash of whatever you pass into it.
// Beware that this doesn't support rolling, so if the input is too big,
// the uint32 will overflow and the hash will be incorrect.
// Returns the base16 representation of the hash as a string.
func New(s string) string {
	bytes := []byte(s)

	a := uint32(1)
	b := uint32(0)
	mod := uint32(65521)

	for _, by := range bytes {
		a += uint32(by) % mod
		b += a % mod
	}

	return fmt.Sprintf("%x", (b<<16)|a)
}
