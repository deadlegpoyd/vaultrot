package generate

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

const (
	DefaultLength  = 32
	charsetAlpha   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetNumeric = "0123456789"
	charsetSpecial = "!@#$%^&*()-_=+[]{}"
	charsetAll     = charsetAlpha + charsetNumeric + charsetSpecial
)

// Options configures secret generation.
type Options struct {
	Length      int
	AlphaOnly   bool
	Base64Encode bool
}

// Generator produces new secret values.
type Generator struct {
	opts Options
}

// New returns a Generator with the given options.
func New(opts Options) *Generator {
	if opts.Length <= 0 {
		opts.Length = DefaultLength
	}
	return &Generator{opts: opts}
}

// Secret generates a cryptographically random secret string.
func (g *Generator) Secret() (string, error) {
	charset := charsetAll
	if g.opts.AlphaOnly {
		charset = charsetAlpha + charsetNumeric
	}

	buf := make([]byte, g.opts.Length)
	for i := range buf {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("generate: random read failed: %w", err)
		}
		buf[i] = charset[idx.Int64()]
	}

	result := string(buf)
	if g.opts.Base64Encode {
		result = base64.StdEncoding.EncodeToString(buf)
	}
	return result, nil
}

// Bytes generates n cryptographically random bytes and returns them base64-encoded.
func Bytes(n int) (string, error) {
	if n <= 0 {
		n = DefaultLength
	}
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate: rand.Read failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}
