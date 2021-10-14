package model

import "terraform-provider-securecn/escher_api/model"

type ConnectionRuleSource struct {
	ID           string
	SourceConfig *model.SecureCNConnectionRuleSource
}
