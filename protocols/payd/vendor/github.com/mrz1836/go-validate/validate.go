/*
Package validate (go-validate) provides validations for struct fields based on a validation tag and offers additional validation functions.
*/
package validate

import (
	"log"
	"reflect"
	"strings"
	"sync"
)

// DefaultMap is the default validation map used to tell if a struct is valid.
var DefaultMap = Map{}

// Interface specifies the necessary methods a validation must implement to be compatible with this package
type Interface interface {

	// SetFieldIndex stores the index of the field
	SetFieldIndex(index int)

	// FieldIndex retrieves the index of the field
	FieldIndex() int

	// SetFieldName stores the name of the field
	SetFieldName(name string)

	// FieldName retrieves the name of the field
	FieldName() string

	// Validate determines if the value is valid. The value nil is returned if it is valid
	Validate(value interface{}, obj reflect.Value) *ValidationError
}

// Validation is an implementation of a Interface and can be used to provide basic functionality to a new validation type through an anonymous field
type Validation struct {

	// Name of the validation
	Name string

	// fieldIndex is the field index location
	fieldIndex int

	// fieldName is the field name
	fieldName string

	// options for the validation
	options string
}

// SetFieldIndex stores the index of the field the validation was applied to
func (v *Validation) SetFieldIndex(index int) {
	v.fieldIndex = index
}

// FieldIndex retrieves the index of the field the validation was applied to
func (v *Validation) FieldIndex() int {
	return v.fieldIndex
}

// SetFieldName stores the name of the field the validation was applied to
func (v *Validation) SetFieldName(name string) {
	v.fieldName = name
}

// FieldName retrieves the name of the field the validation was applied to
func (v *Validation) FieldName() string {
	return v.fieldName
}

// Validate determines if the value is valid. The value nil is returned if it is valid
func (v *Validation) Validate(value interface{}, obj reflect.Value) *ValidationError {
	return &ValidationError{
		Key:     v.fieldName,
		Message: "validation not implemented",
	}
}

// Map is an atomic validation map and when two sets happen at the same time, the latest that started wins.
type Map struct {
	validator               sync.Map // map[reflect.Type][]Interface
	validationNameToBuilder sync.Map // map[string]func(string, reflect.Kind) (Interface, error)
}

// get will get the validator interface
func (m *Map) get(k reflect.Type) []Interface {
	v, ok := m.validator.Load(k)
	if !ok {
		return []Interface{}
	}
	return v.([]Interface)
}

// set will store the validator interface
func (m *Map) set(k reflect.Type, v []Interface) {
	m.validator.Store(k, v)
}

// AddValidation registers the validation specified by key to the known
// validations. If more than one validation registers with the same key, the
// last one will become the validation for that key.
func (m *Map) AddValidation(key string, fn func(string, reflect.Kind) (Interface, error)) {
	m.validationNameToBuilder.Store(key, fn)
}

// IsValid will either store the builder interfaces or run the IsValid based on the reflect object type
func (m *Map) IsValid(object interface{}) (bool, []ValidationError) {

	// Get the object's value and type
	objectValue := reflect.ValueOf(object)
	objectType := reflect.TypeOf(object)

	// Get the validations
	validations := m.get(objectType)

	// If we are a pointer and not nil, run IsValid on it's interface
	if objectValue.Kind() == reflect.Ptr && !objectValue.IsNil() {
		return IsValid(objectValue.Elem().Interface())
	}

	// Do we have some validations?
	if len(validations) == 0 {
		var err error

		// Loop the fields and decrement through the loop
		for i := objectType.NumField() - 1; i >= 0; i-- {
			field := objectType.Field(i)
			validationTag := field.Tag.Get("validation")

			// Do we have a validation tag?
			if len(validationTag) > 0 {
				validationComponent := strings.Split(validationTag, " ")

				// Loop each validation component
				for _, validationSpec := range validationComponent {
					component := strings.Split(validationSpec, "=")
					if len(component) != 2 {
						log.Fatalln("invalid validation specification:", objectType.Name(), field.Name, validationSpec)
					}

					// Create the validation
					var validation Interface
					if builder, ok := m.validationNameToBuilder.Load(component[0]); ok && builder != nil {
						fn := builder.(func(string, reflect.Kind) (Interface, error))
						validation, err = fn(component[1], field.Type.Kind())

						if err != nil {
							log.Fatalln("error creating validation:", objectType.Name(), field.Name, validationSpec, err)
						}
					} else {
						log.Fatalln("unknown validation named:", component[0])
					}

					// Store the other properties and append to validations
					validation.SetFieldName(field.Name)
					validation.SetFieldIndex(i)
					validations = append(validations, validation)
				}
			}
		}

		// Set the validations
		m.set(objectType, validations)
	}

	// Loop and build errors
	var errors []ValidationError
	for _, validation := range validations {
		field := objectValue.Field(validation.FieldIndex())
		value := field.Interface()
		if err := validation.Validate(value, objectValue); err != nil {
			errors = append(errors, *err)
		}
	}

	// Return flag and errors
	return len(errors) == 0, errors
}

// AddValidation registers the validation specified by key to the known
// validations. If more than one validation registers with the same key, the
// last one will become the validation for that key
// using DefaultMap.
func AddValidation(key string, fn func(string, reflect.Kind) (Interface, error)) {
	DefaultMap.AddValidation(key, fn)
}

// IsValid determines if an object is valid based on its validation tags using DefaultMap.
func IsValid(object interface{}) (bool, []ValidationError) {
	return DefaultMap.IsValid(object)
}
