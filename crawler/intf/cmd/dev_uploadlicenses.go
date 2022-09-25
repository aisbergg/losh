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

package cmd

import (
	"context"

	"losh/internal/infra/licensefile"
	"losh/internal/lib/log"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/gookit/gcli/v3"
)

// DevUploadLicensesCommand is the CLI command to upload licenses from license
// files.
var DevUploadLicensesCommand = &gcli.Command{
	Name: "upload-licenses",
	Desc: "Update the license database entries from license files",
	Config: func(c *gcli.Command) {
		c.AddArg("file", "License file to upload", true, false)
	},
	Func: func(cmd *gcli.Command, args []string) error {
		lcsFlePth := pathlib.NewPath(cmd.Arg("file").String())
		if exists, err := lcsFlePth.Exists(); err != nil || !exists {
			return errors.Errorf("license file %s does not exist", lcsFlePth.String())
		}

		_, db, err := initConfigAndDatabase(devOptions.ConfigPath)
		if err != nil {
			return err
		}

		log := log.NewLogger("cmd")
		log.Info("updating licenses now")

		// load licenses
		licenseProvider := licensefile.NewFileRepository(lcsFlePth)
		licenses, err := licenseProvider.GetAllLicenses(context.Background())
		if err != nil {
			return errors.Wrap(err, "failed to load licenses")
		}

		// upload licenses
		err = db.CreateLicenses(context.Background(), licenses)
		if err != nil {
			return errors.Wrap(err, "failed to save licenses")
		}

		log.Info("successfully uploaded licenses")

		return nil
	},
}
