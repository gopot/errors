# subpackage
--
    import "github.com/gopot/errors/examples/pureErrorsExample/subpackage"


## Usage

```go
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
```
Represents Error keys identifiers

#### func  ReturnsCriticalError

```go
func ReturnsCriticalError() errors.Error
```

#### func  ReturnsRetriableError

```go
func ReturnsRetriableError() errors.Error
```
