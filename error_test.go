//   Copyright © 2015-2017 Ivan Kostko (github.com/ivan-kostko; github.com/gopot)

//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at

//       http://www.apache.org/licenses/LICENSE-2.0

//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package errors_test

import (
	e "errors"
	"testing"

	"github.com/gopot/errors"
)

func TestNewErrorImplementserror(t *testing.T) {
	var _ error = errors.New("")
}

func TestNewWithDetails(t *testing.T) {
	// Given
	type myPrivateInt int
	type myPrivateStr string

	errorMessage := "Some message"

	key1 := myPrivateStr("Some string")
	val1 := myPrivateInt(1024)

	key2 := myPrivateInt(2048)

	details := []struct{ Key, Value interface{} }{{key1, val1}, {key2, nil}}

	// When
	err := errors.NewWithDetails(errorMessage, details...)

	// Then
	if err == nil {
		t.Fatalf("errors.NewWithDetails(%s, %#v...) returned nil while was not expected to do so", errorMessage, details)
	}

	if err.Error() != errorMessage {
		t.Errorf("Returned Error() %s \r\n while expected %s", err.Error(), errorMessage)
	}

	if valueInterface, found := err.Get(key1); found {
		if value, ok := valueInterface.(myPrivateInt); ok {
			if value != val1 {
				t.Errorf("Returned value is %#v \r\n while expected %#v", value, val1)
			}
		} else {
			t.Errorf("The valueInterface %#v is not assertable to myPrivateInt while it was expected to be", valueInterface)
		}
	} else {
		t.Errorf("The key %#v not found while expected to be", key1)
	}

	if _, found := err.Get(key2); !found {
		t.Errorf("The key %#v was not found while expected to be", key2)
	}
}

func TestErrorOtput(t *testing.T) {

	type CustomError struct {
		error
	}

	testCases := []struct {
		TestAlias     string
		OriginalError error
		ExpectedError string
	}{
		{
			TestAlias:     `Simple github.com/gopot/errors.New()`,
			OriginalError: errors.New("Test Error"),
			ExpectedError: "Test Error",
		},
		{
			TestAlias:     `Simple github.com/gopot/errors.NewErrorf()`,
			OriginalError: errors.NewErrorf("Test Error %s %#v", "some string", struct{ str string }{str: "Blah"}),
			ExpectedError: "Test Error some string struct { str string }{str:\"Blah\"}",
		},
		{
			TestAlias:     `gopot/errors from error`,
			OriginalError: errors.ConvertToError(e.New("Test error")),
			ExpectedError: "Test error",
		},
		{
			TestAlias:     `New Caused`,
			OriginalError: errors.New("Test Error").Caused("A new github.com/gopot/errors error"),
			ExpectedError: "A new github.com/gopot/errors error caused by: Test Error",
		},
		{
			TestAlias:     `Simple test`,
			OriginalError: errors.ConvertToError(e.New("Test error")).Caused("A new github.com/gopot/errors error"),
			ExpectedError: "A new github.com/gopot/errors error caused by: Test error",
		},
		{
			TestAlias:     `Nested CustomError`,
			OriginalError: errors.ConvertToError(CustomError{e.New("My error")}),
			ExpectedError: "My error",
		},
		{
			TestAlias:     `Nested CustomError`,
			OriginalError: errors.ConvertToError(CustomError{e.New("My error")}).Caused("A new github.com/gopot/errors error"),
			ExpectedError: "A new github.com/gopot/errors error caused by: My error",
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalError := testCase.OriginalError
		expectedError := testCase.ExpectedError

		testFn := func(t *testing.T) {

			actualError := originalError.Error()
			if actualError != expectedError {
				t.Errorf("%s :: %#v.Error() returned \r\n %s \r\n while expected \r\n %s \r\n", testAlias, originalError, actualError, expectedError)
			}

		}

		t.Run(testAlias, testFn)
	}

}
