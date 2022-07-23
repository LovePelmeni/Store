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

func ValidationError() error {
	return errors.New("Validation Error.")
}

func FirebaseDatabaseFailure() error {
	return errors.New("Database Operation Failure")
}
