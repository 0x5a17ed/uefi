// Copyright (c) 2022 Arthur Skowronek <0x5a17ed@tuta.io> and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// <https://www.apache.org/licenses/LICENSE-2.0>
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package efivaraccess

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"syscall"

	"github.com/spf13/afero"
	"go.uber.org/multierr"

	"github.com/0x5a17ed/uefi/efi"
	"github.com/0x5a17ed/uefi/efi/guid"
)

var nameRegex = regexp.MustCompile(`^([^-]+)-([\da-fA-F]{8}-[\da-fA-F]{4}-[\da-fA-F]{4}-[\da-fA-F]{4}-[\da-fA-F]{12})$`)

// iterator is a variable name iterator for FsContext
type iterator struct {
	f       afero.File
	err     error
	current *VariableNameItem
}

func (it *iterator) Close() error {
	return it.f.Close()
}

func (it *iterator) Next() bool {
	for {
		names, err := it.f.Readdirnames(1)
		if err != nil || len(names) < 1 {
			if !errors.Is(err, io.EOF) {
				it.err = err
			}
			it.current = nil
			return false
		}

		matches := nameRegex.FindStringSubmatch(names[0])
		if matches == nil {
			continue
		}

		g, err := guid.FromString(matches[2])
		if err != nil {
			it.err = err
			it.current = nil
			return false
		}

		it.current = &VariableNameItem{Name: matches[1], GUID: g}
		return true
	}
}

func (it *iterator) Value() VariableNameItem {
	return *it.current
}

func (it *iterator) Err() error {
	return it.err
}

func getFileName(name string, guid guid.GUID) string {
	return fmt.Sprintf("%s-%s", name, guid)
}

type FsContext struct {
	fs afero.Fs
}

// Ensure the public facing API in Context is implemented by FsContext.
var _ Context = &FsContext{}

func (c FsContext) Close() error {
	return nil
}

func (c FsContext) VariableNames() (VariableNameIterator, error) {
	f, err := c.fs.Open("")
	if err != nil {
		return nil, err
	}
	return &iterator{f: f}, nil
}

func (c FsContext) readEfiVarFileName(name string, out []byte) (a efi.Attributes, n int, err error) {
	f, err := c.fs.Open(name)
	if err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && pathErr.Err == syscall.ENOENT {
			err = ErrNotFound
		} else {
			err = fmt.Errorf("efivaraccess/get: %w", err)
		}
		return
	}
	defer multierr.AppendInvoke(&err, multierr.Close(f))

	if err = binary.Read(f, binary.LittleEndian, &a); err != nil {
		err = fmt.Errorf("efivaraccess/get: %w", err)
		return
	}

	n, rerr := io.ReadFull(f, out)
	switch rerr {
	case nil:
		err = ErrInsufficientSpace
	case io.ErrUnexpectedEOF:
		err = nil
	default:
	}
	return
}

func (c FsContext) writeEfiVarFileName(name string, value []byte, attrs efi.Attributes) error {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, attrs); err != nil {
		return fmt.Errorf("efivaraccess/set: write attr: %w", err)
	}

	if _, err := buf.Write(value); err != nil {
		return fmt.Errorf("efivaraccess/set: write value: %w", err)
	}

	flags := os.O_WRONLY | os.O_CREATE
	if attrs&efi.AppendWrite != 0 {
		flags |= os.O_APPEND
	}

	f, err := c.fs.OpenFile(name, flags, 0644)
	if err != nil {
		return fmt.Errorf("efivaraccess/set: %w", err)
	}
	defer multierr.AppendInvoke(&err, multierr.Close(f))

	if _, err := buf.WriteTo(f); err != nil {
		return fmt.Errorf("efivaraccess/set: %w", err)
	}

	if err := f.Sync(); err != nil {
		var errno syscall.Errno
		switch {
		case errors.Is(err, fs.ErrInvalid):
			fallthrough
		case errors.As(err, &errno) && errno == syscall.EINVAL:
			// fsync is not implemented by efivarfs yet so
			// calling it here might sound silly, which it
			// actually is.  Lets just ignore it for now.
			return nil
		default:
			return err
		}
	}
	return nil
}

func (c FsContext) GetSizeHint(name string, guid guid.GUID) (int64, error) {
	fi, err := c.fs.Stat(getFileName(name, guid))
	if err != nil {
		return 0, err
	}
	return fi.Size() - 4, nil
}

func (c FsContext) GetWithGUID(name string, guid guid.GUID, out []byte) (a efi.Attributes, n int, err error) {
	return c.readEfiVarFileName(getFileName(name, guid), out)
}

func (c FsContext) Get(name string, out []byte) (a efi.Attributes, n int, err error) {
	return c.GetWithGUID(name, efi.GlobalVariable, out)
}

func (c FsContext) SetWithGUID(name string, guid guid.GUID, attributes efi.Attributes, value []byte) error {
	return c.writeEfiVarFileName(getFileName(name, guid), value, attributes)
}

func (c FsContext) Set(name string, attributes efi.Attributes, value []byte) error {
	return c.SetWithGUID(name, efi.GlobalVariable, attributes, value)
}

func NewFileSystemContext(fs afero.Fs) *FsContext {
	return &FsContext{fs: fs}
}
