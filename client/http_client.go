package client

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"terraform-provider-securecn/escher_api/escherClient"
)

type HttpClientWrapper struct {
	AccessKey    string
	SecretKey    string
	BaseURL      string
	EscherClient *escherClient.MgmtServiceApiCtx
	HttpClient   *http.Client
}

func NewHttpClient(accessKey, secretKey, baseUrl string) HttpClientWrapper {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		}}

	serviceApi, err := escherClient.CreateServiceApi(baseUrl, accessKey, secretKey, httpClient)

	if err != nil {
		log.Print("[ERROR] failed to initialize escher client")
		os.Exit(1)
	}

	return HttpClientWrapper{
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		BaseURL:      baseUrl,
		EscherClient: serviceApi,
		HttpClient:   httpClient,
	}

}
