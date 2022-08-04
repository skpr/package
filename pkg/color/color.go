package color

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

const (
	// Green is an alias for FgGreen.
	Green string = "green"
	// Blue is an alias for FgBlue.
	Blue string = "blue"
	// Red is an alias for FgRed.
	Red string = "red"
)

// Fprint is a wrapper for color.Fprint().
func Fprint(out io.Writer, c, s string) error {
	// Globally disable color output by setting NO_COLOR environment variable.
	if os.Getenv(EnvNoColor) != "" {
		_, err := fmt.Fprint(out, s)
		return err
	}

	cv := getColor(c)
	v := color.New(cv)
	_, err := v.Fprint(out, s)
	return err
}

func getColor(c string) color.Attribute {
	switch c {
	case Green:
		return color.FgGreen

	case Blue:
		return color.FgBlue

	case Red:
		return color.FgRed
	}

	return color.FgWhite
}
