package servererrors

import (
	"fmt"
)

// ===============================================================================

type NotFoundError struct {
	Entity string
	Key    any
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with key '%v' not found", e.Entity, e.Key)
}

// ===============================================================================

type BadRequestError struct {
	Violation string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("submitted form violates server types: %s", e.Violation)
}

// ===============================================================================

type ForbiddenAccessError struct {
	UserId       string
	ResourceType string
	ResourceId   string
}

func (e *ForbiddenAccessError) Error() string {
	return fmt.Sprintf("user %s does not have access to %s %s", e.UserId, e.ResourceType, e.ResourceId)
}

// ===============================================================================

type FieldError struct {
	Field   string
	Message string
}

type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func (e *ValidationError) Add(field, msg string) {
	e.Fields = append(e.Fields, FieldError{Field: field, Message: msg})
}

// ===============================================================================

type AlreadyExistsError struct {
	Field string
	Value string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("'%s' with value '%s' already exists", e.Field, e.Value)
}

// ===============================================================================

type OverlapError struct {
	ApId string
}

func (e *OverlapError) Error() string {
	return fmt.Sprintf("Booking time overlap on apartment %s", e.ApId)
}