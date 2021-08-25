package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ConfigType string

const (
	ConfigTypeInCluster ConfigType = "in-cluster"
	ConfigTypeLocalhost ConfigType = "localhost"
)

func NewClient(inCluster bool) (*kubernetes.Clientset, error) {
	var restCfg *rest.Config
	if inCluster {

		c, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		restCfg = c
	} else {
		homeDir := homedir.HomeDir()
		cfgDir := fmt.Sprintf("%s/.kube/config", homeDir)
		c, err := clientcmd.BuildConfigFromFlags("", cfgDir)
		if err != nil {
			return nil, err
		}
		restCfg = c
	}

	return kubernetes.NewForConfig(restCfg)
}
