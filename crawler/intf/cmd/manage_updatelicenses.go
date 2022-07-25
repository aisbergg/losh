package cmd

import (
	"losh/internal/provider/spdxorg"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"
)

var ManageUpdateLicensesCommand = &gcli.Command{
	Name: "update-licenses",
	Desc: "Download SPDX licenses and update the license database entries",
	Func: func(cmd *gcli.Command, args []string) error {
		cfg, db, err := initConfigAndDatabase(devOptions.Path)
		if err != nil {
			return err
		}

		// download licenses
		licenseProvider := spdxorg.NewSpdxOrgProvider(cfg.Crawler.UserAgent)
		licenses, err := licenseProvider.GetAllLicenses()
		if err != nil {
			return errors.Wrap(err, "failed to download licenses")
		}

		// upload licenses
		err = db.SaveLicenses(licenses)
		if err != nil {
			return errors.Wrap(err, "failed to save licenses")
		}

		return nil
	},
}
