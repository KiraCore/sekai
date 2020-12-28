package database

import "os"

var temp = os.Stdout

// DisableStdout is a function to disable stdout
func DisableStdout() {
	os.Stdout = nil // turn it off
}

// EnableStdout is a function to enable stdout
func EnableStdout() {
	os.Stdout = temp // turn it off
}
