package eda

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNoValidateMethod    = errors.New("no Validate method")
	ErrValidateReturnValue = errors.New("validate method should return only an error")
)

type Validater interface {
	Validate() error
}

func Validate(v Validater) error {
	return validateOne(v)
}

func hasValidateMethod(v any) (reflect.Value, bool) {
	// get Validate method
	fn := reflect.ValueOf(v).MethodByName("Validate")
	return fn, fn.IsValid()
}

func validateOne(obj any) error {
	// get Validate method
	fn, hasValid := hasValidateMethod(obj)
	if !hasValid {
		return fmt.Errorf("%w in %T", ErrNoValidateMethod, obj)
	}
	// call Validate method
	values := fn.Call([]reflect.Value{})
	if len(values) != 1 {
		// should return only an error
		return ErrValidateReturnValue
	}
	value := values[0]
	// check the return value
	switch value.Interface().(type) {
	case error:
		return value.Interface().(error)
	case nil:
		return nil
	default:
		return ErrValidateReturnValue
	}
}

// ValidateAll calls Validate method on all elements of the slice
// TODO : check recursion infinite loop if called in Validate method
func ValidateAll(obj any) error {
	// get all members of v

	err := validateOne(obj)
	if !errors.Is(err, ErrNoValidateMethod) {
		return err
	}

	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	_ = objValue
	objKind := objType.Kind()

	t := fmt.Sprintf("%T", obj) // TODO : remove
	fmt.Println(t, objKind)

	switch objKind {
	case reflect.Struct:
		// get all fields of the struct
		nbFields := objType.NumField()
		for i := 0; i < nbFields; i++ {
			field := objType.Field(i)
			if field.IsExported() {
				obj := objValue.Field(i).Interface()
				_, hasValid := hasValidateMethod(obj)
				if hasValid {
					err := validateOne(obj)
					if err != nil {
						return err
					}
				} else {
					err := ValidateAll(obj)
					if err != nil {
						// TODO : wrap error here
						return err
					}
				}

			}
		}
	default:
		var fieldName string
		fieldName = fmt.Sprintf("%T", obj)
		_ = fieldName
		return fmt.Errorf("name=%s type=%s value=%s obj=%T : %w", fieldName, objType.String(), objValue.String(), obj, ErrNoValidateMethod)
	}

	return nil
}
