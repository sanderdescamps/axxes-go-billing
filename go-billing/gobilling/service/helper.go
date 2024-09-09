package service

import "crypto/rand"

const RAND_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-"

func RandomString(length int) string {
	ll := len(RAND_CHARS)
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = RAND_CHARS[int(b[i])%ll]
	}
	return string(b)
}
