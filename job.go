package main

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/batch"
)

const (
	defaultBoomImage    = "quay.io/arschles/boom:0.1.0"
	defaultJobNamespace = "default"
	defaultJobName      = "megaboom"
	jobKind             = "Job"
	createdByLabelKey   = "create-by"
	createdByLabelValue = "megaboom"
)

func newBoomJob(boomImage string, boomCmd boomCommand, jobNamespace, jobName string, parallelism int) *batch.Job {
	parallelismi32 := int32(parallelism)
	return &batch.Job{
		TypeMeta: unversioned.TypeMeta{
			Kind: jobKind,
		},
		ObjectMeta: api.ObjectMeta{
			Name:      jobName,
			Namespace: jobNamespace,
			Labels:    map[string]string{createdByLabelKey: createdByLabelValue},
		},
		Spec: batch.JobSpec{
			Parallelism: &parallelismi32,
			Completions: &parallelismi32,
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: map[string]string{createdByLabelKey: createdByLabelValue},
				},
				Spec: api.PodSpec{
					RestartPolicy: api.RestartPolicyNever,
					Containers: []api.Container{
						api.Container{
							Image:           boomImage,
							Command:         boomCmd.Slice(),
							ImagePullPolicy: api.PullAlways,
						},
					},
				},
			},
		},
	}
}
