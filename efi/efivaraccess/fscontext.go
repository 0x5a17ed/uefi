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
	"os"
	"syscall"

	"github.com/spf13/afero"
	"go.uber.org/multierr"

	"github.com/0x5a17ed/uefi/efi"
	"github.com/0x5a17ed/uefi/efi/guid"
)

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

	return f.Sync()
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