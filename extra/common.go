package extra

import (
	http "github.com/Danny-Dasilva/fhttp"
	"math/rand"
	"strings"
	"time"
)

func RandomSleep(min, max int) {
	if min > max {
		min, max = max, min
	}
	duration := rand.Intn(max-min+1) + min
	time.Sleep(time.Duration(duration) * time.Second)
}

func ChangeProxyURL(link string) bool {
	for i := 0; i < 3; i++ {
		c := http.Client{}
		req, err := http.NewRequest("GET", link, strings.NewReader(""))

		if err != nil {
			Logger{}.Error("Failed to change mobile proxy IP: %s", err)
			continue
		}

		req.Header = http.Header{
			"user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"},
		}

		_, err = c.Do(req)
		if err != nil {
			Logger{}.Error("Failed to change mobile proxy IP: %s", err)
			continue
		}
		return true
	}
	return false
}
