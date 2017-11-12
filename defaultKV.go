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

// Represents default KV storage implementation as linked list.
type defaultKVStorage struct {
	parent  *defaultKVStorage
	element *struct {
		Key   interface{}
		Value interface{}
	}
}

// Implements KVStorage.GetValue(key interface{})(interface{}, bool) method
func (this *defaultKVStorage) GetValue(key interface{}) (value interface{}, found bool) {
	if this.element != nil && this.element.Key == key {
		return this.element.Value, true
	}

	if this.parent != nil {
		return this.parent.GetValue(key)
	}

	return nil, false
}

// Implements Stringer interface.
//
// It loops over all KVPairs in the list from last to first and concatenates the following:
// * if both Key and Value implement fmt.Stringer or are strings it prints them as 'Key : Value' with new line at the end
// * if only KVPair.Key implements fmt.Stringer interface or is a string, it prints it with new line at the end
// * if only KVPair.Value implements fmt.Stringer interface or is a string, it prints it with new line at the end
func (this *defaultKVStorage) String() string {
	ret := ""
	if this.element != nil {

		var key, value string

		switch ks := (this.element.Key).(type) {
		case fmt.Stringer:
			key = ks.String()
		case string:
			key = ks
		default:
		}

		switch vs := (this.element.Value).(type) {
		case fmt.Stringer:
			value = vs.String()
		case string:
			value = vs
		default:
		}

		if key != "" && value != "" {
			ret += fmt.Sprintln(key, ":", value)
		} else if key+value != "" {
			ret += fmt.Sprintln(key + value)
		}
	}

	if this.parent != nil {
		ret += this.parent.String()
	}

	return ret
}

// Represents default KVStorage factory method.
// The default KVStorage accepts detail keys only of comparable(k1==k2) type(s), otherwise it will panic.
// In case of two or more details have exactly the same Key it returns the last one on GetValue call.
// At the same time, it will print on .String() call all entries.
//
// Implements fmt.Stringer with the following logic:
// It loops over all KVPairs in the list from last to first and concatenates the following:
// * if both Key and Value implement fmt.Stringer or are strings it prints them as 'Key : Value' with new line at the end
// * if only KVPair.Key implements fmt.Stringer interface or is a string, it prints it with new line at the end
// * if only KVPair.Value implements fmt.Stringer interface or is a string, it prints it with new line at the end
func NewDefaultKVStorage(pairs ...struct{ Key, Value interface{} }) KVStorage {
	if len(pairs) == 0 {
		return new(defaultKVStorage)
	}
	var kv *defaultKVStorage
	for _, pair := range pairs {
		kvNew := &defaultKVStorage{
			parent:  kv,
			element: &struct{ Key, Value interface{} }{Key: pair.Key, Value: pair.Value},
		}
		kv = kvNew
	}
	return kv
}
