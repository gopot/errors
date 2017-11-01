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

package detalizers

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

const (
	CallStackDetailKey = callStackKey("Call Stack")
)

// Represents type alias for Call Stack key
type callStackKey string

// Implements fmt.Stringer() which is consumed by DetailedError/Error interfaces
// to generate Detailed() result.
func (this callStackKey) String() string {
	return string(this)
}

// Represents internal implementation of Call Stack value. Should implement fmt.Stringer
type callStackValue struct {
	frames  *runtime.Frames
	once    sync.Once
	printed string
}

// Represents errors.Detalizer implementation to generate call stack detail!
func NewCallStackDetalizer(skipFrames int, nestLevel int) func() []struct{ Key, Value interface{} } {
	return func() []struct{ Key, Value interface{} } {

		pc := make([]uintptr, nestLevel)
		n := runtime.Callers(skipFrames+2, pc)
		if n == 0 {
			// No pcs available. Stop now.
			// This can happen if the first argument to runtime.Callers is large.
			return nil
		}
		pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
		frames := runtime.CallersFrames(pc)

		cs := &callStackValue{frames: frames}

		return []struct{ Key, Value interface{} }{struct{ Key, Value interface{} }{CallStackDetailKey, cs}}
	}
}

// fmt.Stringer implementation.
// Prints callstack with each entry as "[package path]/[package name] [package file]:[line] [functionname]"
func (this *callStackValue) String() string {

	// Frames could be safely read only once
	this.once.Do(func() {
		this.printed = fmt.Sprintln("")
		for {
			frame, more := this.frames.Next()

			functionNameSplitBySlash := strings.Split(frame.Function, "/")
			functionNameParts := strings.Split(functionNameSplitBySlash[len(functionNameSplitBySlash)-1], ".")
			functionName := strings.Join(functionNameParts[1:], ".")

			pkgPath := ""
			if len(functionNameSplitBySlash) > 1 {
				pkgPath = strings.Join(functionNameSplitBySlash[:len(functionNameSplitBySlash)-1], "/") + "/"
			}

			pkgName := functionNameParts[0]

			fileNameParts := strings.Split(frame.File, "/")
			fileName := fileNameParts[len(fileNameParts)-1]

			this.printed += fmt.Sprintln(fmt.Sprintf("\t%s%s.%s %s:%d", pkgPath, pkgName, functionName, fileName, frame.Line))

			if !more {
				break
			}

		}
	})

	return this.printed
}
