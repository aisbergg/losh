package services

import (
	"context"
	"fmt"
	"reflect"

	"losh/internal/core/product/models"
	"losh/internal/lib/util/reflectutil"
)

func s(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *Service) SaveNodes(ctx context.Context, node []models.Node) (err error) {
	for _, n := range node {
		if err = s.SaveNode(ctx, n); err != nil {
			return
		}
	}
	return
}

func (s *Service) SaveNode(ctx context.Context, node models.Node) (err error) {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if node == nil {
		return
	}

	// traverse graph and store into set of nodes
	traversed := models.NewNodeSetFromDepthFirst(node)

	traversed.Range(func(node models.Node) bool {
		xid := node.GetAltID()
		if xid != nil {
			return true
		}
		return true
	})

	// check if nodes already exists in the DB
	traversed.Range(func(node models.Node) bool {
		if err = s.determineID(ctx, node); err != nil {
			return false
		}
		return true
	})
	if err != nil {
		return
	}

	// fmt.Println("before")
	// traversed.Range(func(node models.Node) bool {
	// 	fmt.Println("node type", reflect.TypeOf(node))
	// 	return true
	// })

	// reorder nodes so that mandatory sub nodes will be saved before their
	// parent node
	processed := make(map[models.Node]struct{}, traversed.Len())
	traversed.RangeReverse(func(node models.Node) bool {
		reorderMandatory(node, traversed, processed)
		return true
	})

	// fmt.Println("")
	// fmt.Println("after")
	// traversed.Range(func(node models.Node) bool {
	// 	fmt.Println("node type", reflect.TypeOf(node))
	// 	return true
	// })

	// create or update nodes on best effort basis
	updateRequired := models.NewNodeSet()
	traversed.Range(func(node models.Node) bool {
		// exclude license nodes, because they should not be updated here
		if _, ok := node.(*models.License); ok {
			return true
		}

		// if all sub nodes already exists in the DB, save the node directly
		// without the need for an extra pass
		allSubNodesHaveIDs := doAllSubNodesHaveIDs(node)
		err = s.saveNode(ctx, node, allSubNodesHaveIDs)
		if err != nil {
			return false
		}

		// otherwise queue the node for an extra pass
		if !allSubNodesHaveIDs {
			updateRequired.AddBack(node)
		}

		return true
	})
	if err != nil {
		return
	}

	// update, that have newly created sub nodes
	updateRequired.Range(func(node models.Node) bool {
		err = s.saveNode(ctx, node, true)
		if err != nil {
			return false
		}
		return true
	})

	return err
}

func reorderMandatory(node models.Node, reordered *models.NodeSet, processed map[models.Node]struct{}) {
	if node == nil {
		return
	}
	// already processed -> skip
	if _, ok := processed[node]; ok {
		return
	}
	processed[node] = struct{}{}

	// visit children
	nVal := reflect.Indirect(reflect.ValueOf(node))
	flds := reflectutil.GetStructFields(nVal)
	for _, sfFld := range flds {
		if !sfFld.IsMandatory {
			continue
		}
		f := nVal.Field(sfFld.Index)
		fd := reflectutil.Indirect(f)

		switch fd.Kind() {
		case reflect.Invalid:
			continue

		case reflect.Struct:
			n, ok := f.Interface().(models.Node)
			if ok {
				reorderMandatory(n, reordered, processed)
				// reordered.MoveAfter(node, n)
				reordered.MoveBefore(n, node)
			}

		case reflect.Slice:
			// check type of slice
			st := fd.Type().Elem()
			if st.Implements(models.NodeType) {
				for i := 0; i < f.Len(); i++ {
					n := f.Index(i).Interface().(models.Node)
					reorderMandatory(n, reordered, processed)
					// reordered.MoveAfter(node, n)
					reordered.MoveBefore(n, node)
				}
			}
		}
	}
}

// checkAllSubNodesExist checks if all sub nodes of the given node have an ID.
// If not, it will return an error.
func doAllSubNodesHaveIDs(node models.Node) bool {
	ndeVal := reflect.Indirect(reflect.ValueOf(node))
	for i := 0; i < ndeVal.NumField(); i++ {
		fldVal := ndeVal.Field(i)
		fldValInd := reflectutil.Indirect(fldVal)

		switch fldValInd.Kind() {
		case reflect.Invalid:
			continue

		case reflect.Struct:
			n, ok := fldVal.Interface().(models.Node)
			if ok && n.GetID() == nil {
				return false
			}

		case reflect.Slice:
			// check type of slice
			st := fldValInd.Type().Elem()
			if !st.Implements(models.NodeType) {
				continue
			}
			for i := 0; i < fldValInd.Len(); i++ {
				n := fldValInd.Index(i).Interface().(models.Node)
				if n.GetID() == nil {
					return false
				}
			}
		}
	}
	return true
}

// saveNode creates or updates a node in the DB.
func (s *Service) saveNode(ctx context.Context, node models.Node, update bool) error {
	if node.GetID() != nil && !update {
		return nil
	}

	// save the node itself
	switch n := node.(type) {
	case *models.Product:
		if n.ID == nil {
			return s.repo.CreateProduct(ctx, n)
		}
		return s.repo.UpdateProduct(ctx, n)

	case *models.Component:
		if n.ID == nil {
			return s.repo.CreateComponent(ctx, n)
		}
		return s.repo.UpdateComponent(ctx, n)

	case *models.Repository:
		if n.ID == nil {
			return s.repo.CreateRepository(ctx, n)
		}
		return s.repo.UpdateRepository(ctx, n)

	case *models.TechnologySpecificDocumentationCriteria:
		if n.ID == nil {
			return s.repo.CreateTechnologySpecificDocumentationCriteria(ctx, n)
		}
		return s.repo.UpdateTechnologySpecificDocumentationCriteria(ctx, n)

	case *models.TechnicalStandard:
		if n.ID == nil {
			return s.repo.CreateTechnicalStandard(ctx, n)
		}
		return s.repo.UpdateTechnicalStandard(ctx, n)

	case *models.User:
		if n.ID == nil {
			return s.repo.CreateUser(ctx, n)
		}
		return s.repo.UpdateUser(ctx, n)

	case *models.Group:
		if n.ID == nil {
			return s.repo.CreateGroup(ctx, n)
		}
		return s.repo.UpdateGroup(ctx, n)

	case *models.Software:
		if n.ID == nil {
			return s.repo.CreateSoftware(ctx, n)
		}
		return s.repo.UpdateSoftware(ctx, n)

	case *models.File:
		if n.ID == nil {
			return s.repo.CreateFile(ctx, n)
		}
		return s.repo.UpdateFile(ctx, n)

	case *models.KeyValue:
		if n.ID == nil {
			return s.repo.CreateKeyValue(ctx, n)
		}
		return s.repo.UpdateKeyValue(ctx, n)

	case *models.StringV:
		if n.ID == nil {
			return s.repo.CreateStringV(ctx, n)
		}
		return s.repo.UpdateStringV(ctx, n)

	case *models.FloatV:
		if n.ID == nil {
			return s.repo.CreateFloatV(ctx, n)
		}
		return s.repo.UpdateFloatV(ctx, n)

	case *models.Material:
		if n.ID == nil {
			return s.repo.CreateMaterial(ctx, n)
		}
		return s.repo.UpdateMaterial(ctx, n)

	case *models.ManufacturingProcess:
		if n.ID == nil {
			return s.repo.CreateManufacturingProcess(ctx, n)
		}
		return s.repo.UpdateManufacturingProcess(ctx, n)

	case *models.BoundingBoxDimensions:
		if n.ID == nil {
			return s.repo.CreateBoundingBoxDimensions(ctx, n)
		}
		return s.repo.UpdateBoundingBoxDimensions(ctx, n)

	case *models.OpenSCADDimensions:
		if n.ID == nil {
			return s.repo.CreateOpenSCADDimensions(ctx, n)
		}
		return s.repo.UpdateOpenSCADDimensions(ctx, n)

	case *models.Category:
		if n.ID == nil {
			return s.repo.CreateCategory(ctx, n)
		}
		return s.repo.UpdateCategory(ctx, n)

	case *models.Tag:
		if n.ID == nil {
			return s.repo.CreateTag(ctx, n)
		}
		return s.repo.UpdateTag(ctx, n)

	case *models.License:
		// don't update license
		return nil

	case *models.Host:
		if n.ID == nil {
			return s.repo.CreateHost(ctx, n)
		}
		return s.repo.UpdateHost(ctx, n)

	default:
		panic(fmt.Sprintf("unsupported node type: %T", n))
	}
}
