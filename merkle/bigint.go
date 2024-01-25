package merkle

import (
	"fmt"
	"math/big"
	"strings"
)

type BigInt struct {
	big.Int
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}

	base := 10
	s := string(p)

	// trim quotes
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		s = s[1 : len(s)-1]
	}

	// if it starts with 0x, decode using base 16
	if strings.HasPrefix(s, "0x") {
		base = 16

		s = strings.TrimPrefix(s, "0x")
	}

	var z big.Int
	_, ok := z.SetString(s, base)

	if !ok {
		return fmt.Errorf("not a valid big integer: %s, base=%d, hex=%v", s, base, strings.HasPrefix(s, "0x"))
	}

	b.Int = z
	return nil
}
