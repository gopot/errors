# detalizers
--
    import "github.com/gopot/errors/detalizers"

### Description

This package contains detalizer(s) plugins for github.com/gopot/errors package

## Usage

```go
const (
	CallStackDetailKey = callStackKey("Call Stack")
)
```

#### func  NewCallStackDetalizer

```go
func NewCallStackDetalizer(skipFrames int, nestLevel int) func() []struct{ Key, Value interface{} }
```
Represents errors.Detalizer implementation to generate call stack detail!
