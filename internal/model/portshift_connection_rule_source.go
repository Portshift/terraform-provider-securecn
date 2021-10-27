package model

import (
	"terraform-provider-securecn/internal/escher_api/model"
)

type ConnectionRuleSource struct {
	ID           string
	SourceConfig *model.SecureCNConnectionRuleSource
}
