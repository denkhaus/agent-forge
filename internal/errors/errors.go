// Package errors provides centralized error definitions and utilities for the MCP Planner application.
package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors defined at package level for consistent error handling.
var (
	// ErrNotFound indicates that a requested resource was not found.
	ErrNotFound = errors.New("not found")

	// ErrInvalidInput indicates that provided input is invalid or malformed.
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized indicates that the operation is not authorized.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrTimeout indicates that an operation timed out.
	ErrTimeout = errors.New("operation timed out")

	// ErrCircuitBreakerOpen indicates that the circuit breaker is open and rejecting requests.
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

	// ErrServiceUnavailable indicates that a service is temporarily unavailable.
	ErrServiceUnavailable = errors.New("service unavailable")

	// ErrConfigurationInvalid indicates that configuration is invalid.
	ErrConfigurationInvalid = errors.New("configuration invalid")

	// ErrAgentNotFound indicates that a requested agent was not found.
	ErrAgentNotFound = errors.New("agent not found")

	// ErrToolNotFound indicates that a requested tool was not found.
	ErrToolNotFound = errors.New("tool not found")

	// ErrPromptNotFound indicates that a requested prompt was not found.
	ErrPromptNotFound = errors.New("prompt not found")
)

// ValidationError represents an error that occurs during validation.
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s' with value '%v': %s", e.Field, e.Value, e.Message)
}

// NewValidationError creates a new validation error.
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// ConfigurationError represents an error that occurs during configuration loading or validation.
type ConfigurationError struct {
	Source  string
	Key     string
	Message string
}

// Error implements the error interface for ConfigurationError.
func (e *ConfigurationError) Error() string {
	return fmt.Sprintf("configuration error in %s for key '%s': %s", e.Source, e.Key, e.Message)
}

// NewConfigurationError creates a new configuration error.
func NewConfigurationError(source, key, message string) *ConfigurationError {
	return &ConfigurationError{
		Source:  source,
		Key:     key,
		Message: message,
	}
}

// ProviderError represents an error that occurs in a provider.
type ProviderError struct {
	Provider  string
	Operation string
	Cause     error
}

// Error implements the error interface for ProviderError.
func (e *ProviderError) Error() string {
	return fmt.Sprintf("provider '%s' failed during '%s': %v", e.Provider, e.Operation, e.Cause)
}

// Unwrap returns the underlying cause of the provider error.
func (e *ProviderError) Unwrap() error {
	return e.Cause
}

// NewProviderError creates a new provider error.
func NewProviderError(provider, operation string, cause error) *ProviderError {
	return &ProviderError{
		Provider:  provider,
		Operation: operation,
		Cause:     cause,
	}
}

// IsNotFound checks if an error is a "not found" error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrAgentNotFound) ||
		errors.Is(err, ErrToolNotFound) ||
		errors.Is(err, ErrPromptNotFound)
}

// IsValidation checks if an error is a validation error.
func IsValidation(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsConfiguration checks if an error is a configuration error.
func IsConfiguration(err error) bool {
	var configErr *ConfigurationError
	return errors.As(err, &configErr)
}

// IsProvider checks if an error is a provider error.
func IsProvider(err error) bool {
	var providerErr *ProviderError
	return errors.As(err, &providerErr)
}
