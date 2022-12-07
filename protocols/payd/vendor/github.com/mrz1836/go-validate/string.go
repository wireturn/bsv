package validate

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// emailRegex is common regular expressions
var emailRegex = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?`)

// maxLengthStringValidation type used for string length values
type maxLengthStringValidation struct {
	// Validation is the validation interface
	Validation

	// length is the max string length value
	length int
}

// Validate is for the maxLengthStringValidation type and will test the max string length
func (m *maxLengthStringValidation) Validate(value interface{}, obj reflect.Value) *ValidationError {
	strValue, ok := value.(string)
	if !ok {
		return &ValidationError{
			Key:     m.FieldName(),
			Message: "is not of type string and MaxLengthValidation only accepts strings",
		}
	}

	if len(strValue) > m.length {
		return &ValidationError{
			Key:     m.FieldName(),
			Message: "must be no more than " + strconv.Itoa(m.length) + " characters",
		}
	}

	return nil
}

// minLengthStringValidation type used for string length values
type minLengthStringValidation struct {
	// Validation is the validation interface
	Validation

	// length is the max string length value
	length int
}

// Validate is for the minLengthStringValidation type and will test the min string length
func (m *minLengthStringValidation) Validate(value interface{}, obj reflect.Value) *ValidationError {
	strValue, ok := value.(string)
	if !ok {
		return &ValidationError{
			Key:     m.FieldName(),
			Message: "is not of type string and MinLengthValidation only accepts strings",
		}
	}

	if len(strValue) < m.length {
		return &ValidationError{
			Key:     m.FieldName(),
			Message: "must be at least " + strconv.Itoa(m.length) + " characters",
		}
	}

	return nil
}

// formatStringValidation type used for string pattern testing using regular expressions
type formatStringValidation struct {
	// Validation is the validation interface
	Validation

	// pattern is the regex pattern
	pattern *regexp.Regexp

	// patternName is the name of the pattern
	patternName string
}

// Validate is for the formatStringValidation type and will test the given regular expression
func (f *formatStringValidation) Validate(value interface{}, obj reflect.Value) *ValidationError {
	strValue, ok := value.(string)
	if !ok {
		return &ValidationError{
			Key:     f.FieldName(),
			Message: "is not of type string and FormatValidation only accepts strings",
		}
	}

	if !f.pattern.MatchString(strValue) {
		return &ValidationError{
			Key:     f.FieldName(),
			Message: "does not match " + f.patternName + " format",
		}
	}

	return nil
}

// stringEqualsString string equals string struct
type stringEqualsString struct {
	// Validation is the validation interface
	Validation

	// targetFieldName is the target field name to compare
	targetFieldName string
}

// Validate is for the stringEqualsString type and will test the given field's value and compare
func (s *stringEqualsString) Validate(value interface{}, obj reflect.Value) *ValidationError {

	strValue, ok := value.(string)
	if !ok {
		return &ValidationError{
			Key:     s.FieldName(),
			Message: "is not of type string and StringEqualsStringValidation only accepts strings",
		}
	}

	// Set field name
	compareField := obj.FieldByName(s.targetFieldName)

	// Try to set to string
	compareFieldStrValue, ok := compareField.Interface().(string)
	if !ok {
		return &ValidationError{
			Key:     s.targetFieldName,
			Message: "is not of type string and StringEqualsValidation only accepts strings",
		}
	}

	// Does not compare
	if strValue != compareFieldStrValue {
		return &ValidationError{
			Key:     s.FieldName(),
			Message: "is not the same as the compare field " + s.targetFieldName,
		}
	}

	return nil
}

// maxLengthValidation creates an interface based on the max length value
func maxLengthValidation(maxLength string, kind reflect.Kind) (Interface, error) {
	length, err := strconv.ParseInt(maxLength, 10, 0)
	if err != nil {
		return nil, err
	}

	return &maxLengthStringValidation{
		length: int(length),
	}, nil
}

// minLengthValidation creates an interface based on the minimum length value
func minLengthValidation(minLength string, kind reflect.Kind) (Interface, error) {
	length, err := strconv.ParseInt(minLength, 10, 0)
	if err != nil {
		return nil, err
	}

	return &minLengthStringValidation{
		length: int(length),
	}, nil
}

// formatValidation creates an interface based on the options
func formatValidation(options string, kind reflect.Kind) (Interface, error) {
	if strings.ToLower(options) == "email" {
		return &formatStringValidation{
			pattern:     emailRegex,
			patternName: "email",
		}, nil
	} else if strings.Contains(options, "regexp:") {
		patternStr := options[strings.Index(options, ":")+1:]
		pattern, err := regexp.Compile(patternStr)

		if err != nil {
			return nil, &ValidationError{Key: "regexp:", Message: err.Error()}
		}

		return &formatStringValidation{
			pattern:     pattern,
			patternName: "regexp",
		}, nil
	}

	return nil, &ValidationError{Key: "format", Message: "Has no pattern " + options}
}

// stringEqualsStringValidation creates an interface based on the field name
func stringEqualsStringValidation(fieldName string, kind reflect.Kind) (Interface, error) {
	return &stringEqualsString{
		targetFieldName: fieldName,
	}, nil
}

// init add the string validations when this package is loaded
func init() {

	// Max length validation is len(string) < X
	AddValidation("max_length", maxLengthValidation)

	// Min length validation is len(string) > X
	AddValidation("min_length", minLengthValidation)

	// Format validation uses a given regular expression to match
	AddValidation("format", formatValidation)

	// Compare validation uses another field to compare
	AddValidation("compare", stringEqualsStringValidation)
}
