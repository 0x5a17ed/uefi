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

func ListAllVariables(c efivaraccess.Context) error {
	iter, err := c.VariableNames()
	if err != nil {
		return fmt.Errorf("getIter: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		pp.Println(iter.Value())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("iter/Next: %w", err)
	}
	return nil
}

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

	var listAllVariables bool
	fset.BoolVar(&listAllVariables, "list-all", false, "list all variables")

	var listBootEntries bool
	fset.BoolVar(&listBootEntries, "list-boot", false, "list boot entries")

	var setNextBoot bool
	fset.BoolVar(&setNextBoot, "set-next", false, "set next boot option")

	var nextEntry int
	fset.IntVar(&nextEntry, "next", 0, "boot entry to boot next")

	if err := fset.Parse(args[1:]); err != nil {
		return err
	}

	c := efivaraccess.NewDefaultContext()

	var err error
	switch {
	case listBootEntries:
		err = ReadBootEntries(c)
	case listAllVariables:
		err = ListAllVariables(c)
	case setNextBoot:
		err = efivars.BootNext.Set(c, (uint16)(nextEntry))
	default:
		err = errors.New("no action selected")
	}
	return err
}
