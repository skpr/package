package color

import (
	"os"
	"strconv"

	"github.com/mgutz/ansi"
)

var matched = map[string]string{}

func init() {
	matched = make(map[string]string)
}

// Wrap the string in color.
func Wrap(in string) string {
	// Globally disable color output by setting NO_COLOR environment variable.
	if os.Getenv(EnvNoColor) != "" {
		return in
	}

	if val, ok := matched[in]; ok {
		return ansi.Color(in, val)
	}

	// Ensure that we don't go over 256.
	if len(matched) > 256 {
		return ansi.Color(in, "default")
	}

	// Save it for later.
	matched[in] = strconv.Itoa(len(matched) + 1)

	return ansi.Color(in, matched[in])
}
