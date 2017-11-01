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
	"github.com/gopot/errors"
)

func testHelperErrorsAreEqual(e1, e2 errors.Error) bool {
	if (e1 == nil && e2 == nil) ||
		!(e1 != nil && e2 != nil &&
			(e1.Error() != e2.Error() || e1.Detailed() != e2.Detailed())) {
		return true
	}
	return false
}

type mockKVStorage struct {
	getValue func(key interface{}) (value interface{}, found bool)
	stringFn func() string
}

func (this *mockKVStorage) GetValue(key interface{}) (value interface{}, found bool) {
	return this.getValue(key)
}
func (this *mockKVStorage) String() string {
	return this.stringFn()
}

type mockStringer struct {
	s string
}

func (this mockStringer) String() string {
	return this.s
}
