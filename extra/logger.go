package extra

import (
	"fmt"
	"github.com/gookit/color"
	"time"
)

type Logger struct{}

func (l Logger) Info(format string, a ...any) {
	color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=white;op=bold> | INFO    | </><fg=white>-</> %s\n", time.Now().Format("15:04:05.000"), fmt.Sprintf(format, a...))
}

func (l Logger) Error(format string, a ...any) {
	color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=red;op=bold> | ERROR   | </><fg=white>-</> %s\n", time.Now().Format("15:04:05.000"), fmt.Sprintf(format, a...))
}

func (l Logger) Success(format string, a ...any) {
	color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=green;op=bold> | SUCCESS | </><fg=white>-</> %s\n", time.Now().Format("15:04:05.000"), fmt.Sprintf(format, a...))
}

func (l Logger) Warning(format string, a ...any) {
	color.Printf("<fg=cyan;op=bold>%s</><fg=white></><fg=white></><fg=yellow;op=bold> | WARNING | </><fg=white>-</> %s\n", time.Now().Format("15:04:05.000"), fmt.Sprintf(format, a...))
}
