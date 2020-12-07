package main

import (
	"strconv"

	"github.com/kevin-shelaga/cronjob-cleaner/helpers"
	"github.com/kevin-shelaga/cronjob-cleaner/k8s"
	"github.com/kevin-shelaga/cronjob-cleaner/logging"
)

func main() {

	logging.Information("***** Starting cronjob-cleaner *****")

	var activeDeadlineSeconds = helpers.GetActiveDeadlineSeconds()
	logging.Information("Active deadline seconds set as " + strconv.FormatFloat(activeDeadlineSeconds, 'f', 6, 64))

	var k k8s.KubernetesAPI = k8s.KubernetesAPI{Clientset: nil}

	k.Clientset = k.Connect(helpers.IsInCluster())

	logging.Information("Getting namespaces...")
	namespaces := k.GetNamespaces()

	for _, n := range namespaces {
		logging.Information("Getting jobs in namespace " + n + "...")
		jobsForCleanup := k.GetjobsForCleanup(n, activeDeadlineSeconds)

		logging.Information(strconv.Itoa(len(jobsForCleanup)) + " jobs to cleanup from " + n + " namespace!")
		for _, j := range jobsForCleanup {
			logging.Warning("Cleaning up job " + j.Name + "...")

			pods := k.GetJobsPod(j)

			for _, p := range pods.Items {

				if helpers.ShouldGetPodLogs() {
					tail := helpers.GetLogTail()
					k.GetPodLogs(p, &tail)
				}

				if helpers.ShouldDeletePod() {
					k.DeletePod(p)
				}

				if helpers.ShouldDeleteJob() {
					k.DeleteJob(j)
				}
			}
		}
	}

	logging.Information("***** Finished cronjob-cleaner *****")
}
