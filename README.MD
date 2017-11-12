# errors
--
    import "github.com/gopot/errors"

### Disclaimer

If you are actually hunting for golang library which treats errors as
Exceptions, like in some other languages like Java, .Net, Python, etc. - no good
news for you!


### Description

From authors point of view all errors are of the following representation types:

    1. For User - the error representation to software end-user.
    2. For Code - the error representation options to manage execution flow by code.
    3. For Operations - the error representation in logs/alerts to quickly identify and fix operational issue(s).

This package covers 2nd and 3rd cases. However, it lets developers to simplify
implementations of the 1st.

In the first place - standard errors in GoLang are cool! The only thing they are
missing - is context-like information to take decisions. There are different
approaches to solve it with (+) advantages and (-) disadvantages:

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
interface for all errors all over the code. It provides functionallity to
describe execution flow and description contexts on creation of error. This
could be done by passing `details` to ErrorFactory. Later this details are
accessible via `DetailedError` interface:

    * Detailed() string - returns Operational representation.
    * Get(key interface{}) (interface{}, bool) - returns execution flow representation.

It is also possible to combine both representations for one detail by passing
Key and/or Value which conform(s) fmt.Stringer interface or is of type string.

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


### Build

Usage of build tag "withcallstack" resets Default ErrorFactory whith
NewCallStackDetalizer. It automatically adds call stack information into each
Error produced by default ErrorFactory. See detalizers package for more info.
For.Ex.

    go run -tags "withcallstack" ./examples/pureErrorsExample/pureerrorsexample.go

## Usage

```go
const (
	// Represents Key for value of caused error.
	// Used to extract original error.
	// For more details see Error_causedOriginalError exmple.
	CausedByDetailKey = causedBy("Caused By")
)
```
Represents predefined keys for specific details

#### type DetailedError

```go
type DetailedError interface {

	// Implements standard built-in `error` interface.
	error

	// Returns string representation of detailed error's information.
	Detailed() string

	// Returns value of the detail with given `key` and boolean flag wheather such detail was found.
	Get(key interface{}) (interface{}, bool)
}
```

Represents standard(first class) `error` interface extended with detalization
methods returning context-like information for easier indentification of
specific details.

#### type Detalizer

```go
type Detalizer func() []struct{ Key, Value interface{} }
```

Declares type of functions generating set of details as a slice of `struct{ Key,
Value interface{}}`.

#### type Error

```go
type Error interface {
	DetailedError

	// Returns a new `Error` as caused by current.
	Caused(string, ...struct{ Key, Value interface{} }) Error
}
```

Represents `DetailedError` interface extended with `Caused` method.

#### func  ConvertToError

```go
func ConvertToError(e error) Error
```
Gently converts error `e` into Error trying to preserve original types and data

#### func  New

```go
func New(text string) Error
```
New returns an Error that formats as the given text.

#### func  NewErrorf

```go
func NewErrorf(format string, a ...interface{}) Error
```
Errorf formats according to a format specifier and returns the string as a value
that satisfies Error

#### func  NewWithDetails

```go
func NewWithDetails(text string, details ...struct{ Key, Value interface{} }) Error
```
New returns an Error that formats as the given text with provided details.

#### type ErrorFactory

```go
type ErrorFactory interface {

	// Creates a new Error with `message` as `Error()` content and details.
	New(message string, details ...struct{ Key, Value interface{} }) Error

	// Converts basic built-in `error` into `Error`.
	ConvertToError(e error) Error
}
```

Represents Error factory functionality.

#### func  NewErrorFactory

```go
func NewErrorFactory(kvFactory KVStorageFactory, detalizers ...Detalizer) ErrorFactory
```
Creates a default ErrorFactory which executes provided `detalizers` uppon each
`ErrorFactory.New(...)` call. Then, detalizers results are appended to those
which are passed to ErrorFactory.New(...) and fed to kvFactory to create
back-end details storage.

On ConvertToError(e) call default ErrorFactory works as following:

    * in case error `e` implements `Error` - it returns `e` asserted to `Error`
    * if not, it creates new `Error` with message as `e.Error()`

Errors produced by default ErrorFactory are utilizing KVStorage interface:

    * to Get -> kvStorage.GetValue
    * to print Detailed -> kvStorage.String()

NewErrorFactory(...) panics in case `kvFactory` function is nil.

#### type KVStorage

```go
type KVStorage interface {

	// Returns value by given key.
	// The boolean flag indicates whether the key was found(true) or not(false).
	GetValue(key interface{}) (value interface{}, found bool)

	// Returns string representation of content.
	String() string
}
```

Represents simple immutable Key-Value storage functionality.

#### func  NewDefaultKVStorage

```go
func NewDefaultKVStorage(pairs ...struct{ Key, Value interface{} }) KVStorage
```
Represents default KVStorage factory method. The default KVStorage accepts
detail keys only of comparable(k1==k2) type(s), otherwise it will panic. In case
of two or more details have exactly the same Key it returns the last one on
GetValue call. At the same time, it will print on .String() call all entries.

#### type KVStorageFactory

```go
type KVStorageFactory func(pairs ...struct{ Key, Value interface{} }) KVStorage
```

Represents KVStorage factory.
