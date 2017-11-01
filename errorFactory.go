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

// Represents Default ErrorFactory which is used by package methods
var defaultErrorFactory = NewErrorFactory(NewDefaultKVStorage)

// Represents Error factory functionality.
type ErrorFactory interface {

	// Creates a new Error with `message` as `Error()` content and details.
	New(message string, details ...struct{ Key, Value interface{} }) Error

	// Converts basic built-in `error` into `Error`.
	ConvertToError(e error) Error
}

// Creates a default ErrorFactory which executes provided `detalizers` uppon each
// `ErrorFactory.New(...)` call. Then, detalizers results are appended to those which
// are passed to ErrorFactory.New(...) and fed to kvFactory to create back-end
// details storage.
//
// On ConvertToError(e) call default ErrorFactory works as following:
//  * in case error `e` implements `Error` - it returns `e` asserted to `Error`
//  * if not, it creates new `Error` with message as `e.Error()`
//
// Errors produced by default ErrorFactory are utilizing KVStorage interface:
//  * to Get -> kvStorage.GetValue
//  * to print Detailed -> kvStorage.String()
//
// NewErrorFactory(...) panics in case `kvFactory` function is nil.
func NewErrorFactory(kvFactory KVStorageFactory, detalizers ...Detalizer) ErrorFactory {
	if kvFactory == nil {
		panic("Won't instantiate ErrorFactory with nil KVStorageFactory.")
	}

	return &errorFactory{
		detalizers: detalizers,
		kvFactory:  kvFactory,
	}
}

// Internal default `ErrorFactory` implementation
type errorFactory struct {
	detalizers []Detalizer
	kvFactory  func(pairs ...struct{ Key, Value interface{} }) KVStorage
}

// Creates a new Error with `message` as `Error()` content and details appended with results of added Detalisers (see Add detalizer)
func (this *errorFactory) New(message string, details ...struct{ Key, Value interface{} }) Error {

	for _, detalizer := range this.detalizers {
		if detalizer != nil {
			details = append(details, detalizer()...)
		}
	}

	return &detailedError{
		error:        newBasicError(message),
		details:      this.kvFactory(details...),
		errorFactory: this,
	}
}

// Converts basic built-in `error` into `Error` as following:
// * in case error `e` implements `Error` - it returns `e` asserted to `Error`
// * if not, it creates new `Error` with message as `e.Error()`
func (this *errorFactory) ConvertToError(e error) Error {
	if dec, ok := e.(Error); ok {
		return dec
	}
	if e != nil {
		return this.New(e.Error())
	}
	return Error(nil)
}
