package services

import (
	"context"

	"losh/internal/core/product/models"
)

// GetNode returns a `Node` object by its ID.
func (s *Service) GetNode(ctx context.Context, id string) (interface{}, error) {
	return s.repo.GetNode(ctx, id)
}

// getExistingNode tries to check if the node already exists in the DB. If it exists,
// then its ID will be saved in the node object.
// TODO: only get ID and not all the other stuff
func (s *Service) determineID(ctx context.Context, node models.Node) (err error) {
	// we don't need to continue if the node already has an ID
	if node.GetID() != nil {
		return nil
	}

	switch n := node.(type) {
	case *models.Category:
		n.ID, err = s.repo.GetCategoryID(ctx, n.Xid)

	case *models.Component:
		n.ID, err = s.repo.GetComponentID(ctx, n.Xid)

	case *models.File:
		n.ID, err = s.repo.GetFileID(ctx, n.Xid)

	case *models.Group:
		n.ID, err = s.repo.GetGroupID(ctx, n.Xid)

	case *models.Host:
		n.ID, err = s.repo.GetHostID(ctx, n.Domain)

	case *models.License:
		n.ID, err = s.repo.GetLicenseID(ctx, n.Xid)

	case *models.Product:
		n.ID, err = s.repo.GetProductID(ctx, n.Xid)

	case *models.Repository:
		n.ID, err = s.repo.GetRepositoryID(ctx, n.Xid)

	case *models.Tag:
		n.ID, err = s.repo.GetTagID(ctx, n.Name)

	case *models.TechnicalStandard:
		n.ID, err = s.repo.GetTechnicalStandardID(ctx, n.Xid)

	case *models.TechnologySpecificDocumentationCriteria:
		n.ID, err = s.repo.GetTechnologySpecificDocumentationCriteriaID(ctx, n.Xid)

	case *models.User:
		n.ID, err = s.repo.GetUserID(ctx, n.Xid)
	}

	// other types do not have an alternative ID
	return
}
