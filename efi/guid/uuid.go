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
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/0x5a17ed/uefi/efi/efihex"
)

var (
	ErrBadFormat = errors.New("incorrect format")
	ErrBadLength = errors.New("incorrect length")
)

func encodeHex(dst []byte, uuid GUID) {
	_ = dst[len(uuid)-1] // bounds check hint to compiler
	efihex.Encode(dst, uuid[:4])
	dst[8] = '-'
	efihex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	efihex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	efihex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	efihex.Encode(dst[24:], uuid[10:])
}

type GUID [16]byte

func (u GUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], u)
	return strings.ToUpper(string(buf[:]))
}

func (u GUID) Braced() string {
	var buf [38]byte
	buf[0] = '{'
	encodeHex(buf[1:], u)
	buf[37] = '}'
	return strings.ToUpper(string(buf[:]))
}

// decodeCanonical decodes GUID string in canonical format i.e
// "973a15e0-fa54-4fef-8d93-aebf7ace5d6b".
func (u *GUID) decodeCanonical(b []byte) (err error) {
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return fmt.Errorf("uuid %q: %w", b, ErrBadFormat)
	}

	src := b[:]
	dst := u[:]

	for i, byteGroup := range []int{8, 4, 4, 4, 12} {
		if i > 0 {
			src = src[1:] // skip dash
		}
		_, err = hex.Decode(dst[:byteGroup/2], src[:byteGroup])
		if err != nil {
			return
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	return
}

func (u *GUID) UnmarshalText(b []byte) (err error) {
	switch len(b) {
	case 36:
		return u.decodeCanonical(b)
	default:
		return fmt.Errorf("uuid %q: %w", b, ErrBadLength)
	}
}

func FromString(t string) (u GUID, err error) {
	err = u.UnmarshalText([]byte(t))
	return
}

// Must wrap a call to a function returning (GUID, error) and panics
// if the error is non-nil. Its intended use is for variable
// initializations.
func Must(u GUID, err error) GUID {
	if err != nil {
		panic(err)
	}
	return u
}

func MustFromString(s string) GUID {
	return Must(FromString(s))
}
