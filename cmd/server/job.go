package main

const (
// defaultJobName      = "megaboom"
// jobKind             = "Job"
// createdByLabelKey   = "create-by"
// createdByLabelValue = "megaboom"
// containerName       = "megaboom"
)

// func newBoomJob(boomImage string, boomCmd boomCommand, jobNamespace, jobName string, parallelism int) *batch.Job {
// 	parallelismi32 := int32(parallelism)
// 	return &batch.Job{
// 		TypeMeta: unversioned.TypeMeta{
// 			Kind: jobKind,
// 		},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      jobName,
// 			Namespace: jobNamespace,
// 			Labels:    map[string]string{createdByLabelKey: createdByLabelValue},
// 		},
// 		Spec: batch.JobSpec{
// 			Parallelism: &parallelismi32,
// 			Completions: &parallelismi32,
// 			Template: api.PodTemplateSpec{
// 				ObjectMeta: api.ObjectMeta{
// 					Labels: map[string]string{createdByLabelKey: createdByLabelValue},
// 				},
// 				Spec: api.PodSpec{
// 					RestartPolicy: api.RestartPolicyNever,
// 					Containers: []api.Container{
// 						api.Container{
// 							Name:            containerName,
// 							Image:           boomImage,
// 							Command:         boomCmd.Slice(),
// 							ImagePullPolicy: api.PullAlways,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }
