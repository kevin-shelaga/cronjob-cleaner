package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Setenv("InCusterConfig", "false")
	os.Setenv("ActiveDeadlineSecond", "0")
	os.Setenv("GetPodLogs", "true")
	os.Setenv("LogTail", "0")
	os.Setenv("CleanFailedJob", "true")
	os.Setenv("DeleteJob", "true")
	os.Setenv("DeletePod", "true")

	main()
}
