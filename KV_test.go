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
	"testing"

	"github.com/gopot/errors"
)

// Represents list of KV to be tested.
// All tests should loop over the list and run all test cases for each entry.
var (
	factoriesForTest = []struct {
		Alias     string
		KVFactory errors.KVStorageFactory
	}{
		{
			Alias:     "DefaultKVStorage",
			KVFactory: errors.NewDefaultKVStorage,
		},
	}
)

func Test_KVStorageGetValue(t *testing.T) {

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
		OriginalPairs []struct{ Key, Value interface{} }
		GetKey        interface{}
		ExpectedValue interface{}
		ExpectedOk    bool
	}{
		{
			TestAlias:     "Get from nil pairs",
			OriginalPairs: nil,
			GetKey:        k1,
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			TestAlias: "Get single existing entry",
			OriginalPairs: []struct{ Key, Value interface{} }{
				{
					Key:   k1,
					Value: v1,
				},
			},
			GetKey:        k1,
			ExpectedValue: v1,
			ExpectedOk:    true,
		},
		{
			TestAlias: "Get not existing entry",
			OriginalPairs: []struct{ Key, Value interface{} }{
				{
					Key:   k1,
					Value: v1,
				},
			},
			GetKey:        k2,
			ExpectedValue: nil,
			ExpectedOk:    false,
		},
		{
			TestAlias: "Get by key with multiple values",
			OriginalPairs: []struct{ Key, Value interface{} }{
				{
					Key:   k1,
					Value: v1,
				},
				{
					Key:   k1,
					Value: v2,
				},
			},
			GetKey:        k1,
			ExpectedValue: v2,
			ExpectedOk:    true,
		},
		{
			TestAlias: "Get by existing key = nil",
			OriginalPairs: []struct{ Key, Value interface{} }{
				{
					Key:   nil,
					Value: v1,
				},
				{
					Key:   k2,
					Value: v2,
				},
			},
			GetKey:        nil,
			ExpectedValue: v1,
			ExpectedOk:    true,
		},
	}

	for _, testItem := range factoriesForTest {

		testAliasPrefix := testItem.Alias
		testFactory := testItem.KVFactory

		for _, testCase := range testCases {
			testAlias := testAliasPrefix + "/" + testCase.TestAlias
			originalPairs := testCase.OriginalPairs
			getKey := testCase.GetKey
			expectedValue := testCase.ExpectedValue
			expectedOk := testCase.ExpectedOk

			testFn := func(t *testing.T) {
				kv := testFactory(originalPairs...)

				actualValue, actualOk := kv.GetValue(getKey)

				if actualValue != expectedValue {
					t.Errorf("testFactory(%#v).GetValue(%#v) \r\n returned value %#v \r\n while expected %#v", originalPairs, getKey, actualValue, expectedValue)
				}
				if actualOk != expectedOk {
					t.Errorf("testFactory(%#v).GetValue(%#v) \r\n returned ok as %#v \r\n while expected %#v", originalPairs, getKey, actualOk, expectedOk)
				}
			}

			t.Run(testAlias, testFn)

		}

	}

}

func Test_KVStorageString(t *testing.T) {

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
		TestAlias      string
		OriginalPairs  []struct{ Key, Value interface{} }
		ExpectedString string
	}{
		{
			TestAlias:      "nil entries",
			OriginalPairs:  []struct{ Key, Value interface{} }(nil),
			ExpectedString: "",
		},
		{
			TestAlias: "String single existing entry",
			OriginalPairs: []struct{ Key, Value interface{} }{
				{
					Key:   mockStringer{s: "Key1"},
					Value: mockStringer{s: "Value1"},
				},
			},
			ExpectedString: fmt.Sprintln("Key1", ":", "Value1"),
		},
		{
			TestAlias: "String mixed entries",
			OriginalPairs: []struct{ Key, Value interface{} }{
				{
					Key:   mockStringer{s: "Key1"},
					Value: v1,
				},
				{
					Key:   k1,
					Value: mockStringer{s: "Value1"},
				},
				{
					Key:   k2,
					Value: v2,
				},
				{
					Key:   mockStringer{s: "Key2"},
					Value: mockStringer{s: "Value2"},
				},
			},
			ExpectedString: fmt.Sprintln("Key2", ":", "Value2") + fmt.Sprintln("Value1") + fmt.Sprintln("Key1"),
		},
		{
			TestAlias: "String entries with the same stringer key",
			OriginalPairs: []struct{ Key, Value interface{} }{
				{
					Key:   mockStringer{s: "Key1"},
					Value: v1,
				},
				{
					Key:   mockStringer{s: "Key1"},
					Value: mockStringer{s: "Value1"},
				},
				{
					Key:   k2,
					Value: v2,
				},
				{
					Key:   mockStringer{s: "Key1"},
					Value: mockStringer{s: "Value2"},
				},
			},
			ExpectedString: fmt.Sprintln("Key1", ":", "Value2") + fmt.Sprintln("Key1", ":", "Value1") + fmt.Sprintln("Key1"),
		},
		{
			TestAlias: "String entries with the same non-stringer key",
			OriginalPairs: []struct{ Key, Value interface{} }{
				{
					Key:   k1,
					Value: v1,
				},
				{
					Key:   k1,
					Value: mockStringer{s: "Value1"},
				},
				{
					Key:   mockStringer{s: "Key2"},
					Value: v2,
				},
				{
					Key:   k1,
					Value: mockStringer{s: "Value2"},
				},
			},
			ExpectedString: fmt.Sprintln("Value2") + fmt.Sprintln("Key2") + fmt.Sprintln("Value1"),
		},
	}

	for _, testItem := range factoriesForTest {

		testAliasPrefix := testItem.Alias
		testFactory := testItem.KVFactory

		for _, testCase := range testCases {
			testAlias := testAliasPrefix + "/" + testCase.TestAlias
			originalPairs := testCase.OriginalPairs
			expectedString := testCase.ExpectedString

			testFn := func(t *testing.T) {
				kv := testFactory(originalPairs...)

				actualString := kv.String()

				if actualString != expectedString {
					t.Errorf("testFactory(%#v).String() \r\n returned string %#v \r\n while expected  %#v", originalPairs, actualString, expectedString)
				}
			}

			t.Run(testAlias, testFn)

		}

	}

}
