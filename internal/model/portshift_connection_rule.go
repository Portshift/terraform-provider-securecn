package model

import (
	"terraform-provider-securecn/internal/escher_api/model"
)

type ConnectionRule struct {
	ID         string
	RuleConfig *model.SecureCNConnectionRule
}
