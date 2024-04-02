package random

import (
	"math/rand"
	"time"
)

func NewRandomString(aliasLength int) string {
	rand.Seed(time.Now().UnixNano())
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var str string

	for i := 0; i < aliasLength; i++ {
		str += string(chars[rand.Intn(len(chars))])
	}

	return str
}
