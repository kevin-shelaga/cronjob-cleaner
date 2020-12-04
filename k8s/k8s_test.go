package k8s

import "testing"

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
