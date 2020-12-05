package helpers

import "testing"

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
}

func BenchmarkGetActiveDeadlineSeconds(t *testing.B) {

	GetActiveDeadlineSeconds()
}

func TestShouldGetPodLogs(t *testing.T) {
	result := ShouldGetPodLogs()

	if result != false {
		t.Errorf("ShouldGetPodLogs() should default")
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
}

func BenchmarkGetLogTail(t *testing.B) {

	GetLogTail()
}

func TestShouldDeleteJob(t *testing.T) {
	result := ShouldDeleteJob()

	if result != false {
		t.Errorf("ShouldDeleteJob() should default")
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
}

func BenchmarkShouldDeletePod(t *testing.B) {

	ShouldDeletePod()
}
