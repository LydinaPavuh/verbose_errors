package main

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
