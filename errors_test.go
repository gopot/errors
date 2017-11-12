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

// It is a copy of https://golang.org/src/errors/errors_test.go with the following changes:
//  1. imported "errors" package replaced by "github.com/gopot/errors"
//  2. ExampleNew_errorf renamed into ExampleNewErrorf and uses
//	   github.com/gopot/errors.NewErrorf() function instead of original fmt.NewErrorf()
//
// Which is licensed with the following notice:
//
// Copyright 2011 The Go Authors. All rights reserved.
//
// Use of this source code is governed by https://golang.org/LICENSE License.

package errors_test

import (
	"fmt"
	"testing"

	"github.com/gopot/errors"
)

func TestNewEqual(t *testing.T) {
	// Different allocations should not be equal.
	if errors.New("abc") == errors.New("abc") {
		t.Errorf(`New("abc") == New("abc")`)
	}
	if errors.New("abc") == errors.New("xyz") {
		t.Errorf(`New("abc") == New("xyz")`)
	}

	// Same allocation should be equal to itself (not crash).
	err := errors.New("jkl")
	if err != err {
		t.Errorf(`err != err`)
	}
}

func TestErrorMethod(t *testing.T) {
	err := errors.New("abc")
	if err.Error() != "abc" {
		t.Errorf(`New("abc").Error() = %q, want %q`, err.Error(), "abc")
	}
}

func ExampleNew() {
	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Print(err)
	}
	// Output: emit macho dwarf: elf header corrupted
}

// The errors package's Errorf function lets us use fmt package's formatting
// features to create descriptive error messages.
func ExampleNewErrorf() {
	const name, id = "bimmler", 17
	err := errors.NewErrorf("user %q (id %d) not found", name, id)
	if err != nil {
		fmt.Print(err)
	}
	// Output: user "bimmler" (id 17) not found
}
