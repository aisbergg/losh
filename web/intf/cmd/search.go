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
	"fmt"
	"strings"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/gcli/v3"

	"losh/internal/lib/util/mathutil"
	"losh/internal/lib/util/stringutil"
	"losh/web/core/search"

	searchmodels "losh/web/core/search/models"
)

var searchOptions = struct {
	Path           string
	OrderBy        string
	Descending     bool
	ResultsPerPage int
	Page           int
	Format         string
}{}

// SearchCommand is the CLI command to search for products.
var SearchCommand = &gcli.Command{
	Name: "search",
	Desc: "Search for products",
	Config: func(c *gcli.Command) {
		c.StrOpt(&searchOptions.Path, "config", "c", "", "configuration file path")
		c.StrOpt(&searchOptions.OrderBy, "order", "o", "", "order by")
		c.StrOpt(&searchOptions.Format, "format", "f", "", "export format (accepted values: csv, tsv)")
		c.BoolOpt(&searchOptions.Descending, "descending", "d", false, "descending order")
		c.IntOpt(&searchOptions.ResultsPerPage, "rpp", "n", 100, "results per page")
		c.IntOpt(&searchOptions.Page, "page", "p", 1, "page")
		c.AddArg("queryString", "Search query", true, true)
	},
	Func: func(cmd *gcli.Command, args []string) error {
		_, db, err := initConfigAndDatabase(searchOptions.Path)
		if err != nil {
			return errors.Wrap(err, "failed to load configuration")
		}

		sQry := strings.Join(cmd.Arg("queryString").Strings(), " ")

		// search service
		searchSvc := search.NewService(db, true)
		ctx := context.Background()

		// search
		offset := mathutil.Max(0, (searchOptions.Page-1)*searchOptions.ResultsPerPage)
		first := mathutil.Max(1, searchOptions.ResultsPerPage)
		res, err := searchSvc.Search(
			ctx,
			sQry,
			searchmodels.OrderByFromStr(searchOptions.OrderBy, searchOptions.Descending),
			searchmodels.Pagination{First: first, Offset: offset},
		)
		if err != nil {
			return errors.Wrap(err, "search failed")
		}

		// format output
		format := strings.TrimSpace(strings.ToLower(searchOptions.Format))
		switch format {
		case "csv", "tsv":
			t := searchmodels.ExportTypeCSV
			if format == "tsv" {
				t = searchmodels.ExportTypeTSV
			}
			b, err := searchSvc.ExportResults(
				t,
				ctx,
				sQry,
				searchmodels.OrderByFromStr(searchOptions.OrderBy, searchOptions.Descending),
				searchOptions.ResultsPerPage,
				offset,
			)
			if err != nil {
				return errors.Wrap(err, "export failed")
			}
			fmt.Printf("%s\n", b)
		default:
			fmt.Printf("Number of results: %d\n", res.Count)
			fmt.Printf("Number of retrieved results: %d\n\n", len(res.Items))
			for _, r := range res.Items {
				license := ""
				if r.Release.License != nil {
					license = *r.Release.License.Xid
				}
				fmt.Printf("%s | %s | %s | %s\n", *r.Name, stringutil.Ellipses(*r.Release.Description, 50), license, *r.Release.Repository.URL)
			}
		}

		return nil
	},
}
