package main

import (
	"crypto/rand"
	"math/big"
)

// Generates a random user agent
func randomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:134.0) Gecko/20100101 Firefox/134.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_7_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15",
	}

	max := big.NewInt(int64(len(userAgents)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return userAgents[0]
	}
	return userAgents[n.Int64()]
}

// Generates a random session ID
func generateSessionID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result[i] = charset[index.Int64()]
	}

	return string(result)
}
