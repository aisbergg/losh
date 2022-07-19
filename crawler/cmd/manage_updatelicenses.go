package cmd

import (
	"losh/crawler/config"
	"losh/internal/logging"
	"losh/internal/provider/spdxorg"
	"losh/internal/util/configutil"
	"strings"

	"losh/internal/repository/dgraph"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/gookit/gcli/v3"
)

var ManageUpdateLicensesCommand = &gcli.Command{
	Name: "update-licenses",
	Desc: "Download SPDX licenses and update the license database entries",
	Func: func(cmd *gcli.Command, args []string) error {
		// configuration
		path := pathlib.NewPath(strings.TrimSpace(manageOptions.Path))
		config := config.DefaultConfig()
		err := configutil.Load(path, &config)
		if err != nil {
			return errors.Wrap(err, "failed to load configuration")
		}

		// logging
		err = logging.Initialize(config.Log)
		if err != nil {
			return errors.Wrap(err, "failed to initialize logging")
		}

		// database
		db, err := dgraph.NewDgraphRepository(config.Database)
		if err != nil {
			return errors.Wrap(err, "failed to initialize Dgraph database connection")
		}
		if !db.IsReachable() {
			return errors.New("failed to connect to Dgraph database")
		}

		// download licenses
		licenseProvider := spdxorg.NewSpdxOrgProvider(config.Crawler.UserAgent)
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

// func run() error {
// 	logging.Initialize(logging.AppLogConfig{
// 		Level: "debug",
// 	})
// 	// log := logging.NewLogger("cmd-update-licenses")

// 	// download licenses
// 	licenseProvider := spdxorg.NewSpdxOrgProvider("losh-dev")
// 	licenses, err := licenseProvider.GetAllLicenses()
// 	if err != nil {
// 		return err
// 	}

// 	// l, err := licenseProvider.Get("GPL-2.0-or-later")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// licenses := []models.License{l}

// 	// create dgraph repository
// 	dgraphRepo := dgraph.NewDgraphRepository("http://localhost:8080/graphql")
// 	if !dgraphRepo.IsReachable() {
// 		panic("failed to connect to Dgraph database server")
// 	}

// 	// upload licenses
// 	err = dgraphRepo.SaveLicenses(licenses)
// 	if err != nil {
// 		return eris.Wrap(err, "failed to save licenses")
// 	}

// 	// print number of licenses uploaded
// 	fmt.Printf("%d licenses uploaded\n", len(licenses))

// 	return nil
// }

// func main() {
// 	err := run()
// 	if err != nil {
// 		errStr := eris.ToString(err, true)
// 		fmt.Println(errStr)
// 		os.Exit(1)
// 	}
// }
