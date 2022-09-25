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

package models

import (
	"strings"

	productmodels "losh/internal/core/product/models"
)

type Results struct {
	// Count is the total number of results.
	Count uint64                   `json:"count" liquid:"count"`
	Items []*productmodels.Product `json:"items" liquid:"items"`

	// Operators used in the query (lowercased).
	Operators []string `json:"operators" liquid:"operators"`
}

type ExportType int

const (
	ExportTypeInvalid ExportType = iota
	ExportTypeCSV
	ExportTypeTSV
)

type OrderBy struct {
	Field      OrderByField `json:"field"`
	Descending bool         `json:"ascending"`
}

type OrderByField int

const (
	OrderByName OrderByField = iota
	OrderByDocumentationLanguage
	OrderByState
	OrderByForkCount
	OrderByStarCount
	OrderByVersion
	OrderByWebsite
	OrderByCreatedAt
	OrderByDiscoveredAt
	OrderByLastIndexedAt
	OrderByLastUpdatedAt
	OrderByHasAdditionalLicenses
	OrderByLicenseID
	OrderByLicenseName
	OrderByIsLicenseSpdx
	OrderByIsLicenseDeprecated
	OrderByIsLicenseOsiApproved
	OrderByIsLicenseFsfLibre
	OrderByIsLicenseBlocked
	OrderByLicenseType
	OrderByIsLicenseStrong
	OrderByIsLicenseWeak
	OrderByIsLicensePermissive
	OrderByLicensorFullName
	OrderByLicensorName
	OrderByRepositoryHost
	OrderByRepositoryOwner
	OrderByRepositoryName
	OrderByDatasourceHost
	OrderByDatasourceOwner
	OrderByDatasourceName
	OrderByHasAttestation
	OrderByAttestation
	OrderByHasPublication
	OrderByPublication
	OrderByHasIssueTracker
	OrderByIssueTracker
	OrderByHasComplieswith
	OrderByComplieswith
	OrderByHasCpcpatentclass
	OrderByCpcPatentClass
	OrderByHasTsdc
	OrderByTsdc
	OrderByHasImage
	OrderByImage
	OrderByHasReadme
	OrderByReadme
	OrderByHasContributionGuide
	OrderByContributionGuide
	OrderByHasBom
	OrderByBom
	OrderByHasManufacturingInstructions
	OrderByManufacturingInstructions
	OrderByHasUserManual
	OrderByUserManual
	OrderByHasSource
	OrderBySource
	OrderByHasExport
	OrderByExport
	OrderByHasAuxiliary
	OrderByAuxiliary
)

func OrderByFromStr(s string, descending bool) OrderBy {
	orderBy := OrderBy{Field: OrderByName, Descending: descending}
	switch strings.ToLower(s) {
	case "name":
		orderBy.Field = OrderByName
	case "documentationlanguage", "language":
		orderBy.Field = OrderByDocumentationLanguage
	case "state":
		orderBy.Field = OrderByState
	case "forkcount":
		orderBy.Field = OrderByForkCount
	case "starcount":
		orderBy.Field = OrderByStarCount
	case "version":
		orderBy.Field = OrderByVersion
	case "website":
		orderBy.Field = OrderByWebsite
	case "createdat":
		orderBy.Field = OrderByCreatedAt
	case "discoveredat":
		orderBy.Field = OrderByDiscoveredAt
	case "lastindexedat":
		orderBy.Field = OrderByLastIndexedAt
	case "lastupdatedat":
		orderBy.Field = OrderByLastUpdatedAt
	case "hasadditionallicenses":
		orderBy.Field = OrderByHasAdditionalLicenses
	case "licenseid", "license":
		orderBy.Field = OrderByLicenseID
	case "licensename":
		orderBy.Field = OrderByLicenseName
	case "islicensespdx":
		orderBy.Field = OrderByIsLicenseSpdx
	case "islicensedeprecated":
		orderBy.Field = OrderByIsLicenseDeprecated
	case "islicenseosiapproved":
		orderBy.Field = OrderByIsLicenseOsiApproved
	case "islicensefsflibre":
		orderBy.Field = OrderByIsLicenseFsfLibre
	case "islicenseblocked":
		orderBy.Field = OrderByIsLicenseBlocked
	case "licensetype":
		orderBy.Field = OrderByLicenseType
	case "islicensestrong":
		orderBy.Field = OrderByIsLicenseStrong
	case "islicenseweak":
		orderBy.Field = OrderByIsLicenseWeak
	case "islicensepermissive":
		orderBy.Field = OrderByIsLicensePermissive
	case "licensorfullname", "licensor":
		orderBy.Field = OrderByLicensorFullName
	case "licensorname":
		orderBy.Field = OrderByLicensorName
	case "repositoryhost", "repository", "host":
		orderBy.Field = OrderByRepositoryHost
	case "repositoryowner":
		orderBy.Field = OrderByRepositoryOwner
	case "repositoryname":
		orderBy.Field = OrderByRepositoryName
	case "datasourcehost", "datasource":
		orderBy.Field = OrderByDatasourceHost
	case "datasourceowner":
		orderBy.Field = OrderByDatasourceOwner
	case "datasourcename":
		orderBy.Field = OrderByDatasourceName
	case "hasattestation":
		orderBy.Field = OrderByHasAttestation
	case "attestation":
		orderBy.Field = OrderByAttestation
	case "haspublication":
		orderBy.Field = OrderByHasPublication
	case "publication":
		orderBy.Field = OrderByPublication
	case "hasissuetracker":
		orderBy.Field = OrderByHasIssueTracker
	case "issuetracker":
		orderBy.Field = OrderByIssueTracker
	case "hascomplieswith":
		orderBy.Field = OrderByHasComplieswith
	case "complieswith":
		orderBy.Field = OrderByComplieswith
	case "hascpcpatentclass":
		orderBy.Field = OrderByHasCpcpatentclass
	case "cpcpatentclass":
		orderBy.Field = OrderByCpcPatentClass
	case "hastsdc":
		orderBy.Field = OrderByHasTsdc
	case "tsdc":
		orderBy.Field = OrderByTsdc
	case "hasimage":
		orderBy.Field = OrderByHasImage
	case "image":
		orderBy.Field = OrderByImage
	case "hasreadme":
		orderBy.Field = OrderByHasReadme
	case "readme":
		orderBy.Field = OrderByReadme
	case "hascontributionguide":
		orderBy.Field = OrderByHasContributionGuide
	case "contributionguide":
		orderBy.Field = OrderByContributionGuide
	case "hasbom":
		orderBy.Field = OrderByHasBom
	case "bom":
		orderBy.Field = OrderByBom
	case "hasmanufacturinginstructions":
		orderBy.Field = OrderByHasManufacturingInstructions
	case "manufacturinginstructions":
		orderBy.Field = OrderByManufacturingInstructions
	case "hasusermanual":
		orderBy.Field = OrderByHasUserManual
	case "usermanual":
		orderBy.Field = OrderByUserManual
	case "hassource":
		orderBy.Field = OrderByHasSource
	case "source":
		orderBy.Field = OrderBySource
	case "hasexport":
		orderBy.Field = OrderByHasExport
	case "export":
		orderBy.Field = OrderByExport
	case "hasauxiliary":
		orderBy.Field = OrderByHasAuxiliary
	case "auxiliary":
		orderBy.Field = OrderByAuxiliary
	default:
		orderBy.Field = OrderByName
		orderBy.Descending = false
	}
	return orderBy
}

func OrderByFromCombinedStr(s string) OrderBy {
	if len(s) > 3 {
		s = strings.ToLower(s)
		return OrderByFromStr(s[:len(s)-3], s[len(s)-3:] == "dsc")
	}
	// default to ascending name
	return OrderBy{Field: OrderByName}
}

type Pagination struct {
	First  int `json:"first"`
	Offset int `json:"offset"`
}
