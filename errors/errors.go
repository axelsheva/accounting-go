// Package errors defines custom errors for the application.
package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Standard error definitions
var (
	// ErrNotFound indicates that a requested resource could not be found.
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput indicates that the provided input is invalid.
	ErrInvalidInput = errors.New("invalid input provided")

	// ErrInsufficientFunds indicates that a user doesn't have enough balance for an operation.
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrDuplicateResource indicates that a resource already exists.
	ErrDuplicateResource = errors.New("resource already exists")

	// ErrUnauthorized indicates that the operation is not authorized.
	ErrUnauthorized = errors.New("unauthorized operation")

	// ErrInternal indicates an internal server error.
	ErrInternal = errors.New("internal server error")

	// ErrNegativeBalance indicates that an operation would result in a negative balance.
	ErrNegativeBalance = errors.New("operation would result in negative balance")
)

// WithDetails wraps an error with additional context information.
func WithDetails(err error, format string, args ...interface{}) error {
	details := fmt.Sprintf(format, args...)
	return fmt.Errorf("%w: %s", err, details)
}

// IsNotFound checks if the given error is an ErrNotFound error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsInvalidInput checks if the given error is an ErrInvalidInput error.
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsInsufficientFunds checks if the given error is an ErrInsufficientFunds error.
func IsInsufficientFunds(err error) bool {
	return errors.Is(err, ErrInsufficientFunds)
}

// IsDuplicateResource checks if the given error is an ErrDuplicateResource error.
func IsDuplicateResource(err error) bool {
	return errors.Is(err, ErrDuplicateResource)
}

// IsUnauthorized checks if the given error is an ErrUnauthorized error.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsInternal checks if the given error is an ErrInternal error.
func IsInternal(err error) bool {
	return errors.Is(err, ErrInternal)
}

// IsNegativeBalance checks if the given error is an ErrNegativeBalance error.
func IsNegativeBalance(err error) bool {
	return errors.Is(err, ErrNegativeBalance)
}

// IsNegativeBalanceConstraintError checks if the given error is related to the
// PostgreSQL constraint violation for a negative balance.
func IsNegativeBalanceConstraintError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	return strings.Contains(errMsg, "pq:") &&
		strings.Contains(errMsg, "relation \"balances\"") &&
		strings.Contains(errMsg, "violates check constraint \"balance_amount_non_negative\"")
}

// WrapNegativeBalanceError wraps a PostgreSQL constraint error into a more
// user-friendly ErrNegativeBalance error with additional details.
func WrapNegativeBalanceError(err error, details string) error {
	if IsNegativeBalanceConstraintError(err) {
		if details != "" {
			return WithDetails(ErrNegativeBalance, details)
		}
		return ErrNegativeBalance
	}
	return err
}
