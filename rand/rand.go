// Package for generating random data
package rand

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2015 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"math/rand"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var symbols = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890"

// ////////////////////////////////////////////////////////////////////////////////// //

// String return string with random chars
func String(length int) string {
	if length <= 0 {
		return ""
	}

	symbolsLength := len(symbols)
	result := make([]byte, length)

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < length; i++ {
		result[i] = symbols[rand.Intn(symbolsLength)]
	}

	return string(result)
}

// Int return random int
func Int(n int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(n)
}

// Float64 return random float64
func Float64(n int) float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Float64()
}

// Float32 return random float32
func Float32(n int) float32 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Float32()
}

// Slice return slice with random chars
func Slice(length int) []string {
	if length == 0 {
		return []string{}
	}

	symbolsLength := len(symbols)
	result := make([]string, length)

	for i := 0; i < length; i++ {
		result[i] = string(symbols[rand.Intn(symbolsLength)])
	}

	return result
}

// ////////////////////////////////////////////////////////////////////////////////// //