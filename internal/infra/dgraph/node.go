package dgraph

import (
	"context"
	"losh/internal/core/product/models"
	"losh/internal/repository"
)

var (
	errGetNodeStr    = "failed to get node(s)"
	errSaveNodeStr   = "failed to save node(s)"
	errDeleteNodeStr = "failed to delete node(s)"
)

// GetNode returns a `Node` object by its ID.
func (dr *DgraphRepository) GetNode(id string) (models.Node, error) {
	ctx := context.Background()
	getNode, err := dr.client.GetNode(ctx, id)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetNodeStr).Add("nodeId", id)
	}
	if getNode.GetNode == nil { // not found
		return nil, nil
	}

	var node models.Node
	var copyFrom interface{}
	switch *getNode.GetNode.Typename {
	case "Product":
		node = &models.Product{}
		copyFrom = getNode.GetNode.Product
	case "Component":
		node = &models.Component{}
		copyFrom = getNode.GetNode.Component
	case "License":
		node = &models.License{}
		copyFrom = getNode.GetNode.License
	case "User":
		node = &models.User{}
		copyFrom = getNode.GetNode.User
	case "Group":
		node = &models.Group{}
		copyFrom = getNode.GetNode.Group
	default:
		return nil, repository.WrapRepoError(err, "unsupported type").Add("nodeId", id)
	}
	if err = dr.dataCopier.CopyTo(copyFrom, node); err != nil {
		panic(err)
	}
	return node, nil
}

// // GetNodes returns a list of `Node` objects matching the filter criteria.
// func (dr *DgraphRepository) GetNodes(filter *models.NodeFilter, order *models.NodeOrder, first *int64, offset *int64) ([]*models.Node, error) {
// 	ctx := context.Background()
// 	getNodes, err := dr.client.GetNodes(ctx, filter, order, first, offset)
// 	if err != nil {
// 		return nil, repository.WrapRepoError(err, errGetNodeStr)
// 	}
// 	nodes := make([]*models.Node, 0, len(getNodes.QueryNode))
// 	for _, x := range getNodes.QueryNode {
// 		node := &models.Node{ID: x.ID}
// 		if err = dr.dataCopier.CopyTo(x, node); err != nil {
// 			panic(err)
// 		}
// 		nodes = append(nodes, node)
// 	}
// 	return nodes, nil
// }

// // GetAllNodes returns a list of all `Node` objects.
// func (dr *DgraphRepository) GetAllNodes() ([]*models.Node, error) {
// 	return dr.GetNodes(nil, nil, nil, nil)
// }

// DeleteNode deletes a `Node` object.
func (dr *DgraphRepository) DeleteNode(id *string) error {
	ctx := context.Background()
	delFilter := models.NodeFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteNode(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteNodeStr).Add("nodeId", id)
	}
	return nil
}

// DeleteAllNodes deletes all `Nodes` objects.
func (dr *DgraphRepository) DeleteAllNodes() error {
	return dr.DeleteNode(nil)
}
