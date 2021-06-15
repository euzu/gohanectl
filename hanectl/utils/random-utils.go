package utils

import (
	"math/rand"
	"time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!$%&/()=?{[]}+*#'-_.;:|"

func RandomString(length int) string {
	b := make([]byte, length)
	clen := len(charset)
	for i := range b {
		b[i] = charset[seededRand.Intn(clen)]
	}
	return string(b)
}
