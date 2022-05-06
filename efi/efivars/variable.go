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

	"github.com/0x5a17ed/uefi/efi/efiguid"
	"github.com/0x5a17ed/uefi/efi/efitypes"
	"github.com/0x5a17ed/uefi/efi/efivario"
)

const (
	BootNextName    = "BootNext"
	BootCurrentName = "BootCurrent"
)

type Variable[T any] struct {
	name         string
	guid         efiguid.GUID
	defaultAttrs efivario.Attributes
}

func (e Variable[T]) Get(c efivario.Context) (attrs efivario.Attributes, value T, err error) {
	attrs, data, err := efivario.ReadAll(c, e.name, e.guid)
	if err != nil {
		err = fmt.Errorf("efivars/get(%s): load: %w", e.name, err)
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
		err = fmt.Errorf("efivars/get(%s): parse: %w", e.name, err)
		return
	}

	return
}

func (e Variable[T]) SetWithAttributes(c efivario.Context, attrs efivario.Attributes, value T) (err error) {
	var buf bytes.Buffer

	var valueInterface any = &value
	if writer, ok := (valueInterface).(io.WriterTo); ok {
		_, err = writer.WriteTo(&buf)
	} else {
		err = binary.Write(&buf, binary.LittleEndian, value)
	}
	if err != nil {
		return fmt.Errorf("efivars/set(%s): %w", e.name, err)
	}

	return c.Set(e.name, e.guid, attrs, buf.Bytes())
}

func (e Variable[T]) Set(c efivario.Context, value T) error {
	return e.SetWithAttributes(c, e.defaultAttrs, value)
}

var (
	// BootNext specifies the first boot option on the next boot.
	//
	// <https://uefi.org/sites/default/files/resources/UEFI_Spec_2_9_2021_03_18.pdf#G14.1012867>
	BootNext = Variable[uint16]{
		name:         BootNextName,
		guid:         GlobalVariable,
		defaultAttrs: efivario.NonVolatile | efivario.BootServiceAccess | efivario.RuntimeAccess,
	}

	// BootCurrent defines the Boot#### option that was selected
	// on the current boot.
	BootCurrent = Variable[uint16]{
		name:         BootCurrentName,
		guid:         GlobalVariable,
		defaultAttrs: efivario.NonVolatile | efivario.BootServiceAccess | efivario.RuntimeAccess,
	}
)

// Boot returns an EFI Variable pointing to the boot LoadOption
// for the given index.
//
// <https://uefi.org/sites/default/files/resources/UEFI_Spec_2_9_2021_03_18.pdf#G7.1346720>
func Boot(i int) Variable[efitypes.LoadOption] {
	return Variable[efitypes.LoadOption]{
		name: fmt.Sprintf("Boot%04X", i),
		guid: GlobalVariable,
	}
}
