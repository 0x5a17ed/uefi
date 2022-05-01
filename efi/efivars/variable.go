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

package efivars

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/0x5a17ed/uefi/efi"
	"github.com/0x5a17ed/uefi/efi/efitypes"
	"github.com/0x5a17ed/uefi/efi/efivaraccess"
	"github.com/0x5a17ed/uefi/efi/efivarioutil"
	"github.com/0x5a17ed/uefi/efi/guid"
)

const (
	BootNextName = "BootNext"
)

type Variable[T any] struct {
	name         string
	guid         guid.GUID
	defaultAttrs efi.Attributes
}

func (e Variable[T]) Get(c *efivaraccess.Context) (attrs efi.Attributes, value T, err error) {
	attrs, data, err := efivarioutil.ReadAllWitGuid(c, e.name, e.guid)
	if err != nil {
		err = fmt.Errorf("efivars: failed reading value: %w", err)
		return
	}

	buf := bytes.NewReader(data)

	var valueInterface any = &value
	if reader, ok := (valueInterface).(io.ReaderFrom); ok {
		_, err = reader.ReadFrom(buf)
	} else {
		err = binary.Read(buf, binary.LittleEndian, &value)
	}
	if err != nil {
		err = fmt.Errorf("efivars/get: failed decoding value: %w", err)
		return
	}

	return
}

func (e Variable[T]) SetWithAttributes(c *efivaraccess.Context, attrs efi.Attributes, value T) (err error) {
	var buf bytes.Buffer

	var valueInterface any = &value
	if writer, ok := (valueInterface).(io.WriterTo); ok {
		_, err = writer.WriteTo(&buf)
	} else {
		err = binary.Write(&buf, binary.LittleEndian, value)
	}
	if err != nil {
		return fmt.Errorf("efivars/set: failed encoding value: %w", err)
	}

	return c.SetWithGUID(e.name, e.guid, attrs, buf.Bytes())
}

func (e Variable[T]) Set(c *efivaraccess.Context, value T) error {
	return e.SetWithAttributes(c, e.defaultAttrs, value)
}

var (
	// BootNext specifies the first boot option on the next boot.
	//
	// <https://uefi.org/sites/default/files/resources/UEFI_Spec_2_9_2021_03_18.pdf#G14.1012867>
	BootNext = Variable[uint16]{
		name:         BootNextName,
		guid:         efi.GlobalVariable,
		defaultAttrs: efi.NonVolatile | efi.BootServiceAccess | efi.RuntimeAccess,
	}
)

func Boot(i int) Variable[efitypes.LoadOption] {
	return Variable[efitypes.LoadOption]{
		name: fmt.Sprintf("Boot%04d", i),
		guid: efi.GlobalVariable,
	}
}
