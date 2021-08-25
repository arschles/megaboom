package k8s

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobLister interface {
	List(ctx context.Context, opts metav1.ListOptions) (*batchv1.JobList, error)
}
type JobDeleter interface {
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

func NewJob(
	uid uuid.UUID,
	endpoint string,
	numPods,
	numRequests,
	numConcurrent uint,
	headers []string,
) *batchv1.Job {
	headersStr := strings.Join(headers, ",")
	parallelism := int32(numPods)
	completions := int32(numPods)
	jb := &batchv1.Job{

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
							Env: []corev1.EnvVar{
								{
									Name:  "MEGABOOM_HEADERS",
									Value: headersStr,
								},
							},
						},
					},
				},
			},
		},
	}
	AddUIDToJob(jb, uid)
	AddEndpointToJob(jb, endpoint)
	return jb
}

func DeleteJob(
	ctx context.Context,
	cl JobDeleter,
	uid uuid.UUID,
) error {
	return cl.Delete(ctx, jobName(uid, 1), metav1.DeleteOptions{
		GracePeriodSeconds: i64Ptr(0),
	})
}

func jobName(uid uuid.UUID, jobNum int) string {
	return fmt.Sprintf("megaboom-job-%s-%d", uid.String(), jobNum)
}

func i64Ptr(i int64) *int64 {
	return &i
}

func AddUIDToJob(j *batchv1.Job, uid uuid.UUID) {
	j.ObjectMeta.Labels["uid"] = uid.String()
	j.Spec.Template.Spec.Containers[0].Env = append(
		j.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{
			Name:  "MEGABOOM_UID",
			Value: uid.String(),
		},
	)
}

func GetUIDFromJob(j *batchv1.Job) (uuid.UUID, error) {
	uid, ok := j.ObjectMeta.Labels["uid"]
	if !ok {
		return uuid.Nil, fmt.Errorf("job %s has no uid label", j.Name)
	}
	return uuid.Parse(uid)
}

func AddEndpointToJob(j *batchv1.Job, endpoint string) {
	j.ObjectMeta.Labels["endpoint"] = endpoint
	j.Spec.Template.Spec.Containers[0].Env = append(
		j.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{
			Name:  "MEGABOOM_ENDPOINT",
			Value: endpoint,
		},
	)
}

func GetEndpointFromJob(j *batchv1.Job) (string, error) {
	endpoint, ok := j.ObjectMeta.Labels["endpoint"]
	if !ok {
		return "", fmt.Errorf("job %s has no endpoint label", j.Name)
	}
	return endpoint, nil
}
