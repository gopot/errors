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

/*

Disclaimer

If you are actually hunting for golang library which treats errors as Exceptions, like in some other languages like Java, .Net, Python, etc. - no good news for you!

Description

From authors point of view all errors are of the following representation types:

    1. For User - the error representation to software end-user.
    2. For Code - the error representation options to manage execution flow by code.
    3. For Operations - the error representation in logs/alerts to quickly identify and fix operational issue(s).

This package covers 2nd and 3rd cases. However, it lets developers to simplify implementations of the 1st.

In the first place - standard errors in GoLang are cool! The only thing they are missing - is context-like information to take decisions.
There are different approaches to solve it with (+) advantages and (-) disadvantages:

    0. Define all possible in the package errors as public variables and let consumers compare returned errors to them.
        For.Ex. https://golang.org/pkg/net/http/#pkg-variables ('Errors used by the HTTP server.' block)
        + very easy implementation
        + extremely performant
        - consumer CAN redefine it
        - such errors do not contain explicit information about particular case(s)
    1. Define own Error type as structure with detail information fields.
        For.Ex. https://golang.org/pkg/os/#PathError
        + strongly typed
        + as much information as fields provided
        - consumer must know such Error structure
        - data contract, rather than behaviour
        - direct fields access = cosumer CAN change it
    2. Define own Error interface (in combination with p.2) with contextual properties.
        For.Ex. https://golang.org/pkg/net/#Error
        + strongly typed
        + behaviour contract, rather than data
        Actually, this approach does not have irresolvable disadvantages, however:
        - does not(at least easily) provide good Operational representation.
        - not flexible enough to add new property on the fly.

This package contains implementation of another approach which has its own
advantages and, obviously, disadvantages. This solution declares single
interface for all errors all over the code. It provides functionallity to describe
execution flow and description contexts on creation of error. This could be done by
passing `details` to ErrorFactory. Later this details are accessible via `DetailedError`
interface:

    * Detailed() string - returns Operational representation.
    * Get(key interface{}) (interface{}, bool) - returns execution flow representation.

It is also possible to combine both representations for one detail by passing Key and/or Value which conform(s)
fmt.Stringer interface or is of type string.

Thus, consumer of the following error

    // Expose immutable error property
    const ErrorIsCritical = myPackageErrorDetailStringer("Is Critical")

    type myPackageErrorDetailStringer string
    func (this myPackageErrorDetailStringer) String() string { return this }

    func DoSomething() Error {

        ...

        return errorFactory.New(
            // Some error message
            "Some error",
            // The detail which extends dry GoLang error with context
            struct{ Key, Value interface{}}{ Key: ErrorIsCritical }
        )
    }

Is able to utilize it in any representation.

    err := myPackage.DoSomething()
    if err != nil {
        // In this place Operations representation is utilized. Due to ErrorIsCritical
        // conforms fmt.Stringer() err.Detailed() will log it.
        log.Info(err.Detailed())

        // Here Code representation is utilized.
        // We define execution flow based on existance of myPackage.ErrorIsCritical flag in Error details.
        if _, isCritical := Error.Get(myPackage.ErrorIsCritical); isCritical{
            // Now we can manage CRITICAL error
        } esle {
            // Error is not critical
        }
    }

An complete usage is shown in [examples folder](./examples).

Some detalizers are supplied at [detalizers subpackage](./detalizers).

Build

Usage of build tag "withcallstack" resets Default ErrorFactory whith NewCallStackDetalizer. It automatically
adds call stack information into each Error produced by default ErrorFactory. See detalizers package for more info.
For.Ex.

    go run -tags "withcallstack" ./examples/pureErrorsExample/pureerrorsexample.go

*/
package errors
