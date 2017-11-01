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
	parent *defaultKVStorage
	Key    interface{}
	Value  interface{}
}

// Implements KVStorage.GetValue(key interface{})(interface{}, bool) method
func (this *defaultKVStorage) GetValue(key interface{}) (value interface{}, found bool) {
	if this != nil {
		if this.Key == key {
			value = this.Value
			found = true
		} else if this.parent != nil {
			value, found = this.parent.GetValue(key)
		}
	}
	return
}

// Implements Stringer interface.
//
// It loops over all KVPairs in the list from last to first and concatenates the following:
// * if KVPair.Key implements fmt.Stringer interface, it prints it with new line at the end
// * if KVPair.Value implements fmt.Stringer interface, it prints it with new line at the end
func (this *defaultKVStorage) String() string {
	ret := ""
	if this != nil {
		var key, value string
		if ks, ok := this.Key.(fmt.Stringer); ok {
			key = ks.String()
		} else if ks, ok := this.Key.(string); ok {
			key = ks
		}
		if vs, ok := this.Value.(fmt.Stringer); ok {
			value = vs.String()
		} else if vs, ok := this.Value.(string); ok {
			value = vs
		}
		if key != "" && value != "" {
			ret += fmt.Sprintln(key, ":", value)
		} else {
			if key != "" {
				ret += fmt.Sprintln(key)
			}
			if value != "" {
				ret += fmt.Sprintln(value)
			}
		}

		if this.parent != nil {
			ret += this.parent.String()
		}
	}
	return ret
}

// Represents default KVStorage factory method.
// The default KVStorage accepts detail keys only of comparable(k1==k2) type(s), otherwise it will panic.
// In case of two or more details have exactly the same Key it returns the last one on GetValue call.
// At the same time, it will print on .String() call all entries.
func NewDefaultKVStorage(pairs ...struct{ Key, Value interface{} }) KVStorage {
	kv := new(defaultKVStorage)
	for _, pair := range pairs {
		kvNew := &defaultKVStorage{
			parent: kv,
			Key:    pair.Key,
			Value:  pair.Value,
		}
		kv = kvNew
	}
	return kv
}
