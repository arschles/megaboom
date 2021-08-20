package k8s

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobCreator interface {
	Create(ctx context.Context, job *batchv1.Job, opts metav1.CreateOptions) (*batchv1.Job, error)
}

type JobDeleter interface {
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

func NewJobs(uid uuid.UUID, endpoint string, numPods, numRequests, numConcurrent uint) []*batchv1.Job {
	parallelism := int32(numPods)
	completions := int32(numPods)
	ret := []*batchv1.Job{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: jobName(uid, 1),
				Labels: map[string]string{
					"created-by": "megaboom",
				},
			},
			Spec: batchv1.JobSpec{
				Parallelism: &parallelism,
				Completions: &completions,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"created-by": "megaboom",
						},
					},
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyNever,
						Containers: []corev1.Container{
							{
								Name:  "megaboom-runner",
								Image: "ghcr.io/arschles/hey:latest",
								Command: []string{
									"hey",
									"-c",
									strconv.Itoa(int(numConcurrent)),
									"-n",
									strconv.Itoa(int(numRequests)),
									endpoint,
								},
								ImagePullPolicy: corev1.PullAlways,
							},
						},
					},
				},
			},
		},
	}
	return ret
}

func CreateJobs(
	ctx context.Context,
	cl JobCreator,
	jobs []*batchv1.Job,
) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, job := range jobs {
		curJob := job
		g.Go(func() error {
			_, err := cl.Create(ctx, curJob, metav1.CreateOptions{})
			return err
		})
	}
	return g.Wait()
}

func DeleteJobs(
	ctx context.Context,
	cl JobDeleter,
	uid uuid.UUID,
) error {
	return cl.Delete(ctx, jobName(uid, 1), metav1.DeleteOptions{})
}

func jobName(uid uuid.UUID, jobNum int) string {
	return fmt.Sprintf("megaboom-job-%s-%d", uid.String(), jobNum)
}
