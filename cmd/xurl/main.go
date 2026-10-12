package main

import (
	"fmt"
	"os"
)

func main() {
	if err := NewRootCommand().Execute(); err != nil {
		if xe, ok := err.(interface{ Error() string }); ok {
			fmt.Fprintln(os.Stderr, xe.Error())
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}
}
