package logging

import (
	"testing"
	"errors"
)
func TestInformation(t *testing.T) {

	Information("test")
}

func BenchmarkInformation(t *testing.B) {

	Information("test")
}

func TestWarning(t *testing.T) {

	Warning("test")
}

func BenchmarkWarning(t *testing.B) {

	Warning("test")
}


func TestError(t *testing.T) {

	Error("test")
}

func BenchmarkError(t *testing.B) {

	Error("test")
}


func TestCritical(t *testing.T) {

	err := errors.New("test")
	Critical(err)
}

func BenchmarkCritical(t *testing.B) {

	err := errors.New("test")
	Critical(err)
}