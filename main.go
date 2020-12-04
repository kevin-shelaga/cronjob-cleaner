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

	k8s.Connect(helpers.IsInCluster())

	logging.Information("Getting namespaces...")
	namespaces := k8s.GetNamespaces()

	for _, n := range namespaces {
		logging.Information("Getting jobs in namespace " + n + "...")
		jobsForCleanup := k8s.GetjobsForCleanup(n, 10)

		logging.Information(strconv.Itoa(len(jobsForCleanup)) + " jobs to cleanup from " + n + " namespace!")
		for _, j := range jobsForCleanup {
			logging.Warning("Cleaning up job " + j.Name + "...")

			pods := k8s.GetJobsPod(j)

			for _, p := range pods.Items {
				tail := int64(100)
				k8s.GetPodLogs(p, &tail)

				k8s.DeletePod(p)
			}
		}
	}

	logging.Information("***** Finished cronjob-cleaner *****")
}
