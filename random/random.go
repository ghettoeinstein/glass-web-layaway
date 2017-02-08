package random

import (
	"github.com/google/uuid"
	"math/rand"
)

func GenerateUUID() string {
	return uuid.Must(uuid.NewRandom()).String()

}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomPin(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(randInt(48, 57))
	}
	return string(bytes)
}
