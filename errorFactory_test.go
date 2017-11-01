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
	"fmt"
	"reflect"
	"testing"

	"github.com/gopot/errors"
)

func Test_NewErrorFactoryWithNilKVFactory(t *testing.T) {

	expectedRecover := "Won't instantiate ErrorFactory with nil KVStorageFactory."

	defer func() {
		actualRecover := recover()

		if actualRecover != expectedRecover {
			t.Errorf("errors.NewErrorFactory(nil) panic %#v while expected %#v", actualRecover, expectedRecover)

		}
	}()

	_ = errors.NewErrorFactory(nil)

}

func Test_ErrorFactoryNewMethodPassDetailsToKVFactory(t *testing.T) {

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

	testCases := []struct {
		TestAlias     string
		PassedDetails []struct{ Key, Value interface{} }
		ExpectedPairs []struct{ Key, Value interface{} }
	}{
		{
			TestAlias:     `Nil details`,
			PassedDetails: nil,
			ExpectedPairs: nil,
		},
		{
			TestAlias: `1 detail with nil key and value`,
			PassedDetails: []struct{ Key, Value interface{} }{
				struct{ Key, Value interface{} }{nil, nil}},
			ExpectedPairs: []struct{ Key, Value interface{} }{
				struct{ Key, Value interface{} }{nil, nil}},
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
			},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		passedDetails := testCase.PassedDetails
		expectedPairs := testCase.ExpectedPairs

		testFn := func(t *testing.T) {

			var actualPairs []struct{ Key, Value interface{} }

			mockedKVStorageFactory := func(pairs ...struct{ Key, Value interface{} }) errors.KVStorage {
				actualPairs = pairs
				return nil
			}

			_ = errors.NewErrorFactory(mockedKVStorageFactory).New(`Something`, passedDetails...)

			if !reflect.DeepEqual(actualPairs, expectedPairs) {
				t.Errorf("errors.NewErrorFactory(mockedKVStorageFactory).New(`Something`,%#v) passed %#v to KVStorageFactory \r\n while expected %#v", passedDetails, actualPairs, expectedPairs)
			}
		}

		t.Run(testAlias, testFn)
	}

}

func Test_ErrorByFactoryGetCallsKVStorageGet(t *testing.T) {

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
			_, _ = errors.NewErrorFactory(kvFactory).New("Test").Get(getKey)

			if actualKey != expectedKey {
				t.Errorf("errors.NewErrorFactory(kvFactory).New(\"Test\").Get(%#v) \r\n called KVStorage.Get with %#v \r\n while expected to pass with %#v", getKey, actualKey, expectedKey)
			}
		}

		t.Run(testAlias, testFn)
	}

}

func Test_ErrorByFactoryGetReturnsResultFromKVStorageGet(t *testing.T) {

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
			actualValue, actualFound := errors.NewErrorFactory(kvFactory).New("Test").Get(getKey)

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

func Test_ErrorByFactoryDetailedReturnsResultFromKVStorageString(t *testing.T) {
	testCases := []struct {
		TestAlias        string
		ErrorMessage     string
		KVReturnString   string
		ExpectedDetailed string
	}{
		{
			TestAlias:        `Empty string`,
			ErrorMessage:     "Test",
			KVReturnString:   "",
			ExpectedDetailed: fmt.Sprintln("Test"),
		},
		{
			TestAlias:        `Something`,
			ErrorMessage:     "Test too",
			KVReturnString:   "Something",
			ExpectedDetailed: fmt.Sprintln("Test too") + "Something",
		},
		{
			TestAlias:        `Something with tabs and new lines`,
			ErrorMessage:     "Test too",
			KVReturnString:   "\t\t \r\n \t\tSomething\r\n",
			ExpectedDetailed: fmt.Sprintln("Test too") + "\t\t \r\n \t\tSomething\r\n",
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		errorMessage := testCase.ErrorMessage
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

			actualDetailed := errors.NewErrorFactory(kvFactory).New(errorMessage).Detailed()

			if actualDetailed != expectedDetailed {
				t.Errorf("errors.NewErrorFactory(kvFactory).New(%s).Detailed() \r\n returned string %#v \r\n while expected %#v", errorMessage, actualDetailed, expectedDetailed)
			}

		}

		t.Run(testAlias, testFn)
	}

}

func Test_ErrorFactoryConvertToError(t *testing.T) {

	testCases := []struct {
		TestAlias     string
		OriginalError error
		ExpectedError errors.Error
	}{
		{
			TestAlias:     `Nil Error`,
			OriginalError: nil,
			ExpectedError: nil,
		},
		{
			TestAlias:     `error(nil) to Error`,
			OriginalError: error(nil),
			ExpectedError: errors.Error(nil),
		},
		{
			TestAlias:     `Error(nil) to Error`,
			OriginalError: errors.Error(nil),
			ExpectedError: errors.Error(nil),
		},
		{
			TestAlias:     `e.New("Some error")`,
			OriginalError: e.New("Some error"),
			ExpectedError: errors.NewErrorFactory(errors.KVStorageFactory(errors.NewDefaultKVStorage)).New("Some error"),
		},
		{
			TestAlias:     `errors.New("Some Error")`,
			OriginalError: errors.NewErrorFactory(errors.KVStorageFactory(errors.NewDefaultKVStorage)).New("Some Error"),
			ExpectedError: errors.NewErrorFactory(errors.KVStorageFactory(errors.NewDefaultKVStorage)).New("Some Error"),
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		originalError := testCase.OriginalError
		expectedError := testCase.ExpectedError

		testFn := func(t *testing.T) {

			actualError := errors.NewErrorFactory(errors.KVStorageFactory(errors.NewDefaultKVStorage)).ConvertToError(originalError)
			if !testHelperErrorsAreEqual(actualError, expectedError) {
				t.Errorf("NewErrorFactory(KVStorageFactory(nil)).ConvertToError(%#v) \r\n returned %#v \r\n while expected %#v", originalError, actualError, expectedError)
			}
		}

		t.Run(testAlias, testFn)

	}

}

func Test_ErrorFactoryCallsDetalizers(t *testing.T) {

	// Dummy types
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}
	type t4 struct{}

	d1 := &t1{}
	d2 := &t2{}
	d3 := &t3{}
	d4 := &t4{}

	actuallyCalled := []interface{}{}

	testCases := []struct {
		TestAlias      string
		Detalizers     []errors.Detalizer
		ErrorMsg       string
		ExpectedCalled []interface{}
	}{
		{
			TestAlias: `nil detalizers`,

			Detalizers:     nil,
			ErrorMsg:       "TestError",
			ExpectedCalled: []interface{}{},
		},
		{
			TestAlias: `one detalizers`,

			Detalizers: []errors.Detalizer{
				func() []struct{ Key, Value interface{} } {
					actuallyCalled = append(actuallyCalled, d1)
					return []struct{ Key, Value interface{} }{}
				},
			},
			ErrorMsg:       "TestError",
			ExpectedCalled: []interface{}{d1},
		},
		{
			TestAlias: `three detalizers`,

			Detalizers: []errors.Detalizer{
				func() []struct{ Key, Value interface{} } {
					actuallyCalled = append(actuallyCalled, d2)
					return []struct{ Key, Value interface{} }{}
				},
				func() []struct{ Key, Value interface{} } {
					actuallyCalled = append(actuallyCalled, d3)
					return []struct{ Key, Value interface{} }{}
				},
				func() []struct{ Key, Value interface{} } {
					actuallyCalled = append(actuallyCalled, d4)
					return []struct{ Key, Value interface{} }{}
				},
			},
			ErrorMsg:       "TestError",
			ExpectedCalled: []interface{}{d2, d3, d4},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		detalizers := testCase.Detalizers
		errorMsg := testCase.ErrorMsg
		expectedCalled := testCase.ExpectedCalled

		actuallyCalled = []interface{}{}

		testFn := func(t *testing.T) {
			factory := errors.NewErrorFactory(errors.NewDefaultKVStorage, detalizers...)
			_ = factory.New(errorMsg)

			if len(actuallyCalled) != len(expectedCalled) {
				t.Fatalf("factory.New(%s) actually called %d times while expected %d", errorMsg, len(actuallyCalled), len(expectedCalled))
			}

			for i, ac := range actuallyCalled {
				if ac != expectedCalled[i] {
					t.Errorf("factory.New(%s) on iteration %d \r\n returned %#v \r\n while expected %#v", errorMsg, i, ac, expectedCalled[i])
				}
			}
		}

		t.Run(testAlias, testFn)
	}

}

func Test_ErrorFactoryPassDetalizersResultsToKVFactory(t *testing.T) {

	// Dummy types
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}
	type t4 struct{}
	type t5 struct{}
	type t6 struct{}

	// Keys
	k1 := &t1{}
	k2 := &t3{}
	k3 := &t5{}

	// Values
	v1 := &t2{}
	v2 := &t4{}
	v3 := &t6{}

	testCases := []struct {
		TestAlias       string
		Detalizers      []errors.Detalizer
		ErrorMsg        string
		ErrorDetails    []struct{ Key, Value interface{} }
		ExpectedDetails []struct{ Key, Value interface{} }
	}{
		{
			TestAlias:       `nil detalizers - no details`,
			Detalizers:      nil,
			ErrorMsg:        "TestError",
			ExpectedDetails: []struct{ Key, Value interface{} }{},
		},
		{
			TestAlias:       `nil detalizers - one detail`,
			Detalizers:      nil,
			ErrorMsg:        "TestError",
			ErrorDetails:    []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k1, v1}},
			ExpectedDetails: []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k1, v1}},
		},
		{
			TestAlias:       `nil detalizers - two detail`,
			Detalizers:      nil,
			ErrorMsg:        "TestError",
			ErrorDetails:    []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k1, v1}, struct{ Key, Value interface{} }{k2, v2}},
			ExpectedDetails: []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k1, v1}, struct{ Key, Value interface{} }{k2, v2}},
		}, {
			TestAlias: `one detalizers returns two details`,
			Detalizers: []errors.Detalizer{
				func() []struct{ Key, Value interface{} } {
					return []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k1, v1}, struct{ Key, Value interface{} }{k2, v2}}
				},
			},
			ErrorMsg:        "TestError",
			ExpectedDetails: []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k1, v1}, struct{ Key, Value interface{} }{k2, v2}},
		},
		{
			TestAlias: `two detalizers per one detail`,

			Detalizers: []errors.Detalizer{
				func() []struct{ Key, Value interface{} } {
					return []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k2, v2}}
				},
				func() []struct{ Key, Value interface{} } {
					return []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k1, v1}}
				},
			},
			ErrorMsg:        "TestError",
			ExpectedDetails: []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k2, v2}, struct{ Key, Value interface{} }{k1, v1}},
		},
		{
			TestAlias: `two detalizers : two and one detail`,

			Detalizers: []errors.Detalizer{
				func() []struct{ Key, Value interface{} } {
					return []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k3, v3}, struct{ Key, Value interface{} }{k2, v2}}
				},
				func() []struct{ Key, Value interface{} } {
					return []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k1, v1}}
				},
			},
			ErrorMsg:        "TestError",
			ExpectedDetails: []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k3, v3}, struct{ Key, Value interface{} }{k2, v2}, struct{ Key, Value interface{} }{k1, v1}},
		},
		{
			TestAlias: `two detalizers : one and two detail`,

			Detalizers: []errors.Detalizer{
				func() []struct{ Key, Value interface{} } {
					return []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k3, v3}}
				},
				func() []struct{ Key, Value interface{} } {
					return []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k2, v2}, struct{ Key, Value interface{} }{k1, v1}}
				},
			},
			ErrorMsg:        "TestError",
			ExpectedDetails: []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{k3, v3}, struct{ Key, Value interface{} }{k2, v2}, struct{ Key, Value interface{} }{k1, v1}},
		},
	}

	for _, testCase := range testCases {
		testAlias := testCase.TestAlias
		detalizers := testCase.Detalizers
		errorMsg := testCase.ErrorMsg
		errorDetails := testCase.ErrorDetails
		expectedDetails := testCase.ExpectedDetails

		testFn := func(t *testing.T) {

			var actualDetails []struct{ Key, Value interface{} }

			mockedKVStorageFactory := func(pairs ...struct{ Key, Value interface{} }) errors.KVStorage {
				actualDetails = pairs
				return nil
			}

			_ = errors.NewErrorFactory(mockedKVStorageFactory, detalizers...).New(errorMsg, errorDetails...)

			if len(actualDetails) != len(expectedDetails) {
				t.Fatalf("factory.New(%s) actually called %d times while expected %d", errorMsg, len(actualDetails), len(expectedDetails))
			}

			for i, ac := range actualDetails {
				if ac != expectedDetails[i] {
					t.Errorf("factory.New(%s) on iteration %d \r\n returned %#v \r\n while expected %#v", errorMsg, i, ac, expectedDetails[i])
				}
			}
		}

		t.Run(testAlias, testFn)
	}

}
