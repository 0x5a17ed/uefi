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

package efidevicepath

import (
	"testing"

	"github.com/0x5a17ed/uefi/efi/efiguid"
)

func TestHardDriveMediaDevicePath_Notation(t *testing.T) {
	type fields struct {
		PartitionNumber    uint32
		PartitionStartLBA  uint64
		PartitionSizeLBA   uint64
		PartitionSignature [16]byte
		PartitionFormat    PartitionFormat
		SignatureType      SignatureType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"mbr, whole disk",
			fields{0, 0, 0, [16]byte{}, PCATPartitionFormat, PCATSignatureType},
			"HD(0,MBR,0x00000000)",
		},
		{
			"mbr, first part",
			fields{1, 0x800, 0x2EE000, [16]byte{0x43, 0x12, 0x02, 0xa0}, PCATPartitionFormat, PCATSignatureType},
			"HD(1,MBR,0xa0021243,0x800,0x2ee000)",
		},
		{
			"gpt, first part",
			fields{1, 0x22, 0x2710000, efiguid.GUID{0x00, 0x9A, 0xE3, 0x15, 0xD2, 0x1D, 0x00, 0x10, 0x8D, 0x7F, 0x00, 0xA0, 0xC9, 0x24, 0x08, 0xFC}, GUIDPartitionFormat, GUIDSignatureType},
			"HD(1,GPT,15E39A00-1DD2-1000-8D7F-00A0C92408FC,0x22,0x2710000)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HardDriveMediaDevicePath{
				PartitionNumber:    tt.fields.PartitionNumber,
				PartitionStartLBA:  tt.fields.PartitionStartLBA,
				PartitionSizeLBA:   tt.fields.PartitionSizeLBA,
				PartitionSignature: tt.fields.PartitionSignature,
				PartitionFormat:    tt.fields.PartitionFormat,
				SignatureType:      tt.fields.SignatureType,
			}
			if got := p.Text(); got != tt.want {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}
