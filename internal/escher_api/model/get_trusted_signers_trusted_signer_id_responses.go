// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// GetTrustedSignersTrustedSignerIDReader is a Reader for the GetTrustedSignersTrustedSignerID structure.
type GetTrustedSignersTrustedSignerIDReader struct {
	Formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetTrustedSignersTrustedSignerIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetTrustedSignersTrustedSignerIDOK()
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetTrustedSignersTrustedSignerIDOK creates a GetTrustedSignersTrustedSignerIDOK with default headers values
func NewGetTrustedSignersTrustedSignerIDOK() *GetTrustedSignersTrustedSignerIDOK {
	return &GetTrustedSignersTrustedSignerIDOK{}
}

/* GetTrustedSignersTrustedSignerIDOK describes a response with status code 200, with default header values.

Success
*/
type GetTrustedSignersTrustedSignerIDOK struct {
	Payload *TrustedSigner
}

func (o *GetTrustedSignersTrustedSignerIDOK) Error() string {
	return fmt.Sprintf("[GET /trustedSigners/{trustedSignerId}][%d] getTrustedSignersTrustedSignerIdOK  %+v", 200, o.Payload)
}
func (o *GetTrustedSignersTrustedSignerIDOK) GetPayload() *TrustedSigner {
	return o.Payload
}

func (o *GetTrustedSignersTrustedSignerIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(TrustedSigner)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
