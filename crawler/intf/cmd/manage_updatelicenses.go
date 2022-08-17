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
