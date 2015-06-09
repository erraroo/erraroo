package models

import "fmt"

// ValidationErrors represents api request errors
type ValidationErrors struct {
	Errors map[string][]string
}

// NewValidationErrors returns an empty ValidationErrors
func NewValidationErrors() ValidationErrors {
	return ValidationErrors{
		make(map[string][]string),
	}
}

func (r ValidationErrors) Error() string {
	return fmt.Sprintf("ValidationErrors: %v", r.Errors)
}

// Add adds a new error
func (r *ValidationErrors) Add(name, message string) {
	if r.Errors[name] == nil {
		r.Errors[name] = []string{}
	}

	r.Errors[name] = append(r.Errors[name], message)
}

// Any returns true if contains any errors
func (r ValidationErrors) Any() bool {
	return len(r.Errors) > 0
}
