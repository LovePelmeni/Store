package exceptions

import (
	"errors"
	"fmt"
)

func FailedRequest(reason ...error) error {
	return errors.New(fmt.Sprintf(
		"Failed to Send gRPC Email Request. Reason: %s", reason))
}

// Services Exceptions

func ServiceUnavailable() error {
	return errors.New("Service Unavailable...")
}
