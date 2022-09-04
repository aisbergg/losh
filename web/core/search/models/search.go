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
	OrderByLastIndexedAt
	OrderByDiscoveredAt
	OrderByCreatedAt
	OrderByUpdatedAt
	OrderByDocumentationLanguage
	OrderByState
	OrderByForkCount
	OrderByStarCount
	OrderByLicense
)

func OrderByFromStr(s string, descending bool) OrderBy {
	orderBy := OrderBy{Field: OrderByName, Descending: descending}
	switch strings.ToLower(s) {
	case "createdat":
		orderBy.Field = OrderByCreatedAt
	case "discoveredat":
		orderBy.Field = OrderByCreatedAt
	case "lastindexedat":
		orderBy.Field = OrderByLastIndexedAt
	case "documentationlanguage":
		orderBy.Field = OrderByDocumentationLanguage
	case "state":
		orderBy.Field = OrderByState
	case "forkcount":
		orderBy.Field = OrderByForkCount
	case "starcount":
		orderBy.Field = OrderByStarCount
	case "license":
		orderBy.Field = OrderByLicense
	default:
		orderBy.Field = OrderByName
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
