package graphql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

type DuplicateError struct {
	Field string
	Value string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("duplicate error: %s '%s' already exists", e.Field, e.Value)
}

type NotFoundError struct {
	Resource string
	ID       any
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s with ID %v does not exist", e.Resource, e.ID)
}

func wrapDBError(err error, field string, value string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return &DuplicateError{Field: field, Value: value}
	}
	return err
}

func wrapNotFoundError(err error, resource string, id any) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return &NotFoundError{Resource: resource, ID: id}
	}
	return err
}

func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

func IsDuplicateError(err error) bool {
	var de *DuplicateError
	return errors.As(err, &de)
}

func IsNotFoundError(err error) bool {
	var ne *NotFoundError
	return errors.As(err, &ne)
}

func validateRequiredString(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{Field: field, Message: "is required"}
	}
	return nil
}

func validatePositiveID(field string, id int) error {
	if id <= 0 {
		return &ValidationError{Field: field, Message: "must be a positive integer"}
	}
	return nil
}

func validateCreateBook(input CreateBookInput) error {
	if err := validateRequiredString("title", input.Title); err != nil {
		return err
	}
	if err := validatePositiveID("authorID", input.AuthorID); err != nil {
		return err
	}
	return nil
}

func validateUpdateBook(id int, input UpdateBookInput) error {
	if err := validatePositiveID("id", id); err != nil {
		return err
	}
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		return &ValidationError{Field: "title", Message: "cannot be empty"}
	}
	return nil
}

func validateCreateAuthor(input CreateAuthorInput) error {
	if err := validateRequiredString("name", input.Name); err != nil {
		return err
	}
	return nil
}

func validateUpdateAuthor(id int, input UpdateAuthorInput) error {
	if err := validatePositiveID("id", id); err != nil {
		return err
	}
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return &ValidationError{Field: "name", Message: "cannot be empty"}
	}
	return nil
}

func validateCreateSeries(input CreateSeriesInput) error {
	if err := validateRequiredString("name", input.Name); err != nil {
		return err
	}
	return nil
}

func validateUpdateSeries(id int, input UpdateSeriesInput) error {
	if err := validatePositiveID("id", id); err != nil {
		return err
	}
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return &ValidationError{Field: "name", Message: "cannot be empty"}
	}
	return nil
}

func validateCreateTag(name string) error {
	if err := validateRequiredString("name", name); err != nil {
		return err
	}
	return nil
}

func validatePositiveIDs(fields map[string]int) error {
	for field, id := range fields {
		if err := validatePositiveID(field, id); err != nil {
			return err
		}
	}
	return nil
}
