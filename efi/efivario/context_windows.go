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

//go:build windows

package efivario

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"syscall"

	"golang.org/x/sys/windows"

	"github.com/0x5a17ed/uefi/efi"
	"github.com/0x5a17ed/uefi/efi/binreader"
	"github.com/0x5a17ed/uefi/efi/efiguid"
	"github.com/0x5a17ed/uefi/efi/efivario/efiwindows"
)

type bufferVarEntry struct {
	Length uint32

	Guid efiguid.GUID

	Name []byte
}

func (e *bufferVarEntry) ReadFrom(r io.Reader) (n int64, err error) {
	r = binreader.NewReadTracker(r, &n)

	if _, err = binreader.ReadFields(r, &e.Length, &e.Guid); err != nil {
		return
	}

	e.Name = make([]byte, e.Length-20)
	if _, err = io.ReadFull(r, e.Name); err != nil {
		return
	}

	return
}

type varNameIterator struct {
	buf     *bytes.Buffer
	current *VariableNameItem
	err     error
}

func (it *varNameIterator) Close() error { return nil }

func (it *varNameIterator) Next() bool {
	var entry bufferVarEntry
	if _, err := entry.ReadFrom(it.buf); err != nil {
		if !errors.Is(err, io.EOF) && errors.Is(err, io.ErrUnexpectedEOF) {
			it.err = err
		}
		it.current = nil
		return false
	}

	it.current = &VariableNameItem{
		Name: binreader.UTF16NullBytesToString(entry.Name),
		GUID: entry.Guid,
	}
	return true
}

func (it *varNameIterator) Value() VariableNameItem {
	return *it.current
}

func (it *varNameIterator) Err() error {
	return it.err
}

type WindowsContext struct{}

// Ensure the public facing API in Context is implemented by FsContext.
var _ Context = &WindowsContext{}

func (c WindowsContext) Close() error {
	return nil
}

func (c WindowsContext) VariableNames() (VariableNameIterator, error) {
	var bufLen uint32
	if err := efiwindows.NtEnumerateSystemEnvironmentValuesEx(1, nil, &bufLen); err != nil {
		var ntStatus windows.NTStatus
		if errors.As(err, &ntStatus) && ntStatus != 0xC0000023 {
			return nil, err
		}
	}

	buf := make([]byte, bufLen)
	if err := efiwindows.NtEnumerateSystemEnvironmentValuesEx(1, &buf[0], &bufLen); err != nil {
		return nil, err
	}

	return &varNameIterator{buf: bytes.NewBuffer(buf)}, nil
}

func (c WindowsContext) GetSizeHint(name string, guid efiguid.GUID) (int64, error) {
	lpName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, fmt.Errorf("utf16(%q): %w", name, err)
	}

	var uName windows.NTUnicodeString
	windows.RtlInitUnicodeString(&uName, lpName)

	var bufLen uint32
	err = efiwindows.NtQuerySystemEnvironmentValueEx(&uName, &guid, nil, &bufLen, nil)
	if err != nil && !errors.Is(err, windows.STATUS_BUFFER_TOO_SMALL) {
		return 0, fmt.Errorf("query(%q): %w", name, err)
	}
	return int64(bufLen), nil
}

func (c WindowsContext) GetWithGUID(name string, guid efiguid.GUID, out []byte) (a Attributes, n int, err error) {
	lpName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		err = fmt.Errorf("efivario/utf16(name): %w", err)
		return
	}

	lpGuid, err := syscall.UTF16PtrFromString(guid.Braced())
	if err != nil {
		err = fmt.Errorf("efivario/utf16(guid): %w", err)
		return
	}

	length, err := efiwindows.GetFirmwareEnvironmentVariableEx(lpName, lpGuid, out, (*uint32)(&a))
	if err != nil {
		switch err {
		case windows.ERROR_INSUFFICIENT_BUFFER:
			err = ErrInsufficientSpace
		case windows.ERROR_ENVVAR_NOT_FOUND:
			err = ErrNotFound
		default:
			err = fmt.Errorf("efivario/get: %w", err)
		}
		return
	}
	return a, int(length), err
}

func (c WindowsContext) Get(name string, out []byte) (a Attributes, n int, err error) {
	return c.GetWithGUID(name, efi.GlobalVariable, out)
}

func (c WindowsContext) SetWithGUID(name string, guid efiguid.GUID, attributes Attributes, value []byte) error {
	lpName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return fmt.Errorf("efivario/utf16(name): %w", err)
	}

	lpGuid, err := syscall.UTF16PtrFromString(guid.Braced())
	if err != nil {
		return fmt.Errorf("efivario/utf16(guid): %w", err)
	}

	err = efiwindows.SetFirmwareEnvironmentVariableEx(lpName, lpGuid, value, (uint32)(attributes))
	if err != nil {
		return fmt.Errorf("efivario/set: %w", err)
	}
	return nil
}

func (c WindowsContext) Set(name string, attributes Attributes, value []byte) error {
	return c.SetWithGUID(name, efi.GlobalVariable, attributes, value)
}

func NewDefaultContext() Context {
	return &WindowsContext{}
}
