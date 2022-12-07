package validate

import (
	"reflect"
	"strconv"
)

// intValueValidation type used for integer values
type intValueValidation struct {
	// Validation is the validation interface
	Validation

	// value is a testing value to compare
	value int64

	// less is a boolean for determining if less (min) or not (max)
	less bool
}

// Validate is for the intValueValidation type and will compare the integer value (min/max)
func (i *intValueValidation) Validate(value interface{}, obj reflect.Value) *ValidationError {

	// Compare the value to see if it is convertible to type int64
	var compareValue int64
	switch value := value.(type) {
	case int:
		compareValue = int64(value)
	case int8:
		compareValue = int64(value)
	case int16:
		compareValue = int64(value)
	case int32:
		compareValue = int64(value)
	case int64:
		compareValue = value
	default:
		return &ValidationError{
			Key:     i.FieldName(),
			Message: "is not convertible to type int64",
		}
	}

	// Check min
	if i.less {
		if compareValue < i.value {
			return &ValidationError{
				Key:     i.FieldName(),
				Message: "must be greater than or equal to " + strconv.FormatInt(i.value, 10),
			}
		}
	} else { // Check max
		if compareValue > i.value {
			return &ValidationError{
				Key:     i.FieldName(),
				Message: "must be less than or equal to " + strconv.FormatInt(i.value, 10),
			}
		}
	}

	return nil
}

// uintValueValidation type used for unsigned integer values
type uintValueValidation struct {

	// Validation is the validation interface
	Validation

	// value is a testing value to compare
	value uint64

	// less is a boolean for determining if less (min) or not (max)
	less bool
}

// Validate is for the uintValueValidation type and will compare the unsigned integer value (min/max)
func (u *uintValueValidation) Validate(value interface{}, obj reflect.Value) *ValidationError {

	// Compare the value to see if it is convertible to type int64
	var compareValue uint64
	switch value := value.(type) {
	case uint:
		compareValue = uint64(value)
	case uint8:
		compareValue = uint64(value)
	case uint16:
		compareValue = uint64(value)
	case uint32:
		compareValue = uint64(value)
	case uint64:
		compareValue = value
	default:
		return &ValidationError{
			Key:     u.FieldName(),
			Message: "is not convertible to type uint64",
		}
	}

	// Check min
	if u.less {
		if compareValue < u.value {
			return &ValidationError{
				Key:     u.FieldName(),
				Message: "must be greater than or equal to " + strconv.FormatUint(u.value, 10),
			}
		}
	} else { // Check max
		if compareValue > u.value {
			return &ValidationError{
				Key:     u.FieldName(),
				Message: "must be less than or equal to " + strconv.FormatUint(u.value, 10),
			}
		}
	}

	return nil
}

// floatValueValidation type used for float values
type floatValueValidation struct {

	// Validation is the validation interface
	Validation

	// value is a testing value to compare
	value float64

	// less is a boolean for determining if less (min) or not (max)
	less bool
}

// Validate is for the floatValueValidation type and will compare the float value (min/max)
func (f *floatValueValidation) Validate(value interface{}, obj reflect.Value) *ValidationError {

	// Compare the value to see if it is convertible to type int64
	var compareValue float64
	switch value := value.(type) {
	case float32:
		compareValue = float64(value)
	case float64:
		compareValue = value
	default:
		return &ValidationError{
			Key:     f.FieldName(),
			Message: "is not convertible to type float64",
		}
	}

	// Check min
	if f.less {
		if compareValue < f.value {
			return &ValidationError{
				Key:     f.FieldName(),
				Message: "must be greater than or equal to " + strconv.FormatFloat(f.value, 'E', -1, 64),
			}
		}
	} else { // Check max
		if compareValue > f.value {
			return &ValidationError{
				Key:     f.FieldName(),
				Message: "must be less than or equal to " + strconv.FormatFloat(f.value, 'E', -1, 64),
			}
		}
	}

	return nil
}

// minValueValidation creates an interface based on the "kind" type
func minValueValidation(minValue string, kind reflect.Kind) (Interface, error) {
	switch kind {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		value, err := strconv.ParseInt(minValue, 10, 0)
		if err != nil {
			return nil, err
		}
		return &intValueValidation{
			value: value,
			less:  true,
		}, nil
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		value, err := strconv.ParseUint(minValue, 10, 0)
		if err != nil {
			return nil, err
		}
		return &uintValueValidation{
			value: value,
			less:  true,
		}, nil
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		value, err := strconv.ParseFloat(minValue, 64)
		if err != nil {
			return nil, err
		}
		return &floatValueValidation{
			value: value,
			less:  true,
		}, nil
	default:
		return nil, &ValidationError{
			Key:     "invalid_validation",
			Message: "field is not of numeric type and min validation only accepts numeric types",
		}
	}
}

// maxValueValidation creates an interface based on the "kind" type
func maxValueValidation(maxValue string, kind reflect.Kind) (Interface, error) {
	switch kind {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		value, err := strconv.ParseInt(maxValue, 10, 0)
		if err != nil {
			return nil, err
		}
		return &intValueValidation{
			value: value,
			less:  false,
		}, nil
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		value, err := strconv.ParseUint(maxValue, 10, 0)
		if err != nil {
			return nil, err
		}
		return &uintValueValidation{
			value: value,
			less:  false,
		}, nil
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		value, err := strconv.ParseFloat(maxValue, 64)
		if err != nil {
			return nil, err
		}
		return &floatValueValidation{
			value: value,
			less:  false,
		}, nil
	default:
		return nil, &ValidationError{
			Key:     "invalid_validation",
			Message: "field is not of numeric type and max validation only accepts numeric types",
		}
	}
}

// init add the numeric validations when this package is loaded
func init() {

	// Min validation is where X cannot be less then Y
	AddValidation("min", minValueValidation)

	// Max validation is where X cannot be greater than Y
	AddValidation("max", maxValueValidation)
}
