/*
 * Author: fasion
 * Created time: 2026-02-08 17:41:50
 * Last Modified by: fasion
 * Last Modified time: 2026-02-08 18:49:08
 */

package ed25519x

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

// KeypairBytes creates an Ed25519 keypair from bytes.
func KeypairFromBytes(publicKeyBytes, privateKeyBytes []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return nil, nil, fmt.Errorf("invalid public key size: got %d, want %d", len(publicKeyBytes), ed25519.PublicKeySize)
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, nil, fmt.Errorf("invalid private key size: got %d, want %d", len(privateKeyBytes), ed25519.PrivateKeySize)
	}

	publicKey := ed25519.PublicKey(publicKeyBytes)
	privateKey := ed25519.PrivateKey(privateKeyBytes)
	if !publicKey.Equal(privateKey.Public()) {
		return nil, nil, fmt.Errorf("public key does not match private key")
	}

	return publicKey, privateKey, nil
}

func KeypairFromB64s(publicKeyB64s, privateKeyB64s string, encoding *base64.Encoding) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	if encoding == nil {
		encoding = base64.RawURLEncoding
	}

	publicKeyBytes, err := encoding.DecodeString(publicKeyB64s)
	if err != nil {
		return nil, nil, fmt.Errorf("decode public key bytes: %w", err)
	}

	privateKeyBytes, err := encoding.DecodeString(privateKeyB64s)
	if err != nil {
		return nil, nil, fmt.Errorf("decode private key bytes: %w", err)
	}

	return KeypairFromBytes(publicKeyBytes, privateKeyBytes)
}

func KeypairFromSeed(seed []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	privateKey := ed25519.NewKeyFromSeed(seed)
	return privateKey.Public().(ed25519.PublicKey), privateKey, nil
}

func KeypairFromSeedB64s(seedB64s string, encoding *base64.Encoding) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	if encoding == nil {
		encoding = base64.RawURLEncoding
	}

	seed, err := encoding.DecodeString(seedB64s)
	if err != nil {
		return nil, nil, fmt.Errorf("decode seed: %w", err)
	}

	return KeypairFromSeed(seed)
}
