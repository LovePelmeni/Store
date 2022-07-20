package exceptions

import (
	"errors"
	"fmt"
)

func ValidationError(InvalidData ...[]string) error {
	return errors.New(fmt.Sprintf("ValidationError. Invalid Data: %v", InvalidData))
}

func OperationFailedError(OperationName ...string) error {
	return errors.New(fmt.Sprintf("Operational Error. Operation: %s Failed.", OperationName))
}

func DatabaseOperationFailure() error {
	return errors.New("Database Operation Failed...")
}

func DatabaseRecordNotFound() error {
	return errors.New("Database Record Not Found...")
}

func ServiceUnavailable(ServiceName ...string) error {
	return errors.New(fmt.Sprintf("Service `%s` is Unavailable.", ServiceName))
}
