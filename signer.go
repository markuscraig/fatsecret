package fatsecret

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// Signer signs Oauth1 messages
type Signer interface {
	Name() string
	Sign(key string, msg string) string
}

// HMACSigner that SHA-1 signs a message
type HMACSigner struct {
	ConsumerSecret string
}

// NewHMACSigner creates and returns an HMACSigner instance
func NewHMACSigner(secret string) *HMACSigner {
	return &HMACSigner{
		ConsumerSecret: secret,
	}
}

// Name returns the signer's Oauth1 algorithm name
func (s *HMACSigner) Name() string {
	return "HMAC-SHA1"
}

// Sign calculates the HMAC digest and returns the base64 string
func (s *HMACSigner) Sign(tokenSecret string, msg string) string {
	// generate the oauth sha1 base64 signature
	key := s.ConsumerSecret + "&" + tokenSecret
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(msg))
	sigBytes := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(sigBytes)
}
