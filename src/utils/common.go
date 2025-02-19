package utils

import (
	"math/rand"
	"time"
)

// RandomSleep sleeps for a random duration between min and max seconds
func RandomSleep(min, max int) {
	if min > max {
		min, max = max, min // Swap if min is greater than max
	}
	if min == max {
		time.Sleep(time.Duration(min) * time.Second)
		return
	}

	// Create a new random source each time for true randomness
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seconds := r.Intn(max-min+1) + min
	time.Sleep(time.Duration(seconds) * time.Second)
}

// Sleep sleeps for the exact number of seconds
func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func Contains(options []string, target string) bool {
	for _, option := range options {
		if option == target {
			return true
		}
	}
	return false
}