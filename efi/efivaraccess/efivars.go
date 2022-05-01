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
	"errors"
	"io"

	"github.com/0x5a17ed/uefi/efi"
	"github.com/0x5a17ed/uefi/efi/guid"
)

var (
	ErrInsufficientSpace = errors.New("buffer too small")
	ErrNotFound          = errors.New("variable not found")
)

type contextInterface interface {
	io.Closer

	// GetWithGUID reads a specific EFI variable and returns the content in slice indicated by out.
	GetWithGUID(name string, guid guid.GUID, out []byte) (efi.Attributes, int, error)

	// Get reads a global EFI variable and returns the content in slice indicated by out.
	Get(name string, out []byte) (efi.Attributes, int, error)

	// SetWithGUID writes a specific EFI variable.
	SetWithGUID(name string, guid guid.GUID, attributes efi.Attributes, value []byte) error

	// Set writes a global EFI variable.
	Set(name string, attributes efi.Attributes, value []byte) error
}

// Ensure the public facing API is implemented by Context.
var _ contextInterface = &Context{}
