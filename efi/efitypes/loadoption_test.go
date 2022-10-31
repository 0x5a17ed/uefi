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

package efitypes_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	assertPkg "github.com/stretchr/testify/assert"
	requirePkg "github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"

	"github.com/0x5a17ed/uefi/efi/efitypes"
)

func readHexdump(r io.Reader) ([]byte, error) {
	var out bytes.Buffer

	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()

		var data string
		for _, c := range l {
			if c != ' ' {
				data += string(c)
			}
		}

		lineData, err := hex.DecodeString(data)
		if err != nil {
			return nil, err
		}
		out.Write(lineData)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func TestLoadOption_ReadFrom(t *testing.T) {
	type args struct {
		fileName string
	}
	type want struct {
		n               int64
		description     string
		filePathStrings []string
		optionalData    []byte
	}
	tt := []struct {
		name    string
		args    args
		want    want
		wantErr assertPkg.ErrorAssertionFunc
	}{
		{"80", args{fileName: "LoadOption80-01.txt"}, want{
			n:               53,
			description:     "TestOption80-01",
			filePathStrings: []string{"Path(128,1,0123456789)"},
			optionalData:    []byte{0x0, 0x0},
		}, assertPkg.NoError},
		{"0404-0401", args{fileName: "LoadOption0404-0401.txt"}, want{
			n:               106,
			description:     "Linux",
			filePathStrings: []string{"HD(1,GPT,FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF,0x800,0x32000)/File(EFI\\LINUX\\GRUB.EFI)"},
			optionalData:    []byte{},
		}, assertPkg.NoError},
		{"05", args{fileName: "LoadOption05.txt"}, want{
			n:               45,
			description:     "TestOption01",
			filePathStrings: []string{`BBS(5,"",0)`},
			optionalData:    []byte{},
		}, assertPkg.NoError},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert := assertPkg.New(t)

			f := golden.Open(t, tc.args.fileName)
			defer f.Close()

			inp, err := readHexdump(f)
			requirePkg.NoError(t, err)

			var lopt efitypes.LoadOption
			gotN, err := lopt.ReadFrom(bytes.NewReader(inp))
			if !tc.wantErr(t, err, fmt.Sprintf("ReadFrom(%v)", tc.args.fileName)) {
				return
			}
			assert.Equal(tc.want.n, gotN)
			assert.Equal(tc.want.description, lopt.DescriptionString())
			assert.Equal(tc.want.filePathStrings, lopt.FilePathList.AllText())
			assert.Equal(tc.want.optionalData, lopt.OptionalData)
		})
	}
}
