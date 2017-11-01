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

// Represents predefined keys for specific details
const (
	// Represents Key for value of caused error.
	// Used to extract original error.
	// For more details see Error_causedOriginalError exmple.
	CausedByDetailKey = causedBy("Caused By")
)

type causedBy string

// New returns an Error that formats as the given text.
func New(text string) Error {
	return defaultErrorFactory.New(text)
}

// New returns an Error that formats as the given text with provided details.
func NewWithDetails(text string, details ...struct{ Key, Value interface{} }) Error {
	return defaultErrorFactory.New(text, details...)
}

// Errorf formats according to a format specifier and returns the string as a value that satisfies Error
func NewErrorf(format string, a ...interface{}) Error {
	return defaultErrorFactory.New(fmt.Sprintf(format, a...))
}

// Gently converts error `e` into Error trying to preserve original types and data
func ConvertToError(e error) Error {
	return defaultErrorFactory.ConvertToError(e)
}
