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

package binreader

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadUTF16NullString(t *testing.T) {
	tests := []struct {
		name    string
		inp     []byte
		wantOut []byte
		wantErr bool
	}{
		{
			name:    "golden path",
			inp:     []byte{0x61, 0x00, 0x73, 0x00, 0x64, 0x00, 0x00, 0x00},
			wantOut: []byte{0x61, 0x00, 0x73, 0x00, 0x64, 0x00, 0x00, 0x00},
			wantErr: false,
		},
		{
			name:    "extra after string",
			inp:     []byte{0x61, 0x00, 0x73, 0x00, 0x64, 0x00, 0x00, 0x00, 0x01, 0x02},
			wantOut: []byte{0x61, 0x00, 0x73, 0x00, 0x64, 0x00, 0x00, 0x00},
			wantErr: false,
		},
		{
			name:    "unterminated",
			inp:     []byte{0x61, 0x00, 0x73, 0x00, 0x64, 0x00},
			wantOut: []byte{0x61, 0x00, 0x73, 0x00, 0x64, 0x00},
			wantErr: true,
		},
		{
			name:    "short",
			inp:     []byte{0x61, 0x00, 0x00},
			wantOut: []byte{0x61, 0x00},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := ReadUTF16NullBytes(bytes.NewReader(tt.inp))
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUTF16NullBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("ReadUTF16NullBytes() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestUTF16NullBytesToString(t *testing.T) {
	tests := []struct {
		name string
		inp  []byte
		want string
	}{
		{
			"golden path",
			[]byte{0x61, 0x00, 0x73, 0x00, 0x64, 0x00, 0x00, 0x00},
			"asd",
		},
		{
			"extra data past string",
			[]byte{0x61, 0x00, 0x73, 0x00, 0x64, 0x00, 0x00, 0x00, 0x61, 0x00, 0x73, 0x00},
			"asd",
		},
		{
			"data without end mark",
			[]byte{0x61, 0x00, 0x73, 0x00, 0x64, 0x00},
			"asd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UTF16NullBytesToString(tt.inp); got != tt.want {
				t.Errorf("UTF16NullBytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUTF16BytesToString(t *testing.T) {
	tests := []struct {
		name     string
		provided []byte
		expected string
	}{
		{"", []byte{0x74, 0x00, 0x65, 0x00, 0x73, 0x00, 0x74, 0x00}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UTF16BytesToString(tt.provided); got != tt.expected {
				t.Errorf("UTF16BytesToString() = %v, want %v", got, tt.expected)
			}
		})
	}
}
