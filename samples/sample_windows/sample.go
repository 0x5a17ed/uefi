//go:build windows

package main

import (
	"os"

	"github.com/Microsoft/go-winio"
	"github.com/k0kubun/pp/v3"

	"github.com/0x5a17ed/uefi/samples"
)

func main() {
	err := winio.RunWithPrivilege("SeSystemEnvironmentPrivilege", func() error {
		return samples.Run(os.Args)
	})

	if err != nil {
		pp.Println(err)
	}
}
