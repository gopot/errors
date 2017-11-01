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

package main

import (
	"time"

	e "github.com/gopot/errors"
	"github.com/gopot/errors/detalizers"
)

// Represents set of immutable Error Key identifiers
var (
	ErrorIsSomething = errorDetailKey("Is something")
)

// package private error factory
var errors = e.NewErrorFactory(e.NewDefaultKVStorage, detalizers.NewCallStackDetalizer(1, 1024), timeDetalizer)

type errorDetailKey string
type errorDetailKeyStringer string

func (this errorDetailKeyStringer) String() string {
	return string(this)
}

func timeDetalizer() []struct{ Key, Value interface{} } {
	return []struct{ Key, Value interface{} }{{errorDetailKeyStringer("Timestamp"), time.Now()}}
}

func main() {

	err := ReturnError()
	if err != nil {
		println(err.Detailed())
		if _, isSomething := err.Get(ErrorIsSomething); isSomething {
			// Handle it here
		}
	}

}

func ReturnError() e.Error {
	return errors.New("Some error", []struct{ Key, Value interface{} }{{Key: ErrorIsSomething}, {Value: "Just in case - it is something!"}}...)
}
