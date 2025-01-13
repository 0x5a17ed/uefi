// Copyright (c) 2025 Arthur Skowronek <0x5a17ed@tuta.io> and contributors
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

package efivario

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"
	"unicode"
	"unsafe"

	"github.com/0x5a17ed/itkit/iters/sliceit"
	"github.com/0x5a17ed/itkit/itlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
	"gotest.tools/v3/golden"

	"github.com/0x5a17ed/uefi/efi/efivario/efiwindows"
)

type sysEnvVarsAPIMock struct {
	buf []byte
}

func (m *sysEnvVarsAPIMock) Get(lpName *uint16, lpGuid *uint16, buf []byte, attrs *uint32) (n uint32, err error) {
	// TODO implement me
	panic("implement me")
}

func (m *sysEnvVarsAPIMock) Set(lpName *uint16, lpGuid *uint16, buf []byte, attrs uint32) (err error) {
	// TODO implement me
	panic("implement me")
}

func (m *sysEnvVarsAPIMock) Enumerate(informationClass uint32, buf *byte, bufLen *uint32) (ntstatus error) {
	if informationClass != 1 {
		return windows.ERROR_INVALID_PARAMETER
	}

	var originalBufLen uint32
	if bufLen != nil {
		originalBufLen = *bufLen
		*bufLen = uint32(len(m.buf))
	}

	if originalBufLen < uint32(len(m.buf)) {
		return windows.STATUS_BUFFER_TOO_SMALL
	}

	if buf != nil {
		copy(unsafe.Slice(buf, originalBufLen), m.buf)
	}

	return
}

func (m *sysEnvVarsAPIMock) Query(
	name *windows.NTUnicodeString,
	guid *efiwindows.GUID,
	buf *byte,
	bufLen *uint32,
	attrs *uint32,
) (ntstatus error) {
	// TODO implement me
	panic("implement me")
}

func stripSpace(s []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func newSysEnvVarsAPIMock(t *testing.T, filename string) *sysEnvVarsAPIMock {
	buf, err := io.ReadAll(hex.NewDecoder(bytes.NewBuffer(stripSpace(golden.Get(t, filename)))))
	require.NoError(t, err)

	return &sysEnvVarsAPIMock{
		buf: buf,
	}
}

func TestWindowsContext_VariableNames(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		efiCtx := &WindowsContext{
			api: newSysEnvVarsAPIMock(t, "windows-efivars-dump"),
		}
		defer efiCtx.Close()

		iter, err := efiCtx.VariableNames()
		require.NoError(t, err)

		s := sliceit.To(itlib.Map(iter.Iter(), func(t VariableNameItem) string { return t.Name }))

		assert.Equal(t, []string{"Alice", "Bob", "Charlie"}, s)
	})
}
