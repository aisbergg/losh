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

package dgraph

import (
	"context"

	"losh/internal/infra/dgraph/dgclient"

	"github.com/fenos/dqlx"
)

var CrawlerMetaDQLFragment = dqlx.QueryBuilder{}.Raw(`
	CrawlerMeta.discoveredAt
	CrawlerMeta.lastIndexedAt
`)

var (
	errGetNodeStr    = "failed to get node(s)"
	errSaveNodeStr   = "failed to save node(s)"
	errDeleteNodeStr = "failed to delete node(s)"
)

// GetNode returns a `Node` object by its ID.
func (dr *DgraphRepository) GetNode(ctx context.Context, id string) (interface{}, error) {
	getNode, err := dr.client.GetNodeByID(ctx, id)
	if err != nil {
		return nil, WrapRepoError(err, errGetNodeStr).Add("nodeId", id)
	}
	if getNode.GetNode == nil {
		return nil, nil
	}

	var node interface{}
	switch *getNode.GetNode.Typename {
	case "Category":
		node, err = dr.GetCategory(ctx, &id, nil)
	case "Component":
		node, err = dr.GetComponent(ctx, &id, nil)
	case "File":
		node, err = dr.GetFile(ctx, &id, nil)
	case "KeyValueFragment":
		node, err = dr.GetKeyValue(ctx, &id)
	case "FloatV":
		node, err = dr.GetFloatV(ctx, &id)
	case "StringV":
		node, err = dr.GetStringV(ctx, &id)
	case "Host":
		node, err = dr.GetHost(ctx, &id, nil)
	case "License":
		node, err = dr.GetLicense(ctx, &id, nil)
	case "ManufacturingProcess":
		node, err = dr.GetManufacturingProcess(ctx, &id)
	case "Material":
		node, err = dr.GetMaterial(ctx, &id)
	case "BoundingBoxDimensions":
		node, err = dr.GetBoundingBoxDimensions(ctx, &id)
	case "OpenSCADDimensions":
		node, err = dr.GetOpenSCADDimensions(ctx, &id)
	case "Product":
		node, err = dr.GetProduct(ctx, &id, nil)
	case "Repository":
		node, err = dr.GetRepository(ctx, &id, nil)
	case "Software":
		node, err = dr.GetSoftware(ctx, &id)
	case "Tag":
		node, err = dr.GetTag(ctx, &id, nil)
	case "TechnicalStandard":
		node, err = dr.GetTechnicalStandard(ctx, &id, nil)
	case "TechnologySpecificDocumentationCriteria":
		node, err = dr.GetTechnologySpecificDocumentationCriteria(ctx, &id, nil)
	case "User":
		node, err = dr.GetUser(ctx, &id, nil)
	case "Group":
		node, err = dr.GetGroup(ctx, &id, nil)
	default:
		return nil, WrapRepoError(err, "unsupported type").Add("nodeId", id)
	}
	if err != nil {
		return nil, err
	}

	return node, nil
}

// CheckNode checks if a `Node` object exists in the DB.
func (dr *DgraphRepository) CheckNode(ctx context.Context, id string) (bool, error) {
	getNode, err := dr.client.CheckNode(ctx, id)
	if err != nil {
		return false, WrapRepoError(err, errGetNodeStr).Add("nodeId", id)
	}
	return getNode.GetNode != nil, nil
}

// DeleteNode deletes a `Node` object.
func (dr *DgraphRepository) DeleteNode(ctx context.Context, id *string) error {
	delFilter := dgclient.NodeFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if _, err := dr.client.DeleteNodes(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteNodeStr).Add("nodeId", id)
	}
	return nil
}

// DeleteAllNodes deletes all `Nodes` objects.
func (dr *DgraphRepository) DeleteAllNodes(ctx context.Context) error {
	return dr.DeleteNode(ctx, nil)
}
