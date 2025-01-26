## Tools for tracing error and readable printing of wrapped errors
[![Go Report Card](https://goreportcard.com/badge/github.com/LydinaPavuh/verbose_errors)](https://goreportcard.com/report/github.com/LydinaPavuh/verbose_errors)


## Wrapped errors printing
```go
import (
	"fmt"
	"github.com/LydinaPavuh/verbose_errors/error_three"
)

func main() {
	third := error_three.Error("Third", error_three.Error("s3"), error_three.Error("s3.3"))
	second := error_three.Error("Second", error_three.Error("s2"), error_three.Error("s2.1", third))
	first := error_three.Error("Root err", second)

	fmt.Println(error_three.PrintErrorThree(first))
}
```

## Errors with trace
```go
package main

import (
	"fmt"
	"github.com/LydinaPavuh/verbose_errors/tracer"
)

func A() error {
	return tracer.New("Some error")
}

func B() error {
	return A()
}

func main() {
	err := B()
	fmt.Print(tracer.PrintTrace(err))
}

```