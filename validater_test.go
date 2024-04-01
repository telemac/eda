package eda

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type HasValidateMethod1 struct{}

func (h HasValidateMethod1) Validate() error {
	return nil
}

type HasValidateMethod2 struct{}

func (h HasValidateMethod2) Validate() error {
	return nil
}

type HasBadValidateMethod struct{}

func (h HasBadValidateMethod) Validate() error {
	return errors.New("HasBadValidateMethod validation failed")
}

type MultipleValidateMethodsFails struct {
	HasValidateMethod1
	HasValidateMethod2
}

// Validate method for MultipleValidateMethodsFails
func (m MultipleValidateMethodsFails) Validate() error {
	return errors.New("MultipleValidateMethodsFails validation failed")
}

type MultipleValidateMethods struct {
	HasValidateMethod1
	HasValidateMethod2
}

type MultipleWithOneFialed struct {
	HasValidateMethod1
	HasBadValidateMethod
	HasValidateMethod2
}

type MultipleValidateMethodsUnknownField struct {
	HasValidateMethod1
	HasValidateMethod2
	UnknownField string
}

type MyGoodString string

func (m MyGoodString) Validate() error {
	return nil
}

type MyBadString string

func (m MyBadString) Validate() error {
	return errors.New("MyBadString validation failed")
}

type MyString string

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	var hasValidateMethod1 HasValidateMethod1
	err := Validate(hasValidateMethod1)
	assert.NoError(err)

	var hasValidateMethod2 HasValidateMethod2
	err = validateOne(hasValidateMethod2)
	assert.NoError(err)

	var multipleValidateMethodsFails MultipleValidateMethodsFails
	err = ValidateAll(multipleValidateMethodsFails)
	assert.Equal(errors.New("MultipleValidateMethodsFails validation failed"), err)

	var multipleValidateMethods MultipleValidateMethods
	err = validateOne(multipleValidateMethods)
	assert.ErrorIs(err, ErrNoValidateMethod)

	err = ValidateAll(multipleValidateMethods)
	assert.NoError(err)

	var multipleWithOneFialed MultipleWithOneFialed
	err = ValidateAll(multipleWithOneFialed)
	assert.Equal(errors.New("HasBadValidateMethod validation failed"), err)

	var myGoodString MyGoodString
	err = Validate(myGoodString)
	assert.NoError(err)

	var myBadString MyBadString
	err = Validate(myBadString)
	assert.Equal(errors.New("MyBadString validation failed"), err)

	var myString MyString
	err = ValidateAll(myString)
	assert.ErrorIs(err, ErrNoValidateMethod)

	var multipleValidateMethodsUnknownField MultipleValidateMethodsUnknownField
	err = ValidateAll(multipleValidateMethodsUnknownField)
	assert.ErrorIs(err, ErrNoValidateMethod)
	assert.Equal("name=string type=string value= obj=string : no Validate method", err.Error())
}
