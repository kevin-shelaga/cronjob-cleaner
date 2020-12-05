package helpers

import (
	"os"
	"strconv"
	"time"
)

//H interface for helpers package
type H interface {
	IsInCluster() bool
	GetDateTime() string
	GetActiveDeadlineSeconds() float64
	GetLogTail() int64
}

const (
	activeDeadlineSeconds = 4200
	logTail               = 100
)

//IsInCluster returns true or false based on the InCluster env var
func IsInCluster() bool {
	result := os.Getenv("InCusterConfig")

	if result != "" {
		b, _ := strconv.ParseBool(result)
		return b
	}

	return false
}

//GetDateTime returns date time now as string
func GetDateTime() string {
	return (time.Now().Format("2006-01-02 15:04:05.000"))
}

//GetActiveDeadlineSeconds returns ActiveDeadlineSeconds env var or 4200
func GetActiveDeadlineSeconds() float64 {
	result := os.Getenv("ActiveDeadlineSecond")

	if result != "" {
		f, _ := strconv.ParseFloat(result, 64)
		return f
	}
	return activeDeadlineSeconds
}

//ShouldGetPodLogs returns GetPodLogs env var or false
func ShouldGetPodLogs() bool {
	result := os.Getenv("GetPodLogs")

	if result != "" {
		b, _ := strconv.ParseBool(result)
		return b
	}

	return false
}

//GetLogTail returns LogTail env var or 100
func GetLogTail() int64 {
	result := os.Getenv("LogTail")

	if result != "" {
		f, _ := strconv.ParseInt(result, 10, 64)
		return f
	}
	return logTail
}

//ShouldDeleteJob returns GetPodLogs env var or false
func ShouldDeleteJob() bool {
	result := os.Getenv("DeleteJob")

	if result != "" {
		b, _ := strconv.ParseBool(result)
		return b
	}

	return false
}

//ShouldDeletePod returns GetPodLogs env var or false
func ShouldDeletePod() bool {
	result := os.Getenv("DeletePod")

	if result != "" {
		b, _ := strconv.ParseBool(result)
		return b
	}

	return false
}
