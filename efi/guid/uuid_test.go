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

	"github.com/stretchr/testify/require"
)

func TestGUID_Braced(t *testing.T) {
	tests := []struct {
		name    string
		u       GUID
		want    string
		wantErr bool
	}{
		{
			"1",
			GUID{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x03, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			"{00000001-0002-0003-0001-020304050607}",
			false,
		},
		{
			"2",
			GUID{0x61, 0xdf, 0xe4, 0x8b, 0xca, 0x93, 0xd2, 0x11, 0xaa, 0xd, 0x0, 0xe0, 0x98, 0x3, 0x2b, 0x8c},
			"{8BE4DF61-93CA-11D2-AA0D-00E098032B8C}",
			false,
		},
		{
			"3",
			GUID{0x00, 0x9A, 0xE3, 0x15, 0xD2, 0x1D, 0x00, 0x10, 0x8D, 0x7F, 0x00, 0xA0, 0xC9, 0x24, 0x08, 0xFC},
			"{15E39A00-1DD2-1000-8D7F-00A0C92408FC}",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.u.Braced()

			require.Equal(t, tt.want, got)

			gotUU, err := FromString(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.u, gotUU)
		})
	}
}
