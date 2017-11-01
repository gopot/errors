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
	"time"

	"github.com/gopot/errors"
)

func ExampleConvertToError_simply() {

	doSomething := func() error {
		return fmt.Errorf("Some error %s", "just for test")
	}

	err := doSomething()

	theError := errors.ConvertToError(err)

	// Operate on Error...

	fmt.Println(theError.Detailed())
	// Output: Some error just for test

}

func ExampleConvertToError_causedWithDetail() {

	doSomething := func() error {
		return fmt.Errorf("Some error %s", "just for test")
	}

	err := doSomething()

	// You could define `detail` type or inline it for each detail - both works fine.
	type detail struct {
		Key, Value interface{}
	}

	theError := errors.ConvertToError(err).Caused("Failure doing something", detail{"MyDetail", "My detail value ..."})

	fmt.Println("theError.Detailed:", "{", theError.Detailed(), "}")
	// Output: theError.Detailed: { Failure doing something caused by: Some error just for test
	// MyDetail : My detail value ...
	//  }

}

func ExampleError() {

	doSomething := func() error {
		return fmt.Errorf("Nothing %s", "to do")
	}

	err := doSomething()

	// You could define `detail` type or inline it for each detail - both works fine.
	theError := errors.ConvertToError(err).Caused(
		"Failure doing something",
		struct{ Key, Value interface{} }{"MyDetail", "My detail value ..."})

	fmt.Println("theError.Detailed:", "{", theError.Detailed(), "}")
	// Output: theError.Detailed: { Failure doing something caused by: Nothing to do
	// MyDetail : My detail value ...
	//  }
}

func ExampleError_causedOriginalError() {

	// Just to compare original error
	var ErrorDoingSomething = errors.New("Nothing to do")

	doSomething := func() errors.Error {
		return ErrorDoingSomething
	}

	err := doSomething()

	// You could define detail type or inline it for each detail - both works fine.
	type detail struct {
		Key, Value interface{}
	}

	theError := err.Caused("Failure doing something", detail{"MyDetail", "My detail value ..."})

	if originalErrorValue, found := theError.Get(errors.CausedByDetailKey); found {
		// Detail value is of type interface{}.
		// So, it should be asserted to Error
		originalErr := originalErrorValue.(errors.Error)

		fmt.Println("originalErrorValue == ErrorDoingSomething :", originalErrorValue == ErrorDoingSomething)

		// Actually, the original error is of type Error,
		// thus Detailed() could be printed.
		fmt.Println("originalErr.Detailed:", originalErr.Detailed())
	}
	// Output: originalErrorValue == ErrorDoingSomething : true
	// originalErr.Detailed: Nothing to do

}

func ExampleErrorFactory_checkForDetail() {

	const TimeDetailKey = "Time"

	// One Detalizer to add time context for each error created by following errorFactory.
	timeDetalizer := func() []struct{ Key, Value interface{} } {
		// To let Example pass test the result should be deterministic.
		// In real life you would like to return time.Now().
		myTime, _ := time.Parse(time.RFC3339, "2017-11-01T00:00:05+02:00")
		return []struct{ Key, Value interface{} }{{TimeDetailKey, myTime}}
	}

	errorFactory := errors.NewErrorFactory(errors.NewDefaultKVStorage, timeDetalizer)

	const hasNothingToDoKey = "Has nothing to do"

	doSomething := func() errors.Error {
		// In this case boolean value for flag true/false `hasNothingToDoKey` is provided explicitly.
		// However, it is possible to ommit it and define true as existance of such detail.
		return errorFactory.New("Nothing to do", struct{ Key, Value interface{} }{Key: hasNothingToDoKey, Value: true})
	}

	err := doSomething()
	if err != nil {
		// Check if has nothing to do
		if valueInterface, found := err.Get(hasNothingToDoKey); found {

			// To make code less busy, it is recommended for flag(s) true/false to
			// imply value as existance of the key. So, here it could havebeen enough
			// to assume `found` as `value` and skip next if{}.
			if value, ok := valueInterface.(bool); ok {
				fmt.Println(hasNothingToDoKey, "is", value)
				// Give it something to do ;)
			}
		}

		fmt.Println("Error Detailed:", "{", err.Detailed(), "}")
	}
	// Output: Has nothing to do is true
	// Error Detailed: { Nothing to do
	// Time : 2017-11-01 00:00:05 +0200 +0200
	// Has nothing to do
	//  }

}
