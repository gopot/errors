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

// Represents simple immutable Key-Value storage functionality.
type KVStorage interface {

	// Returns value by given key.
	// The boolean flag indicates whether the key was found(true) or not(false).
	GetValue(key interface{}) (value interface{}, found bool)

	// Returns string representation of content.
	String() string
}

// Represents KVStorage factory.
type KVStorageFactory func(pairs ...struct{ Key, Value interface{} }) KVStorage
