package model

import "terraform-provider-securecn/escher_api/model"

type Cluster struct {
	ID            string
	ClusterConfig *model.KubernetesCluster
}
