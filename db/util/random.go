package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandString(n int) string {

	randRunes := make([]rune, n)

	for i, _ := range randRunes {
		randRunes[i] = rune(RandInt('a', 'z'))
	}

	return string(randRunes)
}

func RandCurrency() string {
	currencies := []string{
		"$$$",
		"€€€",
		"₪₪₪",
		"₴₴₴",
		"₿₿₿",
	}

	return currencies[rand.Intn(len(currencies))]
}
