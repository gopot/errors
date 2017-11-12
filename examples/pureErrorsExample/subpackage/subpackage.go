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

package subpackage

import (
	"github.com/gopot/errors"
)

// Represents Error keys identifiers
const (
	// This definition would allow clients to make a check as following:
	// 	err.Get("Is Retriable")
	ErrorIsRetriable = string("Is Retriable")

	// This definition will push clients to check only for
	//	err.Get(subpackage.ErrorIsCritical).
	// This approach is more secure and strongly recommended.
	// Even it brings from 2 to 4 extra lines of code.
	ErrorIsCritical = isCriticalErrorDetailKey("Is Critical")
)

type isCriticalErrorDetailKey string

func (this isCriticalErrorDetailKey) String() string { return string(this) }

// Make this type private to uniquely identify keys
type subpackageErrorKeyStringer string

func ReturnsRetriableError() errors.Error {
	return internalFunctionReturningJustError().Caused("Failed to do something retriable.", struct{ Key, Value interface{} }{Key: ErrorIsRetriable})
}

func ReturnsCriticalError() errors.Error {
	return internalFunctionReturningJustError().Caused("Failed to do something critical", struct{ Key, Value interface{} }{Key: ErrorIsCritical})
}

func internalFunctionReturningJustError() errors.Error {
	return errors.New("internalFunctionReturningJustError always fail")
}
