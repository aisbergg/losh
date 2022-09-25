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

package search

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph"
	searchmodels "losh/web/core/search/models"
)

const batchSize = 150

func (s *Service) ExportResults(t searchmodels.ExportType, ctx context.Context, queryStr string, orderBy searchmodels.OrderBy, limit, offset int) ([]byte, error) {
	pgdRes := dgraph.NewPaginatedList(
		func(first int, offset int) ([]*models.Product, uint64, error) {
			res, err := s.Search(ctx, queryStr, orderBy, searchmodels.Pagination{First: first, Offset: offset})
			if err != nil {
				return nil, 0, err
			}
			return res.Items, res.Count, nil
		},
		batchSize,
		offset,
	)

	switch t {
	case searchmodels.ExportTypeCSV, searchmodels.ExportTypeTSV:
		return dumpCSVorTSV(pgdRes, t, limit)
	default:
		// should never happen unless we missed something
		panic("invalid export type")
	}
}

func (s *Service) ExportUpTo300Results(t searchmodels.ExportType, ctx context.Context, queryStr string, orderBy searchmodels.OrderBy) ([]byte, error) {
	return s.ExportResults(t, ctx, queryStr, orderBy, 300, 0)
}

func dumpCSVorTSV(results *dgraph.PaginatedList[*models.Product], t searchmodels.ExportType, limit int) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	w.Comma = ';'
	if t == searchmodels.ExportTypeTSV {
		w.Comma = '\t'
	}

	// header
	err := w.Write([]string{
		"name",
		"description",
		"website",
		"state",
		"forkOf",
		"forkCount",
		"starCount",
		"tags",
		"category",
		"version",
		"lastUpdatedAt",
		"createdAt",
		"repository",
		"license",
		"additionalLicenses",
		"licensor",
		"licensorURL",
		"documentationLanguage",
		"technologyReadinessLevel",
		"documentationReadinessLevel",
		"attestation",
		"publication",
		"compliesWith",
		"cpcPatentClass",
		"tsdc",
		"image",
		"readme",
		"contributionGuide",
		"bom",
		"manufacturingInstructions",
		"userManual",
		"source",
		"mass",
		"material",
		"manufacturingProcess",
		"discoveredAt",
		"lastIndexedAt",
	})
	if err != nil {
		return nil, err
	}

	count := 0
	for {
		if count >= limit {
			break
		}
		prd, err := results.Next()
		if err != nil {
			return nil, err
		}
		if prd == nil {
			break
		}
		tags := make([]string, 0, len(prd.Tags))
		for _, tag := range prd.Tags {
			tags = append(tags, *tag.Name)
		}
		var forkOf string
		if prd.ForkOf != nil {
			forkOf = *prd.ForkOf.Release.Repository.URL
		}
		var category string
		if prd.Category != nil {
			category = *prd.Category.FullName
		}
		var image string
		if prd.Release.Image != nil {
			image = *prd.Release.Image.URL
		}
		var additionalLicenses string
		if prd.Release.AdditionalLicenses != nil {
			lcs := make([]string, 0, len(prd.Release.AdditionalLicenses))
			for _, lc := range prd.Release.AdditionalLicenses {
				lcs = append(lcs, *lc.Xid)
			}
			additionalLicenses = strings.Join(lcs, ",")
		}
		var license string
		if prd.Release.License != nil {
			license = *prd.Release.License.Xid
		}
		var licensor string
		var licensorURL string
		if prd.Release.Licensor != nil {
			switch l := prd.Release.Licensor.(type) {
			case *models.User:
				if l.FullName != nil {
					licensor = *l.FullName
				} else {
					licensor = *l.Name
				}
				licensorURL = *l.URL
			case *models.Group:
				if l.FullName != nil {
					licensor = *l.FullName
				} else {
					licensor = *l.Name
				}
				licensorURL = *l.URL
			}
		}
		var compliesWith string
		if prd.Release.CompliesWith != nil {
			compliesWith = *prd.Release.CompliesWith.Name
		}
		var tsdc string
		if prd.Release.Tsdc != nil {
			tsdc = *prd.Release.Tsdc.Name
		}
		var readme string
		if prd.Release.Readme != nil {
			readme = *prd.Release.Readme.URL
		}
		var contributionGuide string
		if prd.Release.ContributionGuide != nil {
			contributionGuide = *prd.Release.ContributionGuide.URL
		}
		var bom string
		if prd.Release.Bom != nil {
			bom = *prd.Release.Bom.URL
		}
		var manufacturingInstructions string
		if prd.Release.ManufacturingInstructions != nil {
			manufacturingInstructions = *prd.Release.ManufacturingInstructions.URL
		}
		var userManual string
		if prd.Release.UserManual != nil {
			userManual = *prd.Release.UserManual.URL
		}
		var source string
		if prd.Release.Source != nil {
			source = *prd.Release.Source.URL
		}
		var material string
		if prd.Release.Material != nil {
			material = *prd.Release.Material.Name
		}
		var manufacturingProcess string
		if prd.Release.ManufacturingProcess != nil {
			manufacturingProcess = *prd.Release.ManufacturingProcess.Name
		}
		var mass string
		if prd.Release.Mass != nil {
			if *prd.Release.Mass > 1000 {
				mass = fmt.Sprintf("%.2f kg", *prd.Release.Mass/1000)
			} else {
				mass = fmt.Sprintf("%.2f g", *prd.Release.Mass)
			}
		}

		err = w.Write([]string{
			pd(prd.Name),
			pd(prd.Release.Description),
			pd(prd.Website),
			string(pd(prd.State)),
			forkOf,
			strconv.FormatInt(pd(prd.ForkCount), 10),
			strconv.FormatInt(pd(prd.StarCount), 10),
			strings.Join(tags, ","),
			category,
			pd(prd.Release.Version),
			pd(prd.LastUpdatedAt).String(),
			pd(prd.Release.CreatedAt).String(),
			pd(prd.Release.Repository.URL),
			license,
			additionalLicenses,
			licensor,
			licensorURL,
			pd(prd.Release.DocumentationLanguage),
			string(pd(prd.Release.TechnologyReadinessLevel)),
			string(pd(prd.Release.DocumentationReadinessLevel)),
			pd(prd.Release.Attestation),
			pd(prd.Release.Publication),
			compliesWith,
			pd(prd.Release.CpcPatentClass),
			tsdc,
			image,
			readme,
			contributionGuide,
			bom,
			manufacturingInstructions,
			userManual,
			source,
			mass,
			material,
			manufacturingProcess,
			prd.DiscoveredAt.String(),
			prd.LastIndexedAt.String(),
		})
		if err != nil {
			return nil, err
		}
		count++
	}

	w.Flush()
	return buf.Bytes(), nil
}

// p returns a pointer to the value. Just for convenience.
func p[T any](v T) *T {
	return &v
}

// pd dereferences a pointer. If the pointer is nil, it returns an empty value.
func pd[T any](v *T) T {
	if v == nil {
		var t T
		return t
	}
	return *v
}
