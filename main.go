package main

import (
	"strconv"
	"time"

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
	namespaces := k.GetNamespaces(helpers.ExcludedNamespaces())

	for _, n := range namespaces {
		logging.Information("Getting jobs in namespace " + n + "...")
		jobsForCleanup := k.GetjobsForCleanup(n, activeDeadlineSeconds)

		logging.Information(strconv.Itoa(len(jobsForCleanup)) + " jobs to cleanup from " + n + " namespace!")
		for _, j := range jobsForCleanup {
			logging.Warning("Cleaning up job " + j.Name + "...")

			pods := k.GetJobsPod(j)

			if pods != nil {

				for _, p := range pods.Items {

					if helpers.ShouldGetPodLogs() {
						tail := helpers.GetLogTail()
						k.GetPodLogs(p, &tail)
					}

					if helpers.ShouldDeletePod() {
						if _, ok := j.Labels["deleted-pod-timestamp"]; !ok {
							k.LabelJob(j)
							k.DeletePod(p)
						} else {
							logging.Warning("Pod previously deleted!")

							i, err := strconv.ParseInt(j.Labels["deleted-pod-timestamp"], 10, 64)
							if err != nil {
								logging.Error(err.Error())
							}
							tm := time.Unix(i, 0)

							if time.Now().Sub(tm.UTC()).Seconds() > activeDeadlineSeconds {
								logging.Warning("Pod was deleted " + strconv.FormatFloat(time.Now().Sub(tm.UTC()).Seconds(), 'f', 6, 64) + " seconds ago and now exceeds to ActiveDeadlineSeconds of " + strconv.FormatFloat(activeDeadlineSeconds, 'f', 6, 64) + " seconds ago, pod will be deleted again!")
								k.LabelJob(j)
								k.DeletePod(p)
							} else {
								logging.Warning("Pod will not be deleted again this time!")
							}

						}
					}
				}
			}

			if helpers.ShouldDeleteJob() {
				k.DeleteJob(j)
			}
		}
	}

	logging.Information("***** Finished cronjob-cleaner *****")
}
