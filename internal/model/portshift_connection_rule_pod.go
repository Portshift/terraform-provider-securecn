package model

import (
	"terraform-provider-securecn/internal/escher_api/model"
)

type ConnectionRulePod struct {
	ID        string
	PodConfig *model.SecureCNConnectionRulePod
}
