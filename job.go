package megaboom

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/batch"
)

const (
	defaultBoomImage = "quay.io/arschles/boom:0.1.0"
)

func newBoomJob(boomImage string) *batch.Job {
	return &batch.Job{
		TypeMeta:   unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{},
		Spec:       batch.JobSpec{},
	}
}
