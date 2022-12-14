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

package log

import (
	"io/fs"
)

// Config contains the common log configuration options.
type Config struct {
	// Level defines the logging priority.
	Level string `json:"level" filter:"trim,lower" validate:"in:debug,info,warning,error,critical"`

	// Format defines the encoding format of the logs.
	Format string `json:"format" filter:"trim,lower" validate:"in:console,json"`

	// Filename is the file to write logs to. Backup log files will be retained
	// in the same directory. If empty logs will be written to stderr instead of
	// a file.
	Filename string `json:"filename"`

	// Rotate determines if the file should be automatically rotated with the
	// given parameters.
	Rotate bool `json:"rotate"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated.
	MaxSize int `json:"maxSize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxAge"`

	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int `json:"maxBackups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.
	LocalTime bool `json:"localTime"`

	// Compress determines if the rotated log files should be compressed
	// using gzip.
	Compress bool `json:"compress"`

	// Permission defines the mode which is used to create the file.
	Permissions fs.FileMode `json:"permissions"`
}

// DefaultConfig returns the default configuration for the logger.
func DefaultConfig() Config {
	return Config{
		Level:       "info",
		Format:      "json",
		Filename:    "",
		Rotate:      false,
		MaxSize:     5,
		MaxAge:      0,
		MaxBackups:  0,
		LocalTime:   true,
		Compress:    false,
		Permissions: 0,
	}
}
