package utils

import (
	"math/rand"
	"time"
)

func GetRandomSliceElement(slice []string) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(slice))
	return slice[index]
}
