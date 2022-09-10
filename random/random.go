package random

import (
	"math/rand"
	"time"
)

const (
	letterIdxBits = 4                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GenAlphaNumSequence(length int, seed int64) string {
	src := rand.NewSource(time.Now().UnixNano())
	if seed != 0 {
		src = rand.NewSource(seed)
	}

	letterBytes := []byte("AaBbCcDdEeFfGgHhIiJjKkLlMm0123456789NnOoPpQqRrSsTtUuVvWwXxYyZz")

	rand.Shuffle(len(letterBytes), func(i, j int) {
		letterBytes[i], letterBytes[j] = letterBytes[j], letterBytes[i]
	})
	b := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
