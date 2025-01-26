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
