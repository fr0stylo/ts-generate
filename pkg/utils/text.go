package utils

import "math/rand"

var lowerRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var upperRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	// rand.Seed(time.Now().UnixNano())
	rand.Seed(1)
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		if i > 0 {
			b[i] = lowerRunes[rand.Intn(len(lowerRunes))]
		} else {
			b[i] = upperRunes[rand.Intn(len(upperRunes))]
		}
	}
	return string(b)
}
