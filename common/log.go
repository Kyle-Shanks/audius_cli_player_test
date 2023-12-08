// Logging and debugging things
package common

import (
	"os"
	"path/filepath"
)

// TODO: Make a logging manager or something

func log(string string) {
	f, err := os.OpenFile(filepath.Join(GetDataPath(), "debug.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err := f.WriteString(string + "\n"); err != nil {
		panic(err)
	}
}

func errorLog(string string) {
	log("ERROR: " + string)
}
