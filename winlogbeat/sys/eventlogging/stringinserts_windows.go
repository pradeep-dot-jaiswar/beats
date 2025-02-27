// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package eventlogging

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/snappyflow/beats/v7/winlogbeat/sys"
)

const (
	// MaxInsertStrings is the maximum number of strings that can be formatted by
	// FormatMessage API.
	MaxInsertStrings = 99
)

var (
	nullPlaceholder    = []byte{'(', 0, 'n', 0, 'u', 0, 'l', 0, 'l', 0, ')', 0, 0, 0}
	nullPlaceholderPtr = uintptr(unsafe.Pointer(&nullPlaceholder[0]))
)

// StringInserts stores the string inserts for an event, as arrays of string
// and pointer to UTF-16 zero-terminated string suitable to be passed to
// the Windows API. The array of pointers has enough entries to ensure that
// a call to FormatMessage will never crash.
type StringInserts struct {
	pointers [MaxInsertStrings]uintptr
	inserts  []string
	address  uintptr
}

// Parse parses the insert strings from buffer which should contain
// an eventLogRecord.
func (b *StringInserts) Parse(record eventLogRecord, buffer []byte) error {
	if b.inserts == nil { // initialise struct
		b.inserts = make([]string, 0, MaxInsertStrings)
		b.address = reflect.ValueOf(&b.pointers[0]).Pointer()
	}
	b.clear()

	n := int(record.numStrings)
	if n > MaxInsertStrings {
		return fmt.Errorf("number of insert strings in the record (%d) is larger than the limit (%d)", n, MaxInsertStrings)
	}

	b.inserts = b.inserts[:n]
	if n == 0 {
		return nil
	}
	offset := int(record.stringOffset)
	bufferPtr := reflect.ValueOf(&buffer[0]).Pointer()

	for i := 0; i < n; i++ {
		if offset > len(buffer) {
			return fmt.Errorf("Failed reading string number %d, "+
				"offset=%d, len(buffer)=%d, record=%+v", i+1, offset,
				len(buffer), record)
		}
		insertStr, length, err := sys.UTF16BytesToString(buffer[offset:])
		if err != nil {
			return err
		}
		b.inserts[i] = insertStr
		b.pointers[i] = bufferPtr + uintptr(offset)
		offset += length
	}

	return nil
}

// Strings returns the array of strings representing the insert strings.
func (b *StringInserts) Strings() []string {
	return b.inserts
}

// Pointer returns a pointer to an array of UTF-16 strings suitable to be
// passed to the FormatMessage API.
func (b *StringInserts) Pointer() uintptr {
	return b.address
}

func (b *StringInserts) clear() {
	for i := 0; i < MaxInsertStrings && b.pointers[i] != nullPlaceholderPtr; i++ {
		b.pointers[i] = nullPlaceholderPtr
	}
}
