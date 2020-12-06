# cronjob-cleaner
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/kevin-shelaga/cronjob-cleaner/branch/master/graph/badge.svg?token=D07EP88G53)](https://codecov.io/gh/kevin-shelaga/cronjob-cleaner)
[![Go Report Card](https://goreportcard.com/badge/github.com/kevin-shelaga/cronjob-cleaner)](https://goreportcard.com/report/github.com/kevin-shelaga/cronjob-cleaner)
![build](https://github.com/kevin-shelaga/cronjob-cleaner/workflows/build/badge.svg)

## Why

I wrote this cleaner as a way to delete pods for jobs that got "stuck", jobs that just seemed to stop doing anything but the job and pod were both in a running state for hours. The activeDeadlineSecond would delete the job, but this wouldnt resolve any alerts from alert manager. This project is my fix, and a good excuse to write go.

## How

### Environment Variables

| Environment Variable | Type    | Description                                                                                                                                         |
| -------------------- | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| InCusterConfig       | bool    | Used to determine which clientset should be used for kubernetes authentication.Defaults to false, set to true to use inside a cluster as a cronjob. |
| ActiveDeadlineSecond | float64 | Used to determine which jobs/pods should be identified for cleanup. Defaults to 4200.                                                               |
| GetPodLogs           | bool    | Used to determine if the pods logs should be relogged as Information. Defaults to false.                                                            |
| LogTail              | int64   | Used to determine the tail of logs that should be relogged, if GetPodLogs is set to true, otherwise it is ignored. Defaults to 100.                 |
| DeleteJob            | bool    | Used to determine if jobs that are identified for cleanup should be deleted. Defaults to false.                                                     |
| DeletePod            | bool    | Used to determine pods that are identified for cleanup should be deleted. Defaults                                                                  |
to false.

### Run as cronjob in kubernetes

#### Use Kustomize

```sh
curl -sfLo kustomize https://github.com/kubernetes-sigs/kustomize/releases/download/v3.1.0/kustomize_30_linux_amd64
chmod u+x ./kustomi

cp -a ./manifests/. .

./kustomize edit set image kevinshelaga/$IMAGE:tag=kevinshelaga/$IMAGE:latest
./kustomize build . | kubectl -n $K8S_NAMESPACE apply -f -
```

## Whats left

### TODO

- [ ] More/better tests
- [ ] Helm chart
- [ ] Policies around contributions
