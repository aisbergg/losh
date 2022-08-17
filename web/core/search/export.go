package search

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"losh/internal/core/product/models"

	orderedmap "github.com/wk8/go-ordered-map"
)

type ExportType int

const (
	ExportTypeInvalid ExportType = iota
	ExportTypeJSON
	ExportTypeCSV
	ExportTypeTSV
)

func (s *Service) ExportResults(query string, results Results, t ExportType) ([]byte, error) {
	switch t {
	case ExportTypeJSON:
		return dumpJSON(query, results)
	case ExportTypeCSV, ExportTypeTSV:
		return dumpCSVorTSV(results, t)
	default:
		return nil, &Error{"invalid export type"}
	}
}

func dumpJSON(query string, results Results) ([]byte, error) {
	items := make([]*models.Product, 0, len(results.Items))
	for _, prd := range results.Items {
		items = append(items, models.AsTree(prd).(*models.Product))
	}
	m := orderedmap.NewWithPairs(
		"query", query,
		"count", results.Count,
		"items", items,
	)
	return json.MarshalIndent(m, "", "  ")
}

func dumpCSVorTSV(results Results, t ExportType) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	w.Comma = ';'
	if t == ExportTypeTSV {
		w.Comma = '\t'
	}

	// header
	err := w.Write([]string{
		"name",
		"description",
		"website",
		"forkOf",
		"tags",
		"category",
		"version",
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

	for _, prd := range results.Items {
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
			forkOf,
			strings.Join(tags, ","),
			category,
			pd(prd.Release.Version),
			pd(prd.Release.CreatedAt).String(),
			pd(prd.Release.Repository.URL),
			pd(prd.Release.License.Xid),
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
