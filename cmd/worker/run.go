package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type result struct {
	// map from error code to count of that code
	RespCodes map[int]int `json:"response_codes"`
}

func run(
	ctx context.Context,
	lggr logr.Logger,
	totalReqs,
	concurrency int,
	endpoint string,
) (result, error) {
	return result{}, nil
}

func configMapFromResult(cmName string, runnerID uuid.UUID, res result) (corev1.ConfigMap, error) {
	jsonBytes, err := json.Marshal(&res)
	if err != nil {
		return corev1.ConfigMap{}, err
	}

	name := fmt.Sprintf("%s-%s", cmName, runnerID.String())
	return corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			"runnerID":  runnerID.String(),
			"respCodes": string(jsonBytes),
		},
	}, nil

}
