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

	"losh/internal/core/product/models"
	"losh/internal/infra/licensefile"
	"losh/internal/infra/spdxorg"
	"losh/internal/lib/log"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/gookit/gcli/v3"
)

// DevDownloadLicensesCommand is the CLI command to download specific licenses
// from SPDX.org and save them to a file.
var DevDownloadLicensesCommand = &gcli.Command{
	Name: "download-licenses",
	Desc: "Download specific licenses from SPDX.org and save them to a file",
	Config: func(c *gcli.Command) {
		c.AddArg("file", "File to save to", true, false)
		c.AddArg("spdx-ids", "SPDX IDs of licenses to download", true, true)
	},
	Func: func(cmd *gcli.Command, args []string) error {
		lcsFlePth := pathlib.NewPath(cmd.Arg("file").String())
		spdxIDs := cmd.Arg("spdx-ids").Strings()

		cfg, err := initConfig(devOptions.ConfigPath)
		if err != nil {
			return errors.Wrap(err, "failed to load configuration")
		}

		// logging
		err = log.Initialize(cfg.Log)
		if err != nil {
			return errors.Wrap(err, "failed to initialize logging")
		}

		log := log.NewLogger("cmd")
		log.Info("downloading licenses now")

		spdx := spdxorg.NewSpdxOrgProvider(cfg.Crawler.UserAgent)
		lcss := make([]*models.License, 0, len(spdxIDs))
		for _, spdxID := range spdxIDs {
			lcs, err := spdx.GetLicense(context.Background(), nil, &spdxID)
			if err != nil {
				return errors.Wrap(err, "failed to download licenses")
			}
			lcss = append(lcss, lcs)
		}

		// save licenses
		lcsRepo := licensefile.NewFileRepository(lcsFlePth)
		err = lcsRepo.SaveLicenses(context.Background(), lcss)
		if err != nil {
			return errors.Wrap(err, "failed to save licenses")
		}

		log.Info("successfully downloaded and saved licenses")

		return nil
	},
}
