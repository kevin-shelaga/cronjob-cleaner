package k8s

import (
	"strconv"
	"testing"
)

func TestConnect(t *testing.T) {
	//out of cluster
	Connect(false)

	if clientset == nil {
		t.Errorf("Clientset should not be nil")
	}
}

func BenchmarkIsInCluster(t *testing.B) {

	Connect(false)
}

func TestGetNamespaces(t *testing.T) {

	Connect(false)
	result := GetNamespaces()

	if result == nil {
		t.Errorf("result should not be nil")
	}
}

func BenchmarkGetNamespaces(t *testing.B) {

	Connect(false)
	GetNamespaces()
}

func TestGetjobsForCleanup(t *testing.T) {

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		GetjobsForCleanup(n, 4200)
	}
}

func BenchmarkGetjobsForCleanup(t *testing.B) {

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		GetjobsForCleanup(n, 4200)
	}
}

func TestGetJobsPod(t *testing.T) {

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		jobs := GetjobsForCleanup(n, 30)

		for _, j := range jobs {
			GetJobsPod(j)
		}
	}
}

func BenchmarkGetJobsPod(t *testing.B) {

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		jobs := GetjobsForCleanup(n, 30)

		for _, j := range jobs {
			GetJobsPod(j)
		}
	}
}

func TestGetPodLogs(t *testing.T) {

	logTail, _ := strconv.ParseInt("100", 10, 64)

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		jobs := GetjobsForCleanup(n, 30)

		for _, j := range jobs {
			pods := GetJobsPod(j)

			for _, p := range pods.Items {
				GetPodLogs(p, &logTail)
			}
		}
	}
}

func BenchmarkGetPodLogs(t *testing.B) {

	logTail, _ := strconv.ParseInt("100", 10, 64)

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		jobs := GetjobsForCleanup(n, 30)

		for _, j := range jobs {
			pods := GetJobsPod(j)

			for _, p := range pods.Items {
				GetPodLogs(p, &logTail)
			}
		}
	}
}

func TestDeletePod(t *testing.T) {

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		jobs := GetjobsForCleanup(n, 4200)

		for _, j := range jobs {
			pods := GetJobsPod(j)

			for _, p := range pods.Items {

				DeletePod(p)
			}
		}
	}
}

func BenchmarkDeletePod(t *testing.B) {

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		jobs := GetjobsForCleanup(n, 4200)

		for _, j := range jobs {
			pods := GetJobsPod(j)

			for _, p := range pods.Items {

				DeletePod(p)
			}
		}
	}
}

func TestDeleteJob(t *testing.T) {

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		jobs := GetjobsForCleanup(n, 4200)

		for _, j := range jobs {

			DeleteJob(j)
		}
	}
}

func BenchmarkDeleteJob(t *testing.B) {

	Connect(false)
	namespaces := GetNamespaces()

	for _, n := range namespaces {
		jobs := GetjobsForCleanup(n, 4200)

		for _, j := range jobs {

			DeleteJob(j)
		}
	}
}
