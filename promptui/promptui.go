// Package promptui provides ui elements for the command line prompt.
package promptui

// ValidateFunc validates the given input. It should return a ValidationError
// if the input is not valid.
type ValidateFunc func(string) error

// ValidationError is the class of errors resulting from invalid inputs,
// returned from a ValidateFunc.
type ValidationError struct {
	msg string
}

// Error implements the error interface for ValidationError.
func (v *ValidationError) Error() string {
	return v.msg
}

// NewValidationError creates a new validation error with the given message.
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{msg: msg}
}