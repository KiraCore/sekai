package database

import (
	"os"

	"github.com/KiraCore/sekai/INTERX/global"
)

var temp = os.Stdout

// DisableStdout is a function to disable stdout
func DisableStdout() {
	global.Mutex.Lock()
	os.Stdout = nil // turn it off
}

// EnableStdout is a function to enable stdout
func EnableStdout() {
	os.Stdout = temp // turn it off
	global.Mutex.Unlock()
}
