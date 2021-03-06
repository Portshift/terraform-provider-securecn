// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// DeleteCdRuleIDServerlessRuleReader is a Reader for the DeleteCdRuleIDServerlessRule structure.
type DeleteCdRuleIDServerlessRuleReader struct {
	Formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteCdRuleIDServerlessRuleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeleteCdRuleIDServerlessRuleNoContent()
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDeleteCdRuleIDServerlessRuleNoContent creates a DeleteCdRuleIDServerlessRuleNoContent with default headers values
func NewDeleteCdRuleIDServerlessRuleNoContent() *DeleteCdRuleIDServerlessRuleNoContent {
	return &DeleteCdRuleIDServerlessRuleNoContent{}
}

/* DeleteCdRuleIDServerlessRuleNoContent describes a response with status code 204, with default header values.

deleted
*/
type DeleteCdRuleIDServerlessRuleNoContent struct {
}

func (o *DeleteCdRuleIDServerlessRuleNoContent) Error() string {
	return fmt.Sprintf("[DELETE /cd/{ruleId}/serverlessRule][%d] deleteCdRuleIdServerlessRuleNoContent ", 204)
}

func (o *DeleteCdRuleIDServerlessRuleNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
