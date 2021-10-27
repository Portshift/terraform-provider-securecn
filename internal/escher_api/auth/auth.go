package auth

import (
	"errors"
	"net/http"

	"github.com/EscherAuth/escher/config"
	escher_request "github.com/EscherAuth/escher/request"
	"github.com/EscherAuth/escher/signer"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

type HmacSha2Auth struct {
	Credentials     Credentials
	CredentialScope string
}

func NewAuth(accessKey string, sharedKey string, credentialScope string) HmacSha2Auth {
	staticCreds := Credentials{accessKey, sharedKey}
	return HmacSha2Auth{
		Credentials:     staticCreds,
		CredentialScope: credentialScope,
	}
}

func (auth HmacSha2Auth) AuthenticateRequest(r runtime.ClientRequest, _ strfmt.Registry) error {
	req := r.(*request)

	c := config.Config{}
	config.SetDefaults(&c)

	c.AccessKeyId = auth.Credentials.AccessKeyID
	c.ApiSecret = auth.Credentials.SecretAccessKey
	c.CredentialScope = auth.CredentialScope

	signerObj := signer.New(c)
	escherReq, err := escher_request.NewFromHTTPRequest(req.request)
	if err != nil {
		return err
	}

	signReq, err := signerObj.SignRequest(escherReq, []string{})
	if err != nil {
		return err
	}

	err = setHeader(req.request, signReq, c.GetAuthHeaderName())
	if err != nil {
		return err
	}

	err = setHeader(req.request, signReq, c.GetDateHeaderName())
	if err != nil {
		return err
	}

	return nil
}

func setHeader(req *http.Request, signReq *escher_request.Request, header string) error {
	authValue, success := signReq.Headers().Get(header)
	if success != true {
		return errors.New("could not get auth header")
	}

	req.Header.Set(header, authValue)
	return nil
}
