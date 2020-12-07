package k8s

import (
	// "strconv"
	"errors"
	"testing"
	"time"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fake "k8s.io/client-go/kubernetes/fake"
	fakecore "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestConnect(t *testing.T) {

	//out of cluster
	var k KubernetesAPI = KubernetesAPI{Clientset: nil}

	k.Clientset = k.Connect(false)

	if k.Clientset == nil {
		t.Errorf("Clientset should not be nil")
	}
}

func BenchmarkIsInCluster(t *testing.B) {

	//out of cluster
	var k KubernetesAPI = KubernetesAPI{Clientset: nil}

	k.Clientset = k.Connect(false)
}

func TestGetNamespaces(t *testing.T) {

	//no namespaces
	var k KubernetesAPI = KubernetesAPI{Clientset: fake.NewSimpleClientset()}
	result := k.GetNamespaces()

	if result != nil {
		t.Errorf("result should be nil")
	}

	//one namespace
	namespace := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	k.Clientset = fake.NewSimpleClientset(namespace)
	result = k.GetNamespaces()

	if result == nil {
		t.Errorf("result should not be nil")
	}

	//error getting namespaces
	k.Clientset = fake.NewSimpleClientset()

	k.Clientset.CoreV1().(*fakecore.FakeCoreV1).PrependReactor("list", "namespaces", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &core.NamespaceList{}, errors.New("Error listing namespaces")
	})

	result = k.GetNamespaces()

	if result != nil {
		t.Errorf("result should be nil")
	}
}

func BenchmarkGetNamespaces(t *testing.B) {

	//no namespaces
	var k KubernetesAPI = KubernetesAPI{Clientset: fake.NewSimpleClientset()}
	result := k.GetNamespaces()

	if result != nil {
		t.Errorf("result should be nil")
	}

	//one namespace
	namespace := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	k.Clientset = fake.NewSimpleClientset(namespace)
	result = k.GetNamespaces()

	if result == nil {
		t.Errorf("result should not be nil")
	}

	//error getting namespaces
	k.Clientset = fake.NewSimpleClientset()

	k.Clientset.CoreV1().(*fakecore.FakeCoreV1).PrependReactor("list", "namespaces", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &core.NamespaceList{}, errors.New("Error listing namespaces")
	})

	result = k.GetNamespaces()

	if result != nil {
		t.Errorf("result should be nil")
	}
}

func TestGetjobsForCleanup(t *testing.T) {

	//no jobs
	var k KubernetesAPI = KubernetesAPI{Clientset: fake.NewSimpleClientset()}
	result := k.GetjobsForCleanup("default", 4200)

	if result != nil {
		t.Errorf("result should be nil")
	}

	//one job - not active
	namespace := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	job := &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Namespace:   "default",
			Annotations: map[string]string{},
		},
	}

	k.Clientset = fake.NewSimpleClientset(namespace, job)
	result = k.GetjobsForCleanup("default", 4200)

	if result == nil {
		t.Errorf("result should not be nil")
	}

	//one job - active - not long enough
	namespace = &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	job = &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "default",
			Namespace:         "default",
			Annotations:       map[string]string{},
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Status: batch.JobStatus{
			Active:         1,
			Failed:         0,
			CompletionTime: nil,
		},
	}

	k.Clientset = fake.NewSimpleClientset(namespace, job)
	result = k.GetjobsForCleanup("default", 4200)

	if result == nil {
		t.Errorf("result should not be nil")
	}

	//one job - active - long enough
	namespace = &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	job = &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "default",
			Namespace:         "default",
			Annotations:       map[string]string{},
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Status: batch.JobStatus{
			Active:         1,
			Failed:         0,
			CompletionTime: nil,
		},
	}

	k.Clientset = fake.NewSimpleClientset(namespace, job)
	result = k.GetjobsForCleanup("default", 0)

	if result == nil {
		t.Errorf("result should not be nil")
	}
}

// func BenchmarkGetjobsForCleanup(t *testing.B) {

// 	Connect(false)
// 	namespaces := GetNamespaces()

// 	for _, n := range namespaces {
// 		GetjobsForCleanup(n, 4200)
// 	}
// }

// func TestGetJobsPod(t *testing.T) {

// 	Connect(false)
// 	namespaces := GetNamespaces()

// 	for _, n := range namespaces {
// 		jobs := GetjobsForCleanup(n, 30)

// 		for _, j := range jobs {
// 			GetJobsPod(j)
// 		}
// 	}
// }

// func BenchmarkGetJobsPod(t *testing.B) {

// 	Connect(false)
// 	namespaces := GetNamespaces()

// 	for _, n := range namespaces {
// 		jobs := GetjobsForCleanup(n, 30)

// 		for _, j := range jobs {
// 			GetJobsPod(j)
// 		}
// 	}
// }

// func TestGetPodLogs(t *testing.T) {

// 	logTail, _ := strconv.ParseInt("100", 10, 64)

// 	Connect(false)
// 	namespaces := GetNamespaces()

// 	for _, n := range namespaces {
// 		jobs := GetjobsForCleanup(n, 30)

// 		for _, j := range jobs {
// 			pods := GetJobsPod(j)

// 			for _, p := range pods.Items {
// 				GetPodLogs(p, &logTail)
// 			}
// 		}
// 	}
// }

// func BenchmarkGetPodLogs(t *testing.B) {

// 	logTail, _ := strconv.ParseInt("100", 10, 64)

// 	Connect(false)
// 	namespaces := GetNamespaces()

// 	for _, n := range namespaces {
// 		jobs := GetjobsForCleanup(n, 30)

// 		for _, j := range jobs {
// 			pods := GetJobsPod(j)

// 			for _, p := range pods.Items {
// 				GetPodLogs(p, &logTail)
// 			}
// 		}
// 	}
// }

// func TestDeletePod(t *testing.T) {

// 	Connect(false)
// 	namespaces := GetNamespaces()

// 	for _, n := range namespaces {
// 		jobs := GetjobsForCleanup(n, 4200)

// 		for _, j := range jobs {
// 			pods := GetJobsPod(j)

// 			for _, p := range pods.Items {

// 				DeletePod(p)
// 			}
// 		}
// 	}
// }

// func BenchmarkDeletePod(t *testing.B) {

// 	Connect(false)
// 	namespaces := GetNamespaces()

// 	for _, n := range namespaces {
// 		jobs := GetjobsForCleanup(n, 4200)

// 		for _, j := range jobs {
// 			pods := GetJobsPod(j)

// 			for _, p := range pods.Items {

// 				DeletePod(p)
// 			}
// 		}
// 	}
// }

// func TestDeleteJob(t *testing.T) {

// 	Connect(false)
// 	namespaces := GetNamespaces()

// 	for _, n := range namespaces {
// 		jobs := GetjobsForCleanup(n, 4200)

// 		for _, j := range jobs {

// 			DeleteJob(j)
// 		}
// 	}
// }

// func BenchmarkDeleteJob(t *testing.B) {

// 	Connect(false)
// 	namespaces := GetNamespaces()

// 	for _, n := range namespaces {
// 		jobs := GetjobsForCleanup(n, 4200)

// 		for _, j := range jobs {

// 			DeleteJob(j)
// 		}
// 	}
// }
