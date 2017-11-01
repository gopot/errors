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

package errors

import "fmt"

// Represents standard(first class) `error` interface extended with detalization methods
// returning context-like information for easier indentification of specific details.
type DetailedError interface {

	// Implements standard built-in `error` interface.
	error

	// Returns string representation of detailed error's information.
	Detailed() string

	// Returns value of the detail with given `key` and boolean flag wheather such detail was found.
	Get(key interface{}) (interface{}, bool)
}

// Represents `DetailedError` interface extended with `Caused` method.
type Error interface {
	DetailedError

	// Returns a new `Error` as caused by current.
	Caused(string, ...struct{ Key, Value interface{} }) Error
}

// Internal default implementation of `Error`
type detailedError struct {
	error
	details      KVStorage
	errorFactory ErrorFactory
}

// Implements `error` interface and returns message of underliying error.
func (this *detailedError) Error() string {
	return this.error.Error()
}

// Returns error message followed by concatenation of following blocks:
// * if Detail `key` implements `Stringer` interface, the block starts with key.String() header, otherwise it is skept.
// * if Detail `value` implements `Stringer` interface, the block contains value.String(), otherwise it is skept.
func (this *detailedError) Detailed() string {
	return fmt.Sprintln(this.Error()) + this.details.String()
}

// Returns last added Detail `value` for provided `key` and boolean flag wheather the key was found.
func (this *detailedError) Get(key interface{}) (interface{}, bool) {
	return this.details.GetValue(key)
}

// Returns a new `Error` as caused by `this`.
// Current implementation appends details with `CausedByDetailKey` key and `this` Error as value.
func (this *detailedError) Caused(text string, details ...struct{ Key, Value interface{} }) Error {
	details = append(details, struct{ Key, Value interface{} }{CausedByDetailKey, this})

	return this.errorFactory.New(text+" caused by: "+this.Error(), details...)
}
