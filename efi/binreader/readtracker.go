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
	"io"
)

type ReadTracker struct {
	reader io.Reader
	offset *int64
}

func NewReadTracker(reader io.Reader, offset *int64) *ReadTracker {
	if offset == nil {
		offset = new(int64)
	}
	return &ReadTracker{reader: reader, offset: offset}
}

func (t *ReadTracker) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)
	*t.offset += int64(n)
	return
}

func (t *ReadTracker) Offset() int64 {
	return *t.offset
}
