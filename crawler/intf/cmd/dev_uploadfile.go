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
	"encoding/json"

	"losh/internal/core/product/models"
	"losh/internal/core/product/services"
	"losh/internal/lib/log"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/gookit/gcli/v3"
)

var devUploadFileOptions = struct {
	ConfigPath string
}{}

// DevUploadFile is the CLI command to upload serialized data to the database.
var DevUploadFile = &gcli.Command{
	Name: "upload-file",
	Desc: "Upload serialized data from a file to the database",
	Config: func(c *gcli.Command) {
		c.AddArg("file", "Serialized data file to upload", true, false)
	},
	Func: func(cmd *gcli.Command, args []string) error {
		path := pathlib.NewPath(cmd.Arg("file").String())
		if exists, err := path.Exists(); err != nil || !exists {
			return errors.Errorf("file %s does not exist", path.String())
		}

		_, db, err := initConfigAndDatabase(devOptions.ConfigPath)
		if err != nil {
			return err
		}

		// load licenses
		svc := services.NewService(db)
		err = svc.ReloadLicenseCache()
		if err != nil {
			return errors.Wrap(err, "failed to load licenses")
		}

		log := log.NewLogger("cmd")
		log.Info("uploading data from file now")

		// read file and deserialize to struct
		fileCnt, err := path.ReadFile()
		if err != nil {
			return errors.Wrap(err, "failed to read file")
		}
		prds := []*models.Product{}
		err = json.Unmarshal(fileCnt, &prds)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal file")
		}

		for _, prd := range prds {
			g := models.AsGraph(prd)
			prd := g.(*models.Product)

			err = svc.SaveNode(context.Background(), prd)
			if err != nil {
				return errors.Wrap(err, "failed to save product")
			}
		}

		log.Info("successfully uploaded test data")

		return nil
	},
}
