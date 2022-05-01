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

package guid

import (
	"testing"
)

func TestGUID_Braced(t *testing.T) {
	tests := []struct {
		u    GUID
		want string
	}{
		{
			GUID{0x39, 0x85, 0x81, 0xc9, 0x48, 0x77, 0x40, 0x85, 0x81, 0xdc, 0xb3, 0x8d, 0x2a, 0x9f, 0x00, 0x6b},
			"{398581C9-4877-4085-81DC-B38D2A9F006B}",
		},
		{
			GUID{0x74, 0x56, 0xa8, 0x3f, 0x7d, 0x80, 0x40, 0x8b, 0x8a, 0x7d, 0xc1, 0x84, 0x63, 0xc7, 0xec, 0x29},
			"{7456A83F-7D80-408B-8A7D-C18463C7EC29}",
		},
		{
			GUID{0x25, 0xdc, 0x20, 0xc9, 0x04, 0x89, 0x4b, 0x23, 0xab, 0x7b, 0xa3, 0xcb, 0xe1, 0x6c, 0x63, 0xbd},
			"{25DC20C9-0489-4B23-AB7B-A3CBE16C63BD}",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.u.Braced(); got != tt.want {
				t.Errorf("Braced() = %v, want %v", got, tt.want)
			}
		})
	}
}
