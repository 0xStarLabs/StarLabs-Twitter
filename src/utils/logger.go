package utils

import (
	"fmt"
	"time"

	"github.com/gookit/color"
)

type Logger struct {
	Prefix string
}

func (l Logger) formatMessage(message string) string {
	if l.Prefix != "" {
		return fmt.Sprintf("%s: %s", l.Prefix, message)
	}
	return message
}

func (l Logger) Info(format string, a ...any) {
	message := l.formatMessage(fmt.Sprintf(format, a...))
	color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=white;op=bold> | INFO    | </><fg=white>-</> %s\n", time.Now().Format("15:04:05.000"), message)
}

func (l Logger) Error(format string, a ...any) {
	message := l.formatMessage(fmt.Sprintf(format, a...))
	color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=red;op=bold> | ERROR   | </><fg=white>-</> %s\n", time.Now().Format("15:04:05.000"), message)
}

func (l Logger) Success(format string, a ...any) {
	message := l.formatMessage(fmt.Sprintf(format, a...))
	color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=green;op=bold> | SUCCESS | </><fg=white>-</> %s\n", time.Now().Format("15:04:05.000"), message)
}

func (l Logger) Warning(format string, a ...any) {
	message := l.formatMessage(fmt.Sprintf(format, a...))
	color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=yellow;op=bold> | WARNING | </><fg=white>-</> %s\n", time.Now().Format("15:04:05.000"), message)
}
