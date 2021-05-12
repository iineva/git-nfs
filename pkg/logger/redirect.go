package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	_stdout = os.Stdout
	_stderr = os.Stderr
)

// hook os.Stdout output from go-nfs or other module
func RedirectStd() {
	f, err := os.OpenFile("/tmp/stdout-and-go-log.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Print()
	os.Stdout = f
	os.Stderr = f
	log.SetOutput(f)
}

func RecoverRedirectStd() {
	os.Stdout = _stdout
	os.Stderr = _stderr
	log.SetOutput(_stderr)
}
