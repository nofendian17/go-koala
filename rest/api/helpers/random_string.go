package helpers

import "math/rand"

var pool = "abcdefghijklmnopqrstuvwxyzABCEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(l int) string {
	bytes := make([]byte, l)

	for i := 0; i < l; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}

	return string(bytes)
}

// main
// rand.Seed(time.Now().UnixNano())
// fmt.Println(randomString(12))
