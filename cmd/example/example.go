package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/k0kubun/pp/v3"

	"github.com/0x5a17ed/uefi/efi/efireader"
	"github.com/0x5a17ed/uefi/efi/efivario"
	"github.com/0x5a17ed/uefi/efi/efivars"
)

var (
	bootRe = regexp.MustCompile(`^Boot([\da-fA-F]{4})$`)
)

func ListBootOrder(c efivario.Context) error {
	_, value, err := efivars.BootOrder.Get(c)
	if err != nil {
		return err
	}

	for i, index := range value {
		_, lo, err := efivars.Boot(index).Get(c)
		if err != nil {
			return fmt.Errorf("entry %d (%04[1]X): %w", index, err)
		}

		pp.Println(map[string]any{
			"Order":       i,
			"Index":       index,
			"Description": efireader.UTF16NullBytesToString(lo.Description),
			"Path":        lo.FilePathList.AllText(),
		})
	}

	return nil
}

func ListAllVariables(c efivario.Context) error {
	iter, err := c.VariableNames()
	if err != nil {
		return fmt.Errorf("getIter: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		v := iter.Value()

		size, err := c.GetSizeHint(v.Name, v.GUID)

		var errString string
		if err != nil {
			errString = err.Error()
		}

		pp.Println(map[string]any{
			"Name": v.Name,
			"GUID": v.GUID.Braced(),
			"Size": []any{size, errString},
		})
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("iter/Next: %w", err)
	}
	return nil
}

func ReadBootEntries(c efivario.Context) error {
	iter, err := c.VariableNames()
	if err != nil {
		return fmt.Errorf("getIter: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		v := iter.Value()

		if v.GUID != efivars.GlobalVariable {
			continue
		}

		matches := bootRe.FindStringSubmatch(v.Name)
		if matches == nil {
			continue
		}

		value, err := strconv.ParseInt(matches[1], 16, 16)
		if err != nil {
			continue
		}

		fmt.Printf("\nEntry Boot%04X(%[1]d):\n", value)

		attrs, lo, err := efivars.Boot(uint16(int(value))).Get(c)
		if err != nil {
			return err
		}

		pp.Println(map[string]any{
			"Attributes":   attrs.String(),
			"Description":  efireader.UTF16NullBytesToString(lo.Description),
			"OptionalData": string(lo.OptionalData),
			"Path":         lo.FilePathList.AllText(),
			"LoadOption":   lo,
		})
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("iter/Next: %w", err)
	}

	return nil
}

func Run(args []string) error {
	fset := flag.NewFlagSet(path.Base(args[0]), flag.ExitOnError)

	var listAllVariables bool
	fset.BoolVar(&listAllVariables, "list-all", false, "list all variables")

	var listBootEntries bool
	fset.BoolVar(&listBootEntries, "list-boot", false, "list boot entries")

	var listBootOrder bool
	fset.BoolVar(&listBootOrder, "list-boot-order", false, "list boot order")

	var setNextBoot bool
	fset.BoolVar(&setNextBoot, "set-next", false, "set next boot option")

	var nextEntry int
	fset.IntVar(&nextEntry, "next", 0, "boot entry to boot next")

	if err := fset.Parse(args[1:]); err != nil {
		return err
	}

	c := efivario.NewDefaultContext()

	var err error
	switch {
	case listBootOrder:
		err = ListBootOrder(c)
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

func main() {
	err := RunWithPrivileges(func() error {
		return Run(os.Args)
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
