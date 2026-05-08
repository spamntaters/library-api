package graphql

import (
	"errors"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
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

func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
