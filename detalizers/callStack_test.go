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

package detalizers_test

import (
	"fmt"
	"testing"

	"github.com/gopot/errors/detalizers"
)

func A() { B() }

func B() { C() }

func C() { D() }

var D = func() {}

func Test_NewCallStackDetalizerHappyPath(t *testing.T) {

	testCase := struct {
		TestAlias           string
		SkipFrames          int
		NestLevel           int
		ExpectedKey         interface{}
		ExpectedValueString string
	}{
		TestAlias:   `Valid detalizer with chain of calls`,
		SkipFrames:  0, // Should not take NewCallStackDetalizer closure into account
		NestLevel:   5, // Must take only test call
		ExpectedKey: detalizers.CallStackDetailKey,
		ExpectedValueString: fmt.Sprintln("") + fmt.Sprintln("\tgithub.com/gopot/errors/detalizers_test.Test_NewCallStackDetalizerHappyPath.func1.1 callStack_test.go:57") +
			fmt.Sprintln("\tgithub.com/gopot/errors/detalizers_test.C callStack_test.go:28") +
			fmt.Sprintln("\tgithub.com/gopot/errors/detalizers_test.B callStack_test.go:26") +
			fmt.Sprintln("\tgithub.com/gopot/errors/detalizers_test.A callStack_test.go:24") +
			fmt.Sprintln("\tgithub.com/gopot/errors/detalizers_test.Test_NewCallStackDetalizerHappyPath.func1 callStack_test.go:60"),
	}

	testFn := func(t *testing.T) {
		csd := detalizers.NewCallStackDetalizer(testCase.SkipFrames, testCase.NestLevel)
		details := []struct{ Key, Value interface{} }{}

		D = func() {
			details = csd()
		}

		A()

		if len(details) != 1 {
			t.Fatalf("Must be one and only one detail, while returned %d details", len(details))
		}

		detail := details[0]

		if detail.Key != testCase.ExpectedKey {
			t.Errorf("Returned Key: %#v\r\n while expected %#v", detail.Key, testCase.ExpectedKey)
		}

		if detail.Value == nil {
			t.Fatal("Returned detail with Value <nil>")
		}

		actualValueStringer, ok := (detail.Value).(fmt.Stringer)
		if !ok {
			t.Fatal("Returned Value not assertable to fmt.Stringer")
		}

		actualValueString := actualValueStringer.String()
		if actualValueString != testCase.ExpectedValueString {
			t.Errorf("Returned Value.String() '%s' \r\n while expected '%s'", actualValueString, testCase.ExpectedValueString)
		}
	}

	t.Run(testCase.TestAlias, testFn)

}

func Test_NewCallStackDetalizerTooMuchToSkip(t *testing.T) {

	t.Parallel()

	testCase := struct {
		TestAlias       string
		SkipFrames      int
		NestLevel       int
		ExpectedDetails []struct{ Key, Value interface{} }
	}{
		TestAlias:       `Detalizer skips everything`,
		SkipFrames:      1000, // Should be arbitrary high
		NestLevel:       1024,
		ExpectedDetails: nil,
	}

	testFn := func(t *testing.T) {
		csd := detalizers.NewCallStackDetalizer(testCase.SkipFrames, testCase.NestLevel)
		details := []struct{ Key, Value interface{} }{}

		D = func() {
			details = csd()
		}

		A()

		if details != nil {
			t.Fatalf("Returned detail %#v while expected to skip everything and return %#v", details, testCase.ExpectedDetails)
		}
	}

	t.Run(testCase.TestAlias, testFn)

}

func Test_CallStackDetailKeyStringer(t *testing.T) {

	const CALLSTACK_DETAIL_KEY_STRING = "Call Stack"

	var callStackDetailKeyStringer fmt.Stringer = detalizers.CallStackDetailKey

	if callStackDetailKeyStringer.String() != CALLSTACK_DETAIL_KEY_STRING {
		t.Errorf("detalizers.CallStackDetailKey.String() is %s \r\n while expected %s", callStackDetailKeyStringer.String(), CALLSTACK_DETAIL_KEY_STRING)
	}
}
