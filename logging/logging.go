package logging

import (
	"fmt"

	helpers "github.com/kevin-shelaga/cronjob-cleaner/helpers"
)

//L interface for logging package
type L interface {
	Information(message string)
	Warning(message string)
	Error(message string)
	Critical(e error)
}

const (
	infoStr   = "INFO"
	warnStr   = "WARNING"
	errorStr  = "ERROR"
	logFormat = "{ \"timestamp\": %q, \"severity\": %q, \"message\": %q }\n"
)

//Information logs an info message a standard way
func Information(message string) {
	fmt.Printf(logFormat, helpers.GetDateTime(), infoStr, message)
}

//Warning logs a warn message a standard way
func Warning(message string) {
	fmt.Printf(logFormat, helpers.GetDateTime(), warnStr, message)
}

//Error logs an error message a standard way
func Error(message string) {
	fmt.Printf(logFormat, helpers.GetDateTime(), errorStr, message)
}

//Critical logs a critical error a standard way
func Critical(e error) {
	fmt.Printf(logFormat, helpers.GetDateTime(), errorStr, e.Error())
}
