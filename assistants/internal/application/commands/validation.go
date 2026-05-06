package commands

import "middleman/assistants/internal/models"

// Validatable marks a command that can self-validate before execution.
// Returning *models.ServiceError with Code="validation-error" allows
// consistent structured error handling for the LLM.
// When validation passes, return nil.
//
// NOTE: Only basic presence checks are implemented per command; business
// logic remains inside handlers / aggregates.
//
// Handlers can test `cmd, ok := any(cmd).(Validatable)` and call Validate().
// The generic helper below does this automatically.
type Validatable interface {
	Validate() *models.ServiceError
}

// ValidateCommand invokes Validate() if the command implements Validatable.
// It returns a ServiceError or nil.
func ValidateCommand(cmd interface{}) *models.ServiceError {
	if v, ok := cmd.(Validatable); ok {
		return v.Validate()
	}
	return nil
}
