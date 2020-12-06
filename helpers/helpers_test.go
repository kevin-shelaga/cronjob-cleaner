package helpers

import (
	"os"
	"testing"
)

func TestIsInCluster(t *testing.T) {
	result := IsInCluster()

	if result != false {
		t.Errorf("IsInCluster() = true; want false")
	}
}

func BenchmarkIsInCluster(t *testing.B) {

	IsInCluster()
}

func TestGetDateTime(t *testing.T) {
	result := GetDateTime()

	if result == "" {
		t.Errorf("GetDateTime() = \"\"; want time.Now formatted")
	}
}

func BenchmarkGetDateTime(t *testing.B) {

	GetDateTime()
}

func TestGetActiveDeadlineSeconds(t *testing.T) {
	result := GetActiveDeadlineSeconds()

	if result != activeDeadlineSeconds {
		t.Errorf("GetActiveDeadlineSeconds() should default")
	}

	os.Setenv("ActiveDeadlineSecond", "3600")

	result = GetActiveDeadlineSeconds()

	if result != 3600 {
		t.Errorf("GetActiveDeadlineSeconds() should match env var")
	}
}

func BenchmarkGetActiveDeadlineSeconds(t *testing.B) {

	GetActiveDeadlineSeconds()
}

func TestShouldGetPodLogs(t *testing.T) {
	result := ShouldGetPodLogs()

	if result != false {
		t.Errorf("ShouldGetPodLogs() should default")
	}

	os.Setenv("GetPodLogs", "true")

	result = ShouldGetPodLogs()

	if result != true {
		t.Errorf("ShouldGetPodLogs() should match env var")
	}
}

func BenchmarkShouldGetPodLogs(t *testing.B) {

	ShouldGetPodLogs()
}

func TestGetLogTail(t *testing.T) {
	result := GetLogTail()

	if result != logTail {
		t.Errorf("GetLogTail() should default")
	}

	os.Setenv("LogTail", "99")

	result = GetLogTail()

	if result != 99 {
		t.Errorf("GetLogTail() should match env var")
	}
}

func BenchmarkGetLogTail(t *testing.B) {

	GetLogTail()
}

func TestShouldDeleteJob(t *testing.T) {
	result := ShouldDeleteJob()

	if result != false {
		t.Errorf("ShouldDeleteJob() should default")
	}

	os.Setenv("DeleteJob", "true")

	result = ShouldDeleteJob()

	if result != true {
		t.Errorf("ShouldDeleteJob() should match env var")
	}
}

func BenchmarkShouldDeleteJob(t *testing.B) {

	ShouldDeleteJob()
}

func TestShouldDeletePod(t *testing.T) {
	result := ShouldDeletePod()

	if result != false {
		t.Errorf("ShouldDeletePod() should default")
	}

	os.Setenv("DeletePod", "true")

	result = ShouldDeletePod()

	if result != true {
		t.Errorf("ShouldDeletePod() should match env var")
	}
}

func BenchmarkShouldDeletePod(t *testing.B) {

	ShouldDeletePod()
}
