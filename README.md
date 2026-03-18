# go-testassert

[![CI](https://github.com/philiprehberger/go-testassert/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-testassert/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-testassert.svg)](https://pkg.go.dev/github.com/philiprehberger/go-testassert)
[![License](https://img.shields.io/github/license/philiprehberger/go-testassert)](LICENSE)

Fluent, type-safe test assertions for Go. Built with generics, zero dependencies

## Installation

```bash
go get github.com/philiprehberger/go-testassert
```

## Usage

```go
import "github.com/philiprehberger/go-testassert"
```

### Generic assertions

```go
testassert.That(t, result).Equals(expected)
testassert.That(t, result).NotEquals(other)
testassert.That[*MyStruct](t, nil).IsNil()
testassert.That(t, ptr).IsNotNil()
```

### Ordered comparisons

```go
testassert.ThatOrdered(t, count).IsGreaterThan(0)
testassert.ThatOrdered(t, score).IsLessThan(100)
testassert.ThatOrdered(t, age).IsGreaterOrEqual(18)
testassert.ThatOrdered(t, temp).IsLessOrEqual(37.5)
```

### String assertions

```go
testassert.ThatString(t, msg).Contains("hello")
testassert.ThatString(t, url).HasPrefix("https://")
testassert.ThatString(t, path).HasSuffix(".go")
testassert.ThatString(t, name).HasLen(5)
testassert.ThatString(t, "").IsEmpty()
testassert.ThatString(t, id).Matches(`^[a-f0-9]{8}$`)
```

### Slice assertions

```go
testassert.ThatSlice(t, items).HasLen(3)
testassert.ThatSlice(t, items).Contains("needle")
testassert.ThatSlice(t, items).IsNotEmpty()
testassert.ThatSlice(t, []int{}).IsEmpty()
```

### Error assertions

```go
testassert.ThatError(t, err).IsNil()
testassert.ThatError(t, err).IsNotNil()
testassert.ThatError(t, err).Is(ErrNotFound)
testassert.ThatError(t, err).Contains("timeout")

var pathErr *PathError
testassert.ThatError(t, err).As(&pathErr)
```

### JSON assertions

```go
testassert.ThatJSON(t, body).Equals(`{"name":"alice","age":30}`)
testassert.ThatJSON(t, body).HasKey("name")
testassert.ThatJSON(t, body).Contains("age", 30)
```

### Chaining

All assertions return the receiver, so you can chain calls:

```go
testassert.ThatString(t, msg).
    Contains("hello").
    HasPrefix("hello").
    HasLen(11)

testassert.ThatOrdered(t, n).
    IsGreaterThan(0).
    IsLessThan(100)
```

## API

### Constructors

| Function | Description |
|----------|-------------|
| `That[T](t, got)` | Generic assertion for any type |
| `ThatOrdered[T](t, got)` | Assertion for ordered types (int, float, string) |
| `ThatString(t, got)` | String-specific assertions |
| `ThatSlice[T](t, got)` | Slice-specific assertions |
| `ThatError(t, got)` | Error-specific assertions |
| `ThatJSON(t, got)` | JSON string assertions |

### Assertion Methods

| Type | Methods |
|------|---------|
| `Assertion[T]` | `Equals`, `NotEquals`, `IsNil`, `IsNotNil` |
| `OrderedAssertion[T]` | `Equals`, `NotEquals`, `IsGreaterThan`, `IsLessThan`, `IsGreaterOrEqual`, `IsLessOrEqual` |
| `StringAssertion` | `Equals`, `Contains`, `HasPrefix`, `HasSuffix`, `IsEmpty`, `HasLen`, `Matches` |
| `SliceAssertion[T]` | `HasLen`, `IsEmpty`, `IsNotEmpty`, `Contains` |
| `ErrorAssertion` | `IsNil`, `IsNotNil`, `Is`, `Contains`, `As` |
| `JSONAssertion` | `Equals`, `Contains`, `HasKey` |

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT
