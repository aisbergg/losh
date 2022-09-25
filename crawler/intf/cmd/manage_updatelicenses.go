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

	"losh/internal/infra/spdxorg"
	"losh/internal/lib/log"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

// ManageUpdateLicensesCommand is the CLI command to update the licenses.
var ManageUpdateLicensesCommand = &gcli.Command{
	Name: "update-licenses",
	Desc: "Download SPDX licenses and update the license database entries",
	Func: func(cmd *gcli.Command, args []string) error {
		cfg, db, err := initConfigAndDatabase(manageOptions.ConfigPath)
		if err != nil {
			return err
		}

		log := log.NewLogger("cmd")
		log.Info("updating licenses now")

		// download licenses
		licenseProvider := spdxorg.NewSpdxOrgProvider(cfg.Crawler.UserAgent)
		licenses, err := licenseProvider.GetAllLicenses(context.Background())
		if err != nil {
			return errors.Wrap(err, "failed to download licenses")
		}

		// upload licenses
		err = db.CreateLicenses(context.Background(), licenses)
		if err != nil {
			return errors.Wrap(err, "failed to save licenses")
		}

		log.Info("successfully updated licenses")

		return nil
	},
}
