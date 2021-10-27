package model

import (
	"terraform-provider-securecn/internal/escher_api/model"
)

type Cluster struct {
	ID            string
	ClusterConfig *model.KubernetesCluster
}
