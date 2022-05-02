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
	"encoding/hex"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/0x5a17ed/uefi/efi/efivaraccess"
)

func TestBoot(t *testing.T) {
	// objective of the test is to ensure that Boot(n) correctly
	// converts n to hexadecimal.
	fs := afero.NewMemMapFs()

	const s = "000000000000000000000000"

	decoded, err := hex.DecodeString(s)
	assert.NoError(t, err)

	f, _ := fs.Create("Boot000F-8BE4DF61-93CA-11D2-AA0D-00E098032B8C")
	f.Write(decoded)
	f.Close()

	c := efivaraccess.NewFileSystemContext(fs)

	_, _, err = Boot(15).Get(c)
	assert.NoError(t, err)
}
