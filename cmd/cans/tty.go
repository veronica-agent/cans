package main

import "os"

// stdinIsTTY is true only for a real terminal. /dev/null is a char device
// (Stat ModeCharDevice is set) but not a TTY; treating it as one made
// `cans say < /dev/null` print "missing text" instead of "empty text".
func stdinIsTTY() bool {
	f, ok := stdin.(*os.File)
	if !ok {
		return false
	}
	return fdIsTTY(f.Fd())
}
