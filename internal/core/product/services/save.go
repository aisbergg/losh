package services

import (
	"context"
	"fmt"
	"reflect"

	"losh/internal/core/product/models"
	"losh/internal/lib/util/reflectutil"
)

// type NodeSet struct {
// 	m *orderedmap.OrderedMap
// }

// // models.NewNodeSet creates a new NodeSet.
// func models.NewNodeSet() *models.NodeSet {
// 	return &NodeSet{
// 		m: orderedmap.New(),
// 	}
// }

// // Len returns the number of elements in the set.
// func (ns *models.NodeSet) Len() int {
// 	return ns.m.Len()
// }

// // add adds a node to the set. If the node is already in the set, it is not
// // added again. The keys are used in order of given priority.
// func (ns *models.NodeSet) add(value interface{}, key ...string) {
// 	if _, ok := ns.get(key...); ok {
// 		return
// 	}
// 	for _, k := range key {
// 		if k == "" {
// 			continue
// 		}
// 		ns.m.Set(k, value)
// 		return
// 	}
// 	return
// }

func s(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// // Add adds a node to the set.
// func (ns *models.NodeSet) Add(node models.Node) {
// 	if node == nil {
// 		return
// 	}
// 	ns.m.Set(node, node)
// }

// // add adds a node to the set. If the node is already in the set, it is not
// // added again. The keys are used in order of given priority.
// func (ns *models.NodeSet) get(key ...string) (models.Node, bool) {
// 	for _, k := range key {
// 		if k == "" {
// 			continue
// 		}
// 		if v, ok := ns.m.Get(k); ok {
// 			return v.(models.Node), true
// 		}
// 	}
// 	return nil, false
// }

// // Get returns a node.
// func (ns *models.NodeSet) Get(node models.Node) (models.Node, bool) {
// 	ns.m.GetPair()
// 	if node == nil {
// 		return nil, false
// 	}
// 	if v, ok := ns.m.Get(node); ok {
// 		return v.(models.Node), true
// 	}
// 	return nil, false
// }

// // RangeReverse works like Range, but in reverse order.
// func (ns *models.NodeSet) RangeReverse(f func(key interface{}, node models.Node) bool) {
// 	fw := func(key, value interface{}) bool {
// 		n, ok := models.AssertNode(value)
// 		if !ok {
// 			panic(fmt.Sprintf("unsupported node type: %T", n))
// 		}
// 		return f(key, value.(models.Node))
// 	}
// 	ns.m.RangeReverse(fw)
// }

// 1. traverse graph
//     1. if visited -> return
//     2. add Node to a 'visited' list
//     5. visit children
// // first pass - create nodes that are not yet in the DB
// 3. For each Node in list reversed:
//     1. has ID
//         1. continue

//     3. try to figure out, if already exists in DB -> get ID (requires knowledge about get query -> map type to get function )
//         1. different for each type -> switch type statement
//         2. type 'xid'
//         3. Host with 'domain'
//         4. function: getNodeID() Node

//     2. if only leaf values
//         1. save values
//     3. else has leaf values
//         1. if all of them have IDs
//             1. save all
//             2. store ID
//         2. else one or more of them don't have an ID
//             1. save mandatory (needs to be recursive)
//                 for each mandatory sub node
//                     1. if sub Node is mandatory
//                         1. mandatory save sub node first (recursion will eventually stop, because they cannot all mandatory require each other)
//                         2. store their ID
//                 1. save mandatory value
//             2. store ID
//             3. add to 'secondPass'

// // second pass - update nodes
// 4. For each Node in Second pass
//     1. add each sub Node and lists of Nodes (non mandatory)

// each save requires knowledge about query - mapping between type and query document

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

	// reorder nodes so that mandatory sub nodes will be saved before their
	// parent node
	processed := make(map[models.Node]struct{}, traversed.Len())
	traversed.RangeReverse(func(node models.Node) bool {
		reorderMandatory(node, traversed, processed)
		return true
	})

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

// func (s *Service) createNode(ctx context.Context, node models.Node) error {
// 	if node.GetID() != "" {
// 		return nil
// 	}

// 	// create mandatory sub nodes first (fields of type struct that are not pointers)
// 	ndeVal := reflect.ValueOf(node)
// 	for i := 0; i < ndeVal.NumField(); i++ {
// 		field := ndeVal.Field(i)
// 		if field.Kind() == reflect.Pointer {
// 			continue
// 		}
// 		if field.Kind() != reflect.Struct || !field.IsValid() {
// 			continue
// 		}
// 		// is expected to implement Node interface
// 		n := field.Interface().(models.Node)
// 		if err := s.saveNode(ctx, n, true); err != nil {
// 			return err
// 		}
// 	}

// 	// create the node itself
// 	return s.saveNode(ctx, node, false)
// }

// // saveMandatory saves mandatory values of a node.
// func (s *Service) saveMandatory(ctx context.Context, node models.Node) error {

// 	switch n := node.(type) {
// 	case *models.Product:
// 		return s.saveMandatoryProduct(ctx, n)

// 	case *models.Component:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.Software:
// 		return ns.get(n.ID, fmt.Sprintf("%v", n))
// 	case *models.Repository:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.TechnologySpecificDocumentationCriteria:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.TechnicalStandard:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.User:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.Group:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.File:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.KeyValue:
// 		return ns.get(n.ID, fmt.Sprintf("%v", n))
// 	case *models.StringV:
// 		return ns.get(n.ID, fmt.Sprintf("%v", n))
// 	case *models.FloatV:
// 		return ns.get(n.ID, fmt.Sprintf("%v", n))
// 	case *models.Material:
// 		return ns.get(n.ID, fmt.Sprintf("%v", n))
// 	case *models.ManufacturingProcess:
// 		return ns.get(n.ID, fmt.Sprintf("%v", n))
// 	case *models.BoundingBoxDimensions:
// 		return ns.get(n.ID, fmt.Sprintf("%v", n))
// 	case *models.OpenSCADDimensions:
// 		return ns.get(n.ID, fmt.Sprintf("%v", n))
// 	case *models.Category:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.Tag:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.License:
// 		return ns.get(n.ID, fmt.Sprintf("%T|%v", n.Xid, n.Xid), fmt.Sprintf("%v", n))
// 	case *models.Host:
// 		return ns.get(n.ID, n.Domain, fmt.Sprintf("%v", n))
// 	default:
// 		panic(fmt.Sprintf("unsupported node type: %T", n))
// 	}
// }
