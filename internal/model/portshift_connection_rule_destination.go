package model

import (
	"terraform-provider-securecn/internal/escher_api/model"
)

type ConnectionRuleDestination struct {
	ID                string
	DestinationConfig *model.SecureCNConnectionRuleDestination
}
