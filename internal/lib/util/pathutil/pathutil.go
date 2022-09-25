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

package pathutil

import (
	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
)

// GetValidFilePath returns true if path is a valid file path.
func GetValidFilePath(pathStr string) (pathlib.Path, error) {
	path := pathlib.NewPath(pathStr)
	if exists, err := path.Exists(); err != nil || !exists {
		return path, errors.New("file does not exist")
	}
	path, err := path.ResolveAll()
	if err != nil {
		return path, errors.New("failed to resolve path")
	}
	if isFile, err := path.IsFile(); err != nil || !isFile {
		return path, errors.New("given path is not a file")
	}

	return path, nil
}
