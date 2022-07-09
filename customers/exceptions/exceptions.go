package exceptions 

import (
	"errors"
	"fmt"
)

func ValidationError(invalidFields ...[]string) error {
	errorMessage := "Validation Error. Incorrect Fields: "
	for _, field := range invalidFields {
	errorMessage += fmt.Sprintf("%s", field)}
	return errors.New(errorMessage)
}

func UpdateFailure(reason ...string) error {
	return errors.New(
	fmt.Sprintf("Failed To Update Customer. Reason: %s", reason))
}	

func CreateFailure(reason ...string) error {
	return errors.New(
	fmt.Sprintf("Failed To Create Customer. Reason: %s", reason))
}

func DeleteFailure(reason ...string) error {
	return errors.New(
	fmt.Sprintf("Failed To Delete Customer. Reason: %s", reason))
}

