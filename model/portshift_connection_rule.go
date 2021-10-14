package model

import "terraform-provider-securecn/escher_api/model"

type ConnectionRule struct {
	ID         string
	RuleConfig *model.SecureCNConnectionRule
}
