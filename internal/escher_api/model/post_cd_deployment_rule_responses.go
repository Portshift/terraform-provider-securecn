// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// PostCdDeploymentRuleReader is a Reader for the PostCdDeploymentRule structure.
type PostCdDeploymentRuleReader struct {
	Formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostCdDeploymentRuleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPostCdDeploymentRuleCreated()
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPostCdDeploymentRuleCreated creates a PostCdDeploymentRuleCreated with default headers values
func NewPostCdDeploymentRuleCreated() *PostCdDeploymentRuleCreated {
	return &PostCdDeploymentRuleCreated{}
}

/*PostCdDeploymentRuleCreated handles this case with default header values.

created.
*/
type PostCdDeploymentRuleCreated struct {
	Payload *CdAppRule
}

func (o *PostCdDeploymentRuleCreated) Error() string {
	return fmt.Sprintf("[POST /cd/deploymentRule][%d] postCdDeploymentRuleCreated  %+v", 201, o.Payload)
}

func (o *PostCdDeploymentRuleCreated) GetPayload() *CdAppRule {
	return o.Payload
}

func (o *PostCdDeploymentRuleCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(CdAppRule)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
