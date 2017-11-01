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
	"bufio"
	"log"
	"os"

	"github.com/gopot/errors/examples/pureErrorsExample/subpackage"
)

// Try to build and run it with and without build tag "withcallstack"
// 	go build -tags "withcallstack" .
// 	go run -tags "withcallstack" pureerrorsexample.go
func main() {

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	logger := log.New(f, "pureErrorsExample logger", log.Flags())

	retries := 5

	for i := 0; i < retries; i++ {
		// try to call
		err := subpackage.ReturnsRetriableError()
		if err != nil {
			logger.Println(err.Detailed())
			if _, retirable := err.Get("Is Retriable"); !retirable {

				// Never happens while subpackage.IsRetriable is defined as constant string
				// and subpackage.ReturnsRetriableError returns error containing detil with this key.
				err := err.Caused("Won't to retry. The Error is not retriable.")
				logger.Println(err.Detailed())
				break
			}
		} else {
			break
		}
	}

	// Lets call Critical
	err := subpackage.ReturnsCriticalError()
	if err != nil {
		logger.Println(err.Detailed())
		if _, critical := err.Get(subpackage.ErrorIsCritical); critical {
			// Handle Critical error
		} else {
			// Handle non Critical error
		}
	}
}
