package colors

type Color uint64

const (
	ColorDefault Color = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

var terminalColorCodes = map[Color]string{
	ColorDefault: "\033[0m",
	ColorBlack:   "\033[30m",
	ColorRed:     "\033[31m",
	ColorGreen:   "\033[32m",
	ColorYellow:  "\033[33m",
	ColorBlue:    "\033[34m",
	ColorMagenta: "\033[35m",
	ColorCyan:    "\033[36m",
	ColorWhite:   "\033[39m",
}

func TerminalColor(c Color) string {
	if code, ok := terminalColorCodes[c]; ok {
		return code
	}

	return "\033[0m"
}
