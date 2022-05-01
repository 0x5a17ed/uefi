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
)

func TestACPIDevicePath_Notation(t *testing.T) {
	type fields struct {
		HID uint32
		UID uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"PNP0A03,0", fields{0x0A0341D0, 0}, "ACPI(PNP0A03,0)"},
		{"PNP0A03,1", fields{0x0A0341D0, 1}, "ACPI(PNP0A03,1)"},
		{"PNP0A03,2", fields{0x0A0341D0, 2}, "ACPI(PNP0A03,2)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ACPIPath{
				HID: tt.fields.HID,
				UID: tt.fields.UID,
			}
			if got := p.Text(); got != tt.want {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}
