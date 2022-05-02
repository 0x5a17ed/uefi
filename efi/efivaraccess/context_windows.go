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

package efivaraccess

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"

	"github.com/0x5a17ed/uefi/efi"
	"github.com/0x5a17ed/uefi/efi/efivaraccess/efiwindows"
	"github.com/0x5a17ed/uefi/efi/guid"
)

type WindowsContext struct{}

// Ensure the public facing API in Context is implemented by FsContext.
var _ Context = &WindowsContext{}

func (c WindowsContext) Close() error {
	return nil
}

func (c WindowsContext) GetWithGUID(name string, guid guid.GUID, out []byte) (a efi.Attributes, n int, err error) {
	lpName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		err = fmt.Errorf("efivaraccess/utf16(name): %w", err)
		return
	}

	lpGuid, err := syscall.UTF16PtrFromString(guid.Braced())
	if err != nil {
		err = fmt.Errorf("efivaraccess/utf16(guid): %w", err)
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
			err = fmt.Errorf("efivaraccess/get: %w", err)
		}
		return
	}
	return a, int(length), err
}

func (c WindowsContext) Get(name string, out []byte) (a efi.Attributes, n int, err error) {
	return c.GetWithGUID(name, efi.GlobalVariable, out)
}

func (c WindowsContext) SetWithGUID(name string, guid guid.GUID, attributes efi.Attributes, value []byte) error {
	lpName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return fmt.Errorf("efivaraccess/utf16(name): %w", err)
	}

	lpGuid, err := syscall.UTF16PtrFromString(guid.Braced())
	if err != nil {
		return fmt.Errorf("efivaraccess/utf16(guid): %w", err)
	}

	err = efiwindows.SetFirmwareEnvironmentVariableEx(lpName, lpGuid, value, (uint32)(attributes))
	if err != nil {
		return fmt.Errorf("efivaraccess/set: %w", err)
	}
	return nil
}

func (c WindowsContext) Set(name string, attributes efi.Attributes, value []byte) error {
	return c.SetWithGUID(name, efi.GlobalVariable, attributes, value)
}

func NewDefaultContext() Context {
	return &WindowsContext{}
}
