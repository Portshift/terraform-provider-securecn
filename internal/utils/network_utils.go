package utils

import (
	"terraform-provider-securecn/internal/client"
	"terraform-provider-securecn/internal/escher_api/escherClient"
)

func GetServiceApi(httpClientWrapper *client.HttpClientWrapper) *escherClient.MgmtServiceApiCtx {
	return httpClientWrapper.EscherClient
}
