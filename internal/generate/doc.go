// Package generate provides cryptographically secure secret generation
// utilities used during secret rotation.
//
// It supports configurable length, character set restrictions (alphanumeric
// only vs full charset including special characters), and optional base64
// encoding of the output. All randomness is sourced from crypto/rand to
// ensure suitability for production secret material.
//
// Basic usage:
//
//	g := generate.New(generate.Options{Length: 32})
//	secret, err := g.Secret()
//
// For raw random bytes encoded as base64:
//
//	encoded, err := generate.Bytes(32)
package generate
