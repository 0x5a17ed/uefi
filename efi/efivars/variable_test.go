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
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/multierr"

	"github.com/0x5a17ed/uefi/efi/efiguid"
	"github.com/0x5a17ed/uefi/efi/efivario"
)

var testGuid = efiguid.MustFromString("3cd99f3f-4b2b-43eb-ac29-f0890a4772b7")

func readFile(fs afero.Fs, name string) (string, error) {
	f, err := fs.Open(name)
	if err != nil {
		return "", fmt.Errorf("readFile/open: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(hex.NewEncoder(&buf), f); err != nil {
		return "", fmt.Errorf("readFile/read: %w", err)
	}

	return buf.String(), nil
}

func writeFile(fs afero.Fs, name, data string) (err error) {
	f, _ := fs.Create(name)
	defer multierr.AppendInvoke(&err, multierr.Close(f))
	_, err = io.Copy(f, hex.NewDecoder(strings.NewReader(data)))
	return
}

func wantedError(target error) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
		return assert.ErrorIs(t, err, target, msgAndArgs...)
	}
}

type testRow[T any] struct {
	variable Variable[T]
	name     string
	data     string
	wanted   T
	wantErr  assert.ErrorAssertionFunc
}

func (r *testRow[T]) fileName() string {
	return fmt.Sprintf("%s-%s", r.variable.name, r.variable.guid)
}

type testEnv[T any] struct {
	t   *testing.T
	fs  afero.Fs
	ctx *efivario.FsContext
}

func (te *testEnv[T]) testGet(row *testRow[T]) (ok bool) {
	gotAttrs, gotValue, err := row.variable.Get(te.ctx)
	if !row.wantErr(te.t, err, "efivars/get") {
		return
	}

	if !assert.Equal(te.t, row.variable.defaultAttrs, gotAttrs) {
		return
	}
	if !assert.Equal(te.t, row.wanted, gotValue) {
		return
	}
	return true
}

func (te testEnv[T]) testSet(row *testRow[T]) (ok bool) {
	err := row.variable.Set(te.ctx, row.wanted)
	if !row.wantErr(te.t, err, "efivars/set") {
		return
	}

	content, err := readFile(te.fs, row.fileName())
	if !assert.NoError(te.t, err) {
		return
	}

	if !assert.Equal(te.t, row.data, content[8:]) {
		return
	}

	return true
}

func newTestEnv[T any](t *testing.T) *testEnv[T] {
	fs := afero.NewMemMapFs()
	return &testEnv[T]{
		t:   t,
		fs:  fs,
		ctx: efivario.NewFileSystemContext(fs),
	}
}

func (te testEnv[T]) setupVarFile(row *testRow[T]) error {
	return writeFile(te.fs, row.fileName(), "07000000"+row.data)
}

type testRunner[T any] func(*testing.T, *testRow[T], *testEnv[T])

func runTests[T any](t *testing.T, tests []*testRow[T], fn testRunner[T]) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newTestEnv[T](t)
			fn(t, tt, e)
		})
	}
}

type VariableTestSuite struct{ suite.Suite }

func (s *VariableTestSuite) TestPrimitive() {
	s.Run("uint16", func() {
		v := Variable[uint16]{
			name:         "TestVar",
			guid:         testGuid,
			defaultAttrs: defaultAttrs,
			unmarshal:    primitiveUnmarshaller[uint16],
			marshal:      primitiveMarshaller[uint16],
		}

		row := &testRow[uint16]{
			v, "", "1234", uint16(0x3412), assert.NoError,
		}

		tfn := newTestEnv[uint16](s.T())
		if !tfn.testSet(row) {
			return
		}
		tfn.testGet(row)
	})

	s.Run("SetGet/bool", func() {
		v := Variable[bool]{
			name:         "TestVar",
			guid:         testGuid,
			defaultAttrs: defaultAttrs,
			unmarshal:    primitiveUnmarshaller[bool],
			marshal:      primitiveMarshaller[bool],
		}

		row := &testRow[bool]{
			v, "", "01", true, assert.NoError,
		}

		tfn := newTestEnv[bool](s.T())
		if !tfn.testSet(row) {
			return
		}
		tfn.testGet(row)
	})
}

func (s *VariableTestSuite) TestSlice() {
	v := Variable[[]uint16]{
		name:         "TestVar",
		guid:         testGuid,
		defaultAttrs: defaultAttrs,
		marshal:      sliceMarshaller[uint16],
		unmarshal:    sliceUnmarshaller[uint16],
	}

	s.Run("SetGet", func() {
		var tests = []*testRow[[]uint16]{
			{v, "Empty", "", []uint16(nil), assert.NoError},
			{v, "One", "1234", []uint16{0x3412}, assert.NoError},
			{v, "Many", "12345678", []uint16{0x3412, 0x7856}, assert.NoError},
		}

		runTests(s.T(), tests, func(t *testing.T, row *testRow[[]uint16], env *testEnv[[]uint16]) {
			if !env.testSet(row) {
				return
			}
			env.testGet(row)
		})
	})

	s.Run("Get/ErrShort", func() {
		rows := []*testRow[[]uint16]{
			{v, "", "1234", []uint16{0x3412}, assert.NoError},
			{v, "", "123456", []uint16{0x3412}, wantedError(io.ErrUnexpectedEOF)},
		}

		runTests(s.T(), rows, func(t *testing.T, row *testRow[[]uint16], env *testEnv[[]uint16]) {
			if !assert.NoError(s.T(), env.setupVarFile(row)) {
				return
			}
			env.testGet(row)
		})
	})
}

type data struct {
	A uint16
	B uint32
}

func (d *data) WriteTo(w io.Writer) (n int64, err error) {
	return 0, binary.Write(w, binary.LittleEndian, d)
}

func (d *data) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, binary.Read(r, binary.LittleEndian, d)
}

func (s *VariableTestSuite) TestGetStruct() {
	v := Variable[*data]{
		name:         "TestVar",
		guid:         testGuid,
		defaultAttrs: defaultAttrs,
		marshal:      structMarshaller[*data],
		unmarshal:    structUnmarshaller[data],
	}

	s.Run("SetGet", func() {
		var tests = []*testRow[*data]{
			{v, "Empty", "", nil, wantedError(io.EOF)},
			{v, "One", "012345678901", &data{0x2301, 0x01896745}, assert.NoError},
			{v, "Many", "12345678", nil, wantedError(io.ErrUnexpectedEOF)},
		}

		runTests(s.T(), tests, func(t *testing.T, row *testRow[*data], env *testEnv[*data]) {
			if row.wanted == nil {
				if !assert.NoError(t, env.setupVarFile(row)) {
					return
				}
			} else {
				if !env.testSet(row) {
					return
				}
			}
			env.testGet(row)
		})
	})
}

func TestVariablesTestSuite(t *testing.T) {
	suite.Run(t, &VariableTestSuite{})
}
