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
	"fmt"
	"reflect"
	"testing"

	"github.com/gopot/errors"
)

// Represents list of ErrorFactory(ies) used to produce Error to be tested
var testFactories = []struct {
	Alias        string
	ErrorFactory errors.ErrorFactory
}{
	{
		Alias:        `Default ErrorFactory with Default KVStorage`,
		ErrorFactory: errors.NewErrorFactory(errors.NewDefaultKVStorage),
	},
	{
		Alias:        `Default ErrorFactory with <nil> KVStorage`,
		ErrorFactory: errors.NewErrorFactory(func(pairs ...struct{ Key, Value interface{} }) errors.KVStorage { return nil }),
	},
}

func Test_ErrorFactoryNewError(t *testing.T) {

	testCases := []struct {
		TestAlias           string
		InputMessage        string
		InputDetails        []struct{ Key, Value interface{} }
		ExpectedErrorOutput string
	}{
		{
			TestAlias:           `nil details`,
			InputMessage:        "TestError",
			InputDetails:        nil,
			ExpectedErrorOutput: "TestError",
		},
		{
			TestAlias:           `empty slice details`,
			InputMessage:        "TestError",
			InputDetails:        []struct{ Key, Value interface{} }{},
			ExpectedErrorOutput: "TestError",
		},
		{
			TestAlias:           `{nil,nil} details`,
			InputMessage:        "TestError",
			InputDetails:        []struct{ Key, Value interface{} }{{nil, nil}},
			ExpectedErrorOutput: "TestError",
		},
		{
			TestAlias:           `{Stringer,Stringer} details`,
			InputMessage:        "TestError",
			InputDetails:        []struct{ Key, Value interface{} }{{mockStringer{s: "Key"}, mockStringer{s: "Value"}}},
			ExpectedErrorOutput: "TestError",
		},
	}

	for _, testFactory := range testFactories {
		for _, testCase := range testCases {
			testAlias := testFactory.Alias + "/" + testCase.TestAlias
			inputMessage := testCase.InputMessage
			inputDetails := testCase.InputDetails
			expectedErrorOutput := testCase.ExpectedErrorOutput

			err := testFactory.ErrorFactory.New(inputMessage, inputDetails...)

			testFn := func(t *testing.T) {

				actualErrorOutput := err.Error()

				if actualErrorOutput != expectedErrorOutput {
					t.Errorf("ErrorFactory(%s, %#v) returned \r\n %s \r\n while expected \r\n %s \r\n", inputMessage, inputDetails, actualErrorOutput, expectedErrorOutput)
				}

			}

			t.Run(testAlias, testFn)
		}
	}

}

// ====================================== //
//                                        //
//          Error.Caused(...)             //
//                                        //
// ====================================== //

func Test_CausedError(t *testing.T) {

	testCases := []struct {
		TestAlias            string
		OriginalErrorMessage string
		OriginalErrorDetails []struct{ Key, Value interface{} }
		InputMessage         string
		InputDetails         []struct{ Key, Value interface{} }
		ExpectedErrorOutput  string
	}{
		{
			TestAlias:            `nil details`,
			OriginalErrorMessage: "Original error",
			OriginalErrorDetails: nil,
			InputMessage:         "TestError",
			InputDetails:         nil,
			ExpectedErrorOutput:  "TestError caused by: Original error",
		},
		{
			TestAlias:            `empty slice details`,
			OriginalErrorMessage: "Original error",
			OriginalErrorDetails: nil,
			InputMessage:         "TestError",
			InputDetails:         []struct{ Key, Value interface{} }{},
			ExpectedErrorOutput:  "TestError caused by: Original error",
		},
		{
			TestAlias:            `{nil,nil} details`,
			OriginalErrorMessage: "Original error",
			OriginalErrorDetails: nil,
			InputMessage:         "TestError",
			InputDetails:         []struct{ Key, Value interface{} }{{nil, nil}},
			ExpectedErrorOutput:  "TestError caused by: Original error",
		},
		{
			TestAlias:            `{Stringer,Stringer} details`,
			OriginalErrorMessage: "Original error",
			OriginalErrorDetails: nil,
			InputMessage:         "TestError",
			InputDetails:         []struct{ Key, Value interface{} }{{mockStringer{s: "Key"}, mockStringer{s: "Value"}}},
			ExpectedErrorOutput:  "TestError caused by: Original error",
		},
	}

	for _, testFactory := range testFactories {
		for _, testCase := range testCases {
			testAlias := testFactory.Alias + "/" + testCase.TestAlias
			originalErrorMessage := testCase.OriginalErrorMessage
			originalErrorDetails := testCase.OriginalErrorDetails
			inputMessage := testCase.InputMessage
			inputDetails := testCase.InputDetails
			expectedErrorOutput := testCase.ExpectedErrorOutput

			originalError := testFactory.ErrorFactory.New(originalErrorMessage, originalErrorDetails...)

			testFn := func(t *testing.T) {
				actualErrorOutput := originalError.Caused(inputMessage, inputDetails...).Error()

				if actualErrorOutput != expectedErrorOutput {
					t.Errorf("ErrorFactory.New(%s, %#v).Caused(%s, %#v) returned \r\n %s \r\n while expected \r\n %s \r\n", originalErrorMessage, originalErrorDetails, inputMessage, inputDetails, actualErrorOutput, expectedErrorOutput)
				}

			}

			t.Run(testAlias, testFn)
		}
	}

}

func Test_ErrorCausedMethodPassDetailsToKVFactory(t *testing.T) {

	// Dummy types
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}
	type t4 struct{}

	// Keys
	k1 := &t1{}
	k2 := &t3{}

	// Values
	v1 := &t2{}
	v2 := &t4{}

	var actualPairs []struct{ Key, Value interface{} }

	mockedKVStorageFactory := func(pairs ...struct{ Key, Value interface{} }) errors.KVStorage {
		actualPairs = pairs
		return nil
	}

	originalError := errors.NewErrorFactory(mockedKVStorageFactory).New(`Something`)

	testCases := []struct {
		TestAlias     string
		PassedDetails []struct{ Key, Value interface{} }
		ExpectedPairs []struct{ Key, Value interface{} }
	}{
		{
			TestAlias:     `Nil details`,
			PassedDetails: nil,
			ExpectedPairs: []struct{ Key, Value interface{} }{
				struct{ Key, Value interface{} }{errors.CausedByDetailKey, originalError},
			},
		},
		{
			TestAlias: `1 detail with nil key and value`,
			PassedDetails: []struct{ Key, Value interface{} }{
				struct{ Key, Value interface{} }{nil, nil}},
			ExpectedPairs: []struct{ Key, Value interface{} }{
				struct{ Key, Value interface{} }{nil, nil},
				struct{ Key, Value interface{} }{errors.CausedByDetailKey, originalError},
			},
		},
		{
			TestAlias: `3 detail with nil key and value`,
			PassedDetails: []struct{ Key, Value interface{} }{
				struct{ Key, Value interface{} }{nil, nil},
				struct{ Key, Value interface{} }{nil, nil},
				struct{ Key, Value interface{} }{nil, nil},
			},
			ExpectedPairs: []struct{ Key, Value interface{} }{
				struct{ Key, Value interface{} }{nil, nil},
				struct{ Key, Value interface{} }{nil, nil},
				struct{ Key, Value interface{} }{nil, nil},
				struct{ Key, Value interface{} }{errors.CausedByDetailKey, originalError},
			},
		},
		{
			TestAlias: `2 detail with not nil key and value`,
			PassedDetails: []struct{ Key, Value interface{} }{
				struct{ Key, Value interface{} }{k1, v1},
				struct{ Key, Value interface{} }{k2, v2},
			},
			ExpectedPairs: []struct{ Key, Value interface{} }{
				struct{ Key, Value interface{} }{k1, v1},
				struct{ Key, Value interface{} }{k2, v2},
				struct{ Key, Value interface{} }{errors.CausedByDetailKey, originalError},
			},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		passedDetails := testCase.PassedDetails
		expectedPairs := testCase.ExpectedPairs

		testFn := func(t *testing.T) {

			// refresh actual pairs
			actualPairs = nil

			_ = originalError.Caused("Caused", passedDetails...)

			if !reflect.DeepEqual(actualPairs, expectedPairs) {
				t.Errorf("factory.New(`Something`).Caused(\"Caused\", %#v) \r\n passed details %#v to KVStorageFactory \r\n while expected %#v", passedDetails, actualPairs, expectedPairs)
			}
		}

		t.Run(testAlias, testFn)
	}

}

func Test_CausedErrorGetCallsKVStorageGet(t *testing.T) {

	// Dummy types
	type t1 struct{}

	// Keys
	k1 := &t1{}

	testCases := []struct {
		TestAlias   string
		GetKey      interface{}
		ExpectedKey interface{}
	}{
		{
			TestAlias:   `nil`,
			GetKey:      nil,
			ExpectedKey: nil,
		},
		{
			TestAlias:   `k1`,
			GetKey:      k1,
			ExpectedKey: k1,
		},
		{
			TestAlias:   `error(nil)`,
			GetKey:      error(nil),
			ExpectedKey: error(nil),
		},
		{
			TestAlias:   `Caused`,
			GetKey:      errors.CausedByDetailKey,
			ExpectedKey: errors.CausedByDetailKey,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		getKey := testCase.GetKey
		expectedKey := testCase.ExpectedKey

		var actualKey interface{}

		kvFactory := func(pairs ...struct{ Key, Value interface{} }) errors.KVStorage {
			return &mockKVStorage{
				getValue: func(key interface{}) (value interface{}, found bool) {
					actualKey = key
					return nil, false
				},
			}
		}

		testFn := func(t *testing.T) {
			_, _ = errors.NewErrorFactory(kvFactory).New("Test").Caused("Caused").Get(getKey)

			if actualKey != expectedKey {
				t.Errorf("errors.NewErrorFactory(kvFactory).New(\"Test\").Caused(\"Caused\").Get(%#v) \r\n called KVStorage.Get with %#v \r\n while expected to pass with %#v", getKey, actualKey, expectedKey)
			}
		}

		t.Run(testAlias, testFn)
	}

}

func Test_CausedErrorGetReturnsResultFromKVStorageGet(t *testing.T) {

	// Dummy types
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}
	type t4 struct{}

	// Keys
	k1 := &t1{}
	// k2 := &t3{}

	// Values
	v1 := &t2{}
	v2 := &t4{}

	testCases := []struct {
		TestAlias     string
		GetKey        interface{}
		KVReturnValue interface{}
		KVReturnFound bool
		ExpectedValue interface{}
		ExpectedFound bool
	}{
		{
			TestAlias:     `k1 - v1 - true`,
			GetKey:        k1,
			KVReturnValue: v1,
			KVReturnFound: true,
			ExpectedValue: v1,
			ExpectedFound: true,
		},
		{
			TestAlias:     `nil - v2 - true`,
			GetKey:        nil,
			KVReturnValue: v2,
			KVReturnFound: true,
			ExpectedValue: v2,
			ExpectedFound: true,
		},
		{
			TestAlias:     `k1 - nil - false`,
			GetKey:        k1,
			KVReturnValue: nil,
			KVReturnFound: false,
			ExpectedValue: nil,
			ExpectedFound: false,
		},
		{
			TestAlias:     `nil - v2 - false`,
			GetKey:        nil,
			KVReturnValue: v2,
			KVReturnFound: false,
			ExpectedValue: v2,
			ExpectedFound: false,
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		getKey := testCase.GetKey
		kVReturnValue := testCase.KVReturnValue
		kVReturnFound := testCase.KVReturnFound
		expectedValue := testCase.ExpectedValue
		expectedFound := testCase.ExpectedFound

		kvFactory := func(pairs ...struct{ Key, Value interface{} }) errors.KVStorage {
			return &mockKVStorage{
				getValue: func(key interface{}) (value interface{}, found bool) {
					return kVReturnValue, kVReturnFound
				},
			}
		}

		testFn := func(t *testing.T) {
			actualValue, actualFound := errors.NewErrorFactory(kvFactory).New("Test").Caused("Caused").Get(getKey)

			if actualValue != expectedValue {
				t.Errorf("errors.NewErrorFactory(kvFactory).New(\"Test\").Get(%#v) \r\n returned value %#v \r\n while expected %#v", getKey, actualValue, expectedValue)
			}

			if actualFound != expectedFound {
				t.Errorf("errors.NewErrorFactory(kvFactory).New(\"Test\").Get(%#v) \r\n returned found %#v \r\n while expected %#v", getKey, actualFound, expectedFound)
			}
		}

		t.Run(testAlias, testFn)
	}

}

func Test_CausedErrorDetailedReturnsResultFromKVStorageString(t *testing.T) {

	testCases := []struct {
		TestAlias            string
		OriginalErrorMessage string
		CausedErrorMessage   string
		KVReturnString       string
		ExpectedDetailed     string
	}{
		{
			TestAlias:            `Empty string`,
			OriginalErrorMessage: "OriginalError",
			CausedErrorMessage:   "Test",
			KVReturnString:       "",
			ExpectedDetailed:     fmt.Sprintln("Test caused by: OriginalError"),
		},
		{
			TestAlias:            `Something`,
			OriginalErrorMessage: "OriginalError",
			CausedErrorMessage:   "Test too",
			KVReturnString:       "Something",
			ExpectedDetailed:     fmt.Sprintln("Test too caused by: OriginalError") + "Something",
		},
		{
			TestAlias:            `Something with tabs and new lines`,
			OriginalErrorMessage: "OriginalError",
			CausedErrorMessage:   "Test too",
			KVReturnString:       "\t\t \r\n \t\tSomething\r\n",
			ExpectedDetailed:     fmt.Sprintln("Test too caused by: OriginalError") + "\t\t \r\n \t\tSomething\r\n",
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalErrorMessage := testCase.OriginalErrorMessage
		causedErrorMessage := testCase.CausedErrorMessage
		kVReturnString := testCase.KVReturnString
		expectedDetailed := testCase.ExpectedDetailed

		testFn := func(t *testing.T) {

			kvFactory := func(pairs ...struct{ Key, Value interface{} }) errors.KVStorage {
				return &mockKVStorage{
					stringFn: func() string {
						return kVReturnString
					},
				}
			}

			actualDetailed := errors.NewErrorFactory(kvFactory).New(originalErrorMessage).Caused(causedErrorMessage).Detailed()

			if actualDetailed != expectedDetailed {
				t.Errorf("errors.NewErrorFactory(kvFactory).New(%s).Caused(%s).Detailed() \r\n returned string %#v \r\n while expected %#v", originalErrorMessage, causedErrorMessage, actualDetailed, expectedDetailed)
			}

		}

		t.Run(testAlias, testFn)
	}

}
