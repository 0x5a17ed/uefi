package main

import (
	"fmt"
	"os"

	"github.com/0x5a17ed/uefi/samples"
)

func main() {
	if err := samples.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
