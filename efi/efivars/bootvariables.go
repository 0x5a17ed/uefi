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
	"fmt"

	"github.com/0x5a17ed/uefi/efi/efitypes"
	"github.com/0x5a17ed/uefi/efi/efivario"
)

const (
	BootNextName    = "BootNext"
	BootCurrentName = "BootCurrent"
)

var (
	// BootNext specifies the first boot option on the next boot.
	//
	// <https://uefi.org/sites/default/files/resources/UEFI_Spec_2_9_2021_03_18.pdf#G7.1346720>
	BootNext = Variable[uint16]{
		name:         BootNextName,
		guid:         GlobalVariable,
		defaultAttrs: efivario.NonVolatile | efivario.BootServiceAccess | efivario.RuntimeAccess,
	}

	// BootCurrent defines the Boot#### option that was selected
	// on the current boot.
	//
	// <https://uefi.org/sites/default/files/resources/UEFI_Spec_2_9_2021_03_18.pdf#G7.1346720>
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
