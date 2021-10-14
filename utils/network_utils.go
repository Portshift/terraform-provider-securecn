package utils

import (
	"terraform-provider-securecn/client"
	"terraform-provider-securecn/escher_api/escherClient"
)

func GetServiceApi(httpClientWrapper *client.HttpClientWrapper) *escherClient.MgmtServiceApiCtx {
	return httpClientWrapper.EscherClient
}
