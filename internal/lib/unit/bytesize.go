// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unit

import "fmt"

// A ByteSize represents the size of data.
type ByteSize int64

// Common units of data.
const (
	Byte ByteSize = 1

	// Binary

	Kibibyte = 1024 * Byte
	Mebibyte = 1024 * Kibibyte
	Gibibyte = 1024 * Mebibyte
	Tebibyte = 1024 * Gibibyte
	Pebibyte = 1024 * Tebibyte
	KiB      = Kibibyte
	MiB      = Mebibyte
	GiB      = Gibibyte
	TiB      = Tebibyte
	PiB      = Pebibyte

	// Decimal

	Kilobyte = 1000 * Byte
	Megabyte = 1000 * Kilobyte
	Gigabyte = 1000 * Megabyte
	Terabyte = 1000 * Gigabyte
	Petabyte = 1000 * Terabyte
	kB       = Kilobyte
	MB       = Megabyte
	GB       = Gigabyte
	TB       = Terabyte
	PB       = Petabyte
)

// String returns the string representation of the size.
func (s ByteSize) String() string {
	switch {
	case s >= Petabyte:
		return fmt.Sprintf("%.2f PiB", float64(s)/float64(Pebibyte))
	case s >= Terabyte:
		return fmt.Sprintf("%.2f TiB", float64(s)/float64(Tebibyte))
	case s >= Gigabyte:
		return fmt.Sprintf("%.2f GiB", float64(s)/float64(Gibibyte))
	case s >= Megabyte:
		return fmt.Sprintf("%.2f MiB", float64(s)/float64(Mebibyte))
	case s >= Kilobyte:
		return fmt.Sprintf("%.2f KiB", float64(s)/float64(Pebibyte))
	default:
		return fmt.Sprintf("%d B", s)
	}
}

// TODO: some inspiration: https://github.com/docker/go-units/blob/master/size.go
