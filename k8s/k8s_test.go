package k8s

import (
	// "strconv"
	"errors"
	"os"
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

func TestGetNamespaces(t *testing.T) {

	//no namespaces
	var k KubernetesAPI = KubernetesAPI{Clientset: fake.NewSimpleClientset()}
	result := k.GetNamespaces(nil)

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
	result = k.GetNamespaces(nil)

	if result == nil {
		t.Errorf("result should not be nil")
	}

	//error getting namespaces
	k.Clientset = fake.NewSimpleClientset()

	k.Clientset.CoreV1().(*fakecore.FakeCoreV1).PrependReactor("list", "namespaces", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &core.NamespaceList{}, errors.New("Error listing namespaces")
	})

	result = k.GetNamespaces(nil)

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

	if result != nil {
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

	if result != nil {
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

	//one job - failed
	os.Setenv("CleanFailedJob", "true")

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
			Active:         0,
			Failed:         1,
			CompletionTime: nil,
		},
	}

	k.Clientset = fake.NewSimpleClientset(namespace, job)
	result = k.GetjobsForCleanup("default", 0)

	if result == nil {
		t.Errorf("result should not be nil")
	}

	//no job - error
	k.Clientset = fake.NewSimpleClientset()

	k.Clientset.CoreV1().(*fakecore.FakeCoreV1).PrependReactor("list", "jobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &batch.JobList{}, errors.New("Error deleting jobs")
	})

	result = k.GetjobsForCleanup("default", 0)

	if result != nil {
		t.Errorf("result should be nil")
	}
}

func TestGetJobsPod(t *testing.T) {

	//one pod
	var k KubernetesAPI = KubernetesAPI{Clientset: nil}

	namespace := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	var label = map[string]string{}
	label["job-name"] = "default"

	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "default",
			Namespace:         "default",
			Annotations:       map[string]string{},
			Labels:            label,
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
	}

	job := &batch.Job{
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

	k.Clientset = fake.NewSimpleClientset(namespace, job, pod)
	result := k.GetJobsPod(*job)

	if result == nil {
		t.Errorf("result should not be nil")
	}

	//two pods
	namespace = &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	label = map[string]string{}
	label["job-name"] = "default"

	pod1 := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "default1",
			Namespace:         "default",
			Annotations:       map[string]string{},
			Labels:            label,
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
	}

	pod2 := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "default2",
			Namespace:         "default",
			Annotations:       map[string]string{},
			Labels:            label,
			CreationTimestamp: metav1.NewTime(time.Now()),
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

	k.Clientset = fake.NewSimpleClientset(namespace, job, pod1, pod2)
	result = k.GetJobsPod(*job)

	if result != nil {
		t.Errorf("result should be nil")
	}

	//no pods
	namespace = &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	label = map[string]string{}
	label["job-name"] = "default"

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
	result = k.GetJobsPod(*job)

	if result != nil {
		t.Errorf("result should be nil")
	}

	//error getting pods
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

	k.Clientset = fake.NewSimpleClientset(namespace, job, pod1, pod2)

	k.Clientset.CoreV1().(*fakecore.FakeCoreV1).PrependReactor("list", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &core.PodList{}, errors.New("Error listing pods")
	})

	result = k.GetJobsPod(*job)

	if result != nil {
		t.Errorf("result should be nil")
	}

}

func TestGetPodLogs(t *testing.T) {

	//get one pods logs
	var k KubernetesAPI = KubernetesAPI{Clientset: nil}

	namespace := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	var label = map[string]string{}
	label["job-name"] = "default"

	var containers []core.Container
	var container = new(core.Container)
	container.Name = "default"
	containers = append(containers, *container)

	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "default",
			Namespace:         "default",
			Annotations:       map[string]string{},
			Labels:            label,
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Spec: core.PodSpec{
			Containers: containers,
		},
	}

	k.Clientset = fake.NewSimpleClientset(namespace, pod)
	tail := new(int64)
	k.GetPodLogs(*pod, tail)

	//error
	// k.Clientset = fake.NewSimpleClientset(namespace, pod)

	// k.Clientset.CoreV1().(*fakecore.FakeCoreV1).PrependReactor("get", "podlogs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
	// 	return true, &rest.Request{}, errors.New("Error getting pod logs")
	// })

	// k.GetPodLogs(*pod, tail)
}

func TestDeletePod(t *testing.T) {

	//delete one pod
	var k KubernetesAPI = KubernetesAPI{Clientset: nil}

	namespace := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	var label = map[string]string{}
	label["job-name"] = "default"

	var containers []core.Container
	var container = new(core.Container)
	container.Name = "default"
	containers = append(containers, *container)

	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "default",
			Namespace:         "default",
			Annotations:       map[string]string{},
			Labels:            label,
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Spec: core.PodSpec{
			Containers: containers,
		},
	}

	k.Clientset = fake.NewSimpleClientset(namespace, pod)
	k.DeletePod(*pod)

	//delete pod error
	k.Clientset = fake.NewSimpleClientset()

	k.Clientset.CoreV1().(*fakecore.FakeCoreV1).PrependReactor("delete", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &core.Pod{}, errors.New("Error deleting pod")
	})

	k.DeletePod(*pod)
}

func TestDeleteJob(t *testing.T) {

	//delete one pod
	var k KubernetesAPI = KubernetesAPI{Clientset: nil}

	namespace := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "default",
			Annotations: map[string]string{},
		},
	}

	job := &batch.Job{
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
	k.DeleteJob(*job)

	//delete pod error
	k.Clientset = fake.NewSimpleClientset()

	k.Clientset.CoreV1().(*fakecore.FakeCoreV1).PrependReactor("delete", "jobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &batch.Job{}, errors.New("Error deleting jobs")
	})

	k.DeleteJob(*job)
}
