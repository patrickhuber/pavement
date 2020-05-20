package types

import "strings"

// AggregateError provides an error struct that aggregates child errors into a single structure
type AggregateError interface {
	error
	Errors() []error
	Append(e error)
	Length() int
	IsEmpty() bool
	Join(e AggregateError) AggregateError
}

type aggregateError struct {
	errors []error
}

// Append appends the given error to the aggregate error
func (ae *aggregateError) Append(e error) {
	ae.errors = append(ae.errors, e)
}

// Join combines two aggregate errors into a new single aggregate error
func (ae *aggregateError) Join(e AggregateError) AggregateError {
	if ae.IsEmpty() {
		return e
	}
	if e.IsEmpty() {
		return ae
	}
	errors := []error{}
	errors = append(errors, ae.errors...)
	errors = append(errors, e.Errors()...)
	return &aggregateError{
		errors: errors,
	}
}

// Length returns the length of the error array
func (ae *aggregateError) Length() int {
	return len(ae.errors)
}

// IsEmpty returns true if the there are no errors contained in this aggregate
func (ae *aggregateError) IsEmpty() bool {
	return ae.Length() == 0
}

// Error implements the error interface
func (ae *aggregateError) Error() string {
	errorStrings := []string{}
	for _, e := range ae.errors {
		errorStrings = append(errorStrings, e.Error())
	}
	return strings.Join(errorStrings, "\n")
}

// Errors returns the underlying error slice
func (ae *aggregateError) Errors() []error {
	return ae.errors
}

// NewAggregateError returns anew aggregate error
func NewAggregateError(e ...error) AggregateError {
	ae := &aggregateError{
		errors: []error{},
	}
	if e == nil || len(e) == 0 {
		return ae
	}
	ae.errors = append(ae.errors, e...)
	return ae
}
