package extra

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

var mutex sync.Mutex

func AppendToFile(filepath string, data string) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(data + "\n")
	if err != nil {
		return
	}
}

func RewriteChannelToFile(filepath string, usernames chan string) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	for username := range usernames {
		_, err := file.WriteString(username + "\n")
		if err != nil {
			return
		}
	}
}

func ReplaceLineInFile(filename, oldLine, newLine string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == oldLine {
			line = newLine
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}
