package model

import "terraform-provider-securecn/escher_api/model"

type ConnectionRulePod struct {
	ID        string
	PodConfig *model.SecureCNConnectionRulePod
}
