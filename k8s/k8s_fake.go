package k8s

import (
	"k8s.io/client-go/kubernetes/fake"
)

func newClientset() *k8s {
	client := k8s{}
	client.clientset = fake.NewSimpleClientset()
	return &client.clientset
}
