package k8s

import (
	"bytes"
	"context"
	"flag"
	"io"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kevin-shelaga/cronjob-cleaner/logging"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"

	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	_ "k8s.io/client-go/plugin/pkg/client/auth" //all auths
)

//k8s interface for k8s package
type k8s interface {
	Connect(inCluster bool) *kubernetes.Clientset
	GetNamespaces() []string
	GetjobsForCleanup(namespace string, activeDeadlineSeconds float64) []batch.Job
	GetJobsPod(job batch.Job) *core.PodList
	GetPodLogs(pod core.Pod, tail *int64)
	DeletePod(pod core.Pod)
}

//KubernetesAPI is the struct for k8s
type KubernetesAPI struct {
	Clientset kubernetes.Interface
}

//Connect returns new kubernetes Client
func (k KubernetesAPI) Connect(inCluster bool) kubernetes.Interface {

	var clientset kubernetes.Interface

	if k.Clientset == nil && inCluster {
		logging.Information("In cluster config will be used!")
		logging.Information("Connecting...")
		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			logging.Error(err.Error())
			panic(err.Error())
		}
		// creates the Clientset
		cs, err := kubernetes.NewForConfig(config)
		if err != nil {
			logging.Error(err.Error())
			panic(err.Error())
		}

		clientset = cs
		logging.Information("Connected!")
	} else if k.Clientset == nil && !inCluster {
		logging.Information("Out of cluster config will be used!")
		logging.Information("Connecting...")
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			logging.Error(err.Error())
			panic(err.Error())
		}

		// creates the Clientset
		cs, err := kubernetes.NewForConfig(config)
		if err != nil {
			logging.Error(err.Error())
			panic(err.Error())
		}

		clientset = cs
		logging.Information("Connected!")
	}

	return clientset
}

//GetNamespaces returns slice of all namespaces
func (k KubernetesAPI) GetNamespaces() []string {
	var result []string

	namespaces, err := k.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logging.Error(err.Error())
	}
	logging.Information("There is " + strconv.Itoa(len(namespaces.Items)) + " namespace(s) in the cluster!")

	for _, n := range namespaces.Items {
		logging.Information(n.Name)
		result = append(result, n.Name)
	}

	return result
}

//GetjobsForCleanup return list of jobs in the requested namespace that exceed the max run time in second
func (k KubernetesAPI) GetjobsForCleanup(namespace string, activeDeadlineSeconds float64) []batch.Job {

	var jobsToCleanup []batch.Job

	jobs, err := k.Clientset.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logging.Error(err.Error())
	}
	logging.Information("There is " + strconv.Itoa(len(jobs.Items)) + " job(s) in the " + namespace + " namespace!")

	for _, j := range jobs.Items {
		if j.Status.Active == 1 && j.Status.Failed == 0 && j.Status.CompletionTime == nil {
			if time.Now().Sub(j.CreationTimestamp.UTC()).Seconds() > activeDeadlineSeconds {
				logging.Warning("Job " + j.Name + " has been running for " + strconv.FormatFloat(time.Now().Sub(j.CreationTimestamp.UTC()).Seconds(), 'f', 6, 64) + " seconds and has been flagged for cleanup due to exceeding the active deadline seconds of " + strconv.FormatFloat(activeDeadlineSeconds, 'f', 6, 64))
				jobsToCleanup = append(jobsToCleanup, j)
			} else {
				logging.Information("Job " + j.Name + " has been running for " + strconv.FormatFloat(time.Now().Sub(j.CreationTimestamp.UTC()).Seconds(), 'f', 6, 64) + " seconds and has not exceeded the active deadline seconds of " + strconv.FormatFloat(activeDeadlineSeconds, 'f', 6, 64))
			}
		} else {
			logging.Information("Job " + j.Name + " is not active")
		}
	}

	return jobsToCleanup
}

//GetJobsPod gets the pod associated to a job based on the job-name label
func (k KubernetesAPI) GetJobsPod(job batch.Job) *core.PodList {

	pods, err := k.Clientset.CoreV1().Pods(job.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "job-name=" + job.Name})
	if err != nil {
		logging.Error(err.Error())

		return nil
	}

	if len(pods.Items) > 1 {
		logging.Error("Found more than 1 pod matching label " + ("job-name:" + job.Name) + ", cleanup will be skipped")

		return nil
	} 

		logging.Information("Found pod " + pods.Items[0].Name + " for job " + job.Name + " with label " + ("job-name:" + job.Name))
	
	return pods
}

//GetPodLogs logs a list of logs for a given pod name
func (k KubernetesAPI) GetPodLogs(pod core.Pod, tail *int64) {

	logging.Information("Getting the last " + strconv.FormatInt(*tail, 10) + " logs for pod " + pod.Name)
	for _, c := range pod.Spec.Containers {

		req := k.Clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &core.PodLogOptions{Container: c.Name, TailLines: tail})
		podLogs, err := req.Stream(context.TODO())
		if err != nil {
			logging.Error("error in opening stream")
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			logging.Error("error in copy information from podLogs to buf")
		}
		logging.Information(buf.String())
	}
}

//DeletePod deletes pod from kubernetes
func (k KubernetesAPI) DeletePod(pod core.Pod) {

	logging.Warning("Deleting pod " + pod.Name + "...")
	derr := k.Clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
	if derr != nil {
		logging.Error(derr.Error())
	}
	logging.Warning("Deleted pod " + pod.Name + "!")
}

//DeleteJob deletes pod from kubernetes
func (k KubernetesAPI) DeleteJob(job batch.Job) {

	logging.Warning("Deleting job " + job.Name + "...")
	derr := k.Clientset.BatchV1().Jobs(job.Namespace).Delete(context.TODO(), job.Name, metav1.DeleteOptions{})
	if derr != nil {
		logging.Error(derr.Error())
	}
	logging.Warning("Deleted jod " + job.Name + "!")
}
