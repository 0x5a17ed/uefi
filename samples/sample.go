package samples

import (
	"errors"
	"flag"
	"fmt"

	"github.com/k0kubun/pp/v3"

	"github.com/0x5a17ed/uefi/efi/binreader"
	"github.com/0x5a17ed/uefi/efi/efivaraccess"
	"github.com/0x5a17ed/uefi/efi/efivars"
)

func ReadBootEntries(c efivaraccess.Context) error {
	for i := 0; i < 10; i++ {
		fmt.Println(fmt.Sprintf("\nEntry Boot%04d: ", i))

		attrs, lo, err := efivars.Boot(i).Get(c)
		if err != nil {
			if errors.Is(err, efivaraccess.ErrNotFound) {
				fmt.Println("EOF")
				return nil
			}
			return err
		}

		pp.Println(map[string]any{
			"Attributes":   attrs.String(),
			"Description":  binreader.UTF16NullBytesToString(lo.Description),
			"OptionalData": string(lo.OptionalData),
			"Path":         lo.FilePathList.AllText(),
			"LoadOption":   lo,
		})
	}

	return nil
}

func Run(args []string) error {
	fset := flag.NewFlagSet(args[0], flag.ExitOnError)

	var listEntries bool
	fset.BoolVar(&listEntries, "list", false, "list entries")

	var setNextBoot bool
	fset.BoolVar(&setNextBoot, "set-next", false, "set next boot option")

	var nextEntry int
	fset.IntVar(&nextEntry, "next", 0, "boot entry to boot next")

	if err := fset.Parse(args[1:]); err != nil {
		return err
	}

	c := efivaraccess.NewDefaultContext()

	if listEntries {
		if err := ReadBootEntries(c); err != nil {
			return err
		}
	}

	if setNextBoot {
		if err := efivars.BootNext.Set(c, (uint16)(nextEntry)); err != nil {
			return err
		}
	}

	return nil
}
