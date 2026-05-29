package output

import (
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	Width  int
	Height int
}

func NewTerminal() Terminal {
	fd := int(os.Stdin.Fd())

	width, height, err := term.GetSize(fd)
	if err != nil {
		return Terminal{
			Width:  120,
			Height: 24,
		}
	}

	return Terminal{
		Width:  width,
		Height: height,
	}
}
