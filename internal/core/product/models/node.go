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
	"container/list"
	"fmt"
	"reflect"
	"strings"

	"losh/internal/lib/util/reflectutil"
)

type CrawlerMeta interface {
	IsCrawlerMeta()
}

type Node interface {
	IsNode()
	GetID() *string
	GetAltID() *string
}

// -----------------------------------------------------------------------------

// AssertNode performs a type assertion on the given interface and returns the
// Node interface.
func AssertNode(value interface{}) (Node, bool) {
	if value == nil {
		return nil, false
	}
	if n, ok := value.(Node); ok {
		return n, true
	}
	rval := reflect.ValueOf(value)
	if rval.Kind() == reflect.Interface {
		if rval.IsNil() {
			return nil, false
		}
		if n, ok := rval.Elem().Interface().(Node); ok {
			return n, true
		}
	}
	return nil, false
}

var NodeType = reflect.TypeOf((*Node)(nil)).Elem()

type NodeSet struct {
	elms map[Node]struct{}
	list *list.List
}

// NewNodeSet creates a new NodeSet.
func NewNodeSet() *NodeSet {
	return &NodeSet{
		elms: make(map[Node]struct{}),
		list: list.New(),
	}
}

func NewNodeSetFromDepthFirst(n Node) *NodeSet {
	ns := NewNodeSet()
	traverseDepthFirstRec(n, ns)
	return ns
}

func traverseDepthFirstRec(node Node, visited *NodeSet) {
	if node == nil {
		return
	}

	// check if already visited - if not, add to visited and continue
	if visited.In(node) {
		return
	}
	visited.AddFront(node)

	// visit children
	v := reflectutil.Indirect(reflect.ValueOf(node))
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fd := reflectutil.Indirect(f)

		switch fd.Kind() {
		case reflect.Invalid:
			continue

		case reflect.Struct:
			n, ok := f.Interface().(Node)
			if ok {
				traverseDepthFirstRec(n, visited)
			}

		case reflect.Slice:
			// check type of slice
			st := fd.Type().Elem()
			if st.Implements(NodeType) {
				for i := 0; i < f.Len(); i++ {
					n := f.Index(i).Interface().(Node)
					traverseDepthFirstRec(n, visited)
				}
			}
		}
	}
}

func NewNodeSetFromBreadthFirst(node Node) *NodeSet {
	if node == nil {
		return nil
	}
	queue := NewNodeSet()
	traversed := NewNodeSet()

	// add start node to queue
	queue.AddBack(node)

	// traverse breadth first
	for queue.Len() > 0 {
		n := queue.PopFront()
		traversed.AddBack(n)

		// add children to queue
		rv := reflect.ValueOf(n)
		if rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				continue
			}
			rv = rv.Elem()
		}
		for i := 0; i < rv.NumField(); i++ {
			f := rv.Field(i)
			fd := reflectutil.Indirect(f)

			switch fd.Kind() {
			case reflect.Invalid:
				continue

			case reflect.Struct:
				n, ok := f.Interface().(Node)
				if ok && !traversed.In(n) {
					queue.AddBack(n)
				}

			case reflect.Slice:
				// check type of slice
				st := fd.Type().Elem()
				if st.Implements(NodeType) {
					for i := 0; i < f.Len(); i++ {
						n := f.Index(i).Interface().(Node)
						if !traversed.In(n) {
							queue.AddBack(n)
						}
					}
				}
			}
		}
	}

	return traversed
}

// String returns a string representation of the node set.
func (ns *NodeSet) String() string {
	var b strings.Builder
	b.WriteString("{")
	ns.Range(func(n Node) bool {
		b.WriteString(fmt.Sprintf("%T ", n))
		return true
	})
	b.WriteString("}")
	return b.String()
}

// In returns true if the given node is in the set.
func (ns *NodeSet) In(key Node) bool {
	_, present := ns.elms[key]
	return present
}

// AddFront adds a node to the front of the set if it is not already in the set.
func (ns *NodeSet) AddFront(node Node) {
	if _, present := ns.elms[node]; present {
		return
	}
	ns.list.PushFront(node)
	ns.elms[node] = struct{}{}
}

// AddBack adds a node to the back of the set if it is not already in the set.
func (ns *NodeSet) AddBack(node Node) {
	if _, present := ns.elms[node]; present {
		return
	}
	ns.list.PushBack(node)
	ns.elms[node] = struct{}{}
}

// PopFront removes and returns the front node.
func (ns *NodeSet) PopFront() Node {
	if ns.list.Len() == 0 {
		return nil
	}
	n := ns.list.Front()
	ns.list.Remove(n)
	delete(ns.elms, n.Value.(Node))
	return n.Value.(Node)
}

// PopBack removes and returns the back node.
func (ns *NodeSet) PopBack() Node {
	if ns.list.Len() == 0 {
		return nil
	}
	n := ns.list.Back()
	ns.list.Remove(n)
	delete(ns.elms, n.Value.(Node))
	return n.Value.(Node)
}

// Len returns the length of the ordered map.
func (ns *NodeSet) Len() int {
	return len(ns.elms)
}

// Range calls f sequentially for each node in the set. If f returns false,
// Range stops the iteration.
func (ns *NodeSet) Range(f func(n Node) bool) {
	for elm := ns.list.Front(); elm != nil; elm = elm.Next() {
		n := elm.Value.(Node)
		if !f(n) {
			return
		}
	}
}

// RangeReverse works like Range, but in reverse order.
func (ns *NodeSet) RangeReverse(f func(n Node) bool) {
	for elm := ns.list.Back(); elm != nil; elm = elm.Prev() {
		n := elm.Value.(Node)
		if !f(n) {
			return
		}
	}
}

// MoveBefore moves the node before the markNode if it isn't already.
func (ns *NodeSet) MoveBefore(node, markNode Node) error {
	_, present := ns.elms[node]
	if !present {
		panic("node not in set")
	}
	_, present = ns.elms[markNode]
	if !present {
		panic("node not in set")
	}
	i, in, im := 0, -1, -1
	var ln *list.Element
	var lm *list.Element
	for elm := ns.list.Front(); elm != nil; elm = elm.Next() {
		if elm.Value == node {
			in = i
			ln = elm
		}
		if elm.Value == markNode {
			im = i
			lm = elm
		}
		if in != -1 && im != -1 {
			break
		}
		i++
	}
	if in > im {
		ns.list.MoveBefore(ln, lm)
	}
	return nil
}

// MoveAfter moves the node after the markNode if it isn't already.
func (ns *NodeSet) MoveAfter(node, markNode Node) error {
	_, present := ns.elms[node]
	if !present {
		panic("node not in set")
	}
	_, present = ns.elms[markNode]
	if !present {
		panic("node not in set")
	}
	i, in, im := 0, -1, -1
	var ln *list.Element
	var lm *list.Element
	for elm := ns.list.Front(); elm != nil; elm = elm.Next() {
		if elm.Value == node {
			in = i
			ln = elm
		}
		if elm.Value == markNode {
			im = i
			lm = elm
		}
		if in != -1 && im != -1 {
			break
		}
		i++
	}
	if in < im {
		ns.list.MoveAfter(ln, lm)
	}
	return nil
}

// AsGraph returns a copy of the node as a graph with cyclic references. For
// example:
//
//   &Product{
//     Xid: "foo"
//     Owner: &User{
//       Xid: "bar",
//       Name: "Bar",
//       Products: []*Product{
//         &Product{ // new node with just an ID
//           Xid: "foo",
//         }
//       }
//     }
//   }
//
// will be turned into:
//
//   &Product{
//     Xid: "foo"
//     Owner: &User{
//       Xid: "bar",
//       Name: "Bar",
//       Products: []*Product{&Product{...}} // cyclic reference
//     }
//   }
func AsGraph(tree Node) Node {
	if tree == nil {
		return nil
	}
	traversed := NewNodeSet()
	processed := make(map[string]Node)
	return asGraphRec(tree, traversed, processed)
}

func asGraphRec(node Node, visited *NodeSet, processed map[string]Node) Node {
	if node == nil {
		return nil
	}
	// graph := reflectutil.ZeroValue(reflect.TypeOf(node)).Interface().(Node)

	// node
	// create new node with same kind
	//   if same: return existing and continue; same ptr, ID, Type + AltID
	//   else: create new node and continue
	// add node to traversed
	// add new node to processed

	// for each field
	// node: recurse
	// non-node: copy

	// add to processed

	// getFromProcessed := func(node Node) (Node, bool) {
	// 	if node == nil {
	// 		return nil, false
	// 	}
	// 	id := node.GetID()
	// 	if id != nil {
	// 		if n, present := processed[*id]; present {
	// 			return n, true
	// 		}
	// 	}
	// 	altID := node.GetAltID()
	// 	if altID != nil {
	// 		key := fmt.Sprintf("%T:%s", node, *altID)
	// 		if n, present := processed[key]; present {
	// 			return n, true
	// 		}
	// 	}
	// 	return nil, false
	// }

	// addToPorcessed := func(node Node) {
	// 	if node == nil {
	// 		return
	// 	}
	// 	id := node.GetID()
	// 	if id != nil {
	// 		if _, present := processed[*id]; !present {
	// 			processed[*id] = node
	// 		}
	// 	}
	// 	altID := node.GetAltID()
	// 	if altID != nil {
	// 		key := fmt.Sprintf("%T:%s", node, *altID)
	// 		if _, present := processed[key]; !present {
	// 			processed[key] = node
	// 		}
	// 	}
	// 	processed[fmt.Sprint(node)] = node
	// }

	// create new node of same kind
	on := node
	t := reflect.TypeOf(node)
	nn, ok := getFromProcessed(processed, node.GetID(), node.GetAltID(), t)
	if !ok {
		nn = reflectutil.ZeroFromValue(reflect.ValueOf(on)).Interface().(Node)
		addToProcessed(processed, nn, node.GetID(), node.GetAltID(), t)
	}

	// check if already visited - if not, add to visited and continue
	if visited.In(node) {
		return nn
	}
	visited.AddBack(node)

	// visit children
	ov := reflect.ValueOf(node)     // original value
	ovd := reflectutil.Indirect(ov) // original value dereferenced
	nv := reflect.ValueOf(nn)       // new value
	nvd := reflectutil.Indirect(nv) // new value dereferenced
	for i := 0; i < ovd.NumField(); i++ {
		of := ovd.Field(i)              // original field
		ofd := reflectutil.Indirect(of) // original field dereferenced
		nf := nvd.Field(i)              // new field
		nfd := reflectutil.Indirect(nf) // new field dereferenced

		switch ofd.Kind() {
		case reflect.Invalid:
			continue

		case reflect.Struct:
			n, ok := of.Interface().(Node)
			if ok {
				nv := asGraphRec(n, visited, processed)
				// keep values that are already set
				if nv != nil && !nfd.IsValid() {
					nf.Set(reflect.ValueOf(nv))
				}
				continue
			}

		case reflect.Slice:
			// check type of slice
			st := ofd.Type().Elem()
			if st.Implements(NodeType) {
				if ofd.Len() > 0 {
					nslv := reflectutil.ZeroFromValue(ofd) // new slice value
					nslvd := reflectutil.Indirect(nslv)    // new slice value dereferenced
					nslvd.Set(reflect.MakeSlice(ofd.Type(), 0, ofd.Len()))

					// iterate over slice
					for i := 0; i < ofd.Len(); i++ {
						elm := of.Index(i)
						node := elm.Interface().(Node)
						n := asGraphRec(node, visited, processed)
						nslvd.Set(reflect.Append(nslvd, reflect.ValueOf(n)))
					}
					nf.Set(nslv)
				}
				continue
			}
			// not of type Node -> use as is
		}

		// set directly
		if !nfd.IsValid() {
			nf.Set(of)
		}
	}

	return nn
}

func getFromProcessed(processed map[string]Node, id, altID *string, t reflect.Type) (Node, bool) {
	if id != nil {
		if n, present := processed[*id]; present {
			return n, true
		}
	}
	if altID != nil {
		key := fmt.Sprintf("%s:%s", t, *altID)
		if n, present := processed[key]; present {
			return n, true
		}
	}
	return nil, false
}

func addToProcessed(processed map[string]Node, node Node, id, altID *string, t reflect.Type) {
	if id != nil {
		if _, present := processed[*id]; !present {
			processed[*id] = node
		}
	}
	if altID != nil {
		key := fmt.Sprintf("%s:%s", t, *altID)
		if _, present := processed[key]; !present {
			processed[key] = node
		}
	}
	processed[fmt.Sprint(node)] = node
}

// AsTree returns a copy of the node as kind of a tree without any cyclic
// references. For example:
//
//   &Product{
//     Xid: "foo"
//     Owner: &User{
//       Xid: "bar",
//       Name: "Bar",
//       Products: []*Product{&Product{...}} // cyclic reference
//     }
//   }
//
// will be turned into:
//
//   &Product{
//     Xid: "foo"
//     Owner: &User{
//       Xid: "bar",
//       Name: "Bar",
//       Products: []*Product{
//         &Product{ // new node with just an ID
//           Xid: "foo",
//         }
//       }
//     }
//   }
func AsTree(graph Node) Node {
	if graph == nil {
		return nil
	}
	queue := NewNodeQueue()
	queueCopies := NewNodeQueue()
	traversed := NewNodeSet()
	tree := reflectutil.ZeroFromValue(reflect.ValueOf(graph)).Interface().(Node)

	// add start node to queue
	queue.Push(graph)
	queueCopies.Push(tree)

	// traverse breadth first
	for queue.Len() > 0 {
		// pop values from queue
		nOrg := queue.Pop()
		nCpy := queueCopies.Pop()

		ov := reflect.ValueOf(nOrg) // original node value
		ovd := reflect.Indirect(ov) // original node value dereferenced
		nv := reflect.ValueOf(nCpy) // new node value
		nvd := reflect.Indirect(nv) // new node value dereferenced

		// already traversed -> copy ID only
		if traversed.In(nOrg) {
			flds := reflectutil.GetStructFields(ovd)
			for _, fld := range flds {
				fv := ovd.Field(fld.Index)
				if fv.IsNil() {
					continue
				}
				if fld.IsID || fld.IsAltID {
					nvd.Field(fld.Index).Set(fv)
				}
			}
			continue
		}

		// consider popped node as traversed
		traversed.AddBack(nOrg)

		// add children to queue
		for i := 0; i < ovd.NumField(); i++ {
			of := ovd.Field(i)              // original node field
			ofd := reflectutil.Indirect(of) // original node field dereferenced
			nf := nvd.Field(i)              // new node field

			switch ofd.Kind() {
			// skip null fields
			case reflect.Invalid:
				continue

			case reflect.Struct:
				node, ok := of.Interface().(Node)
				if ok {
					nfv := reflectutil.ZeroFromValue(of)
					// add to queue for later processing
					queue.Push(node)
					queueCopies.Push(nfv.Interface().(Node))
					nf.Set(nfv)
					continue
				}
				// not of type Node -> use as is

			case reflect.Slice:
				// check type of slice
				st := ofd.Type().Elem()
				if st.Implements(NodeType) {
					if ofd.Len() > 0 {
						nslv := reflectutil.ZeroFromValue(ofd) // new slice value
						nslvd := reflectutil.Indirect(nslv)    // new slice value dereferenced
						nslvd.Set(reflect.MakeSlice(ofd.Type(), 0, ofd.Len()))

						// iterate over slice
						for i := 0; i < ofd.Len(); i++ {
							elm := of.Index(i)

							node := elm.Interface().(Node)
							nfv := reflectutil.ZeroFromValue(elm)
							// add to queue for later processing
							queue.Push(node)
							queueCopies.Push(nfv.Interface().(Node))
							nslvd.Set(reflect.Append(nslvd, nfv))
						}
						nf.Set(nslv)
					}
					continue
				}
				// not of type Node -> use as is
			}

			// copy field as is
			nf.Set(reflectutil.ZeroFromValue(of))
			nfd := reflectutil.Indirect(nf)
			nfd.Set(ofd)
		}
	}

	return tree
}

// -----------------------------------------------------------------------------
//
// NodeQueue
//
// -----------------------------------------------------------------------------

type NodeQueue struct {
	nodes list.List
}

// NewNodeQueue creates a new node queue.
func NewNodeQueue() *NodeQueue {
	return &NodeQueue{}
}

// Len returns the length of the queue.
func (q *NodeQueue) Len() int {
	return q.nodes.Len()
}

// Push adds a node to the back of the queue.
func (q *NodeQueue) Push(node Node) {
	q.nodes.PushBack(node)
}

// Pop removes and returns the front node.
func (q *NodeQueue) Pop() Node {
	if q.nodes.Len() == 0 {
		return nil
	}
	n := q.nodes.Front()
	q.nodes.Remove(n)
	return n.Value.(Node)
}

// -----------------------------------------------------------------------------
//
// Traverse Functions
//
// -----------------------------------------------------------------------------

type Queue struct {
	items list.List
}

// NewQueue creates a new node queue.
func NewQueue() *Queue {
	return &Queue{}
}

// Len returns the length of the queue.
func (q *Queue) Len() int {
	return q.items.Len()
}

// Push adds a node to the back of the queue.
func (q *Queue) Push(v interface{}) {
	q.items.PushBack(v)
}

// Pop removes and returns the front node.
func (q *Queue) Pop() interface{} {
	if q.items.Len() == 0 {
		return nil
	}
	n := q.items.Front()
	q.items.Remove(n)
	return n.Value
}

type TraverseFunc func(key []string, value interface{}) error

func traverseBreadthFirst(graph Node, fn TraverseFunc) error {
	return traverseBreadthFirstRec(graph, fn, []string{})
}

type traversedNode struct {
	key  []string
	node Node
}

// appendNew creates a new slice (with new backing array) and appends the
// elements to the end of it.
func appendNew[T any](slice []T, elms ...T) []T {
	new := make([]T, len(slice)+len(elms))
	copy(new, slice)
	for i, elm := range elms {
		new[len(slice)+i] = elm
	}
	return new
}

func traverseBreadthFirstRec(graph Node, fn TraverseFunc, key []string) error {
	if graph == nil {
		return nil
	}
	queue := NewQueue()
	traversed := NewNodeSet()

	// add start node to queue
	queue.Push(traversedNode{[]string{}, graph})

	// traverse breadth first
	for queue.Len() > 0 {
		// pop values from queue
		t := queue.Pop().(traversedNode)

		// consider popped node as traversed
		traversed.AddBack(t.node)

		v := reflect.ValueOf(t.node) // node value
		vd := reflect.Indirect(v)    // node value dereferenced

		// add children to queue
		for i := 0; i < vd.NumField(); i++ {
			fldName := vd.Type().Field(i).Name
			fld := vd.Field(i)
			fldDrf := reflectutil.Indirect(fld)

			switch fldDrf.Kind() {
			// skip null fields
			case reflect.Invalid:
				continue

			case reflect.Struct:
				node, ok := fld.Interface().(Node)
				if ok {
					// add to queue for later processing
					queue.Push(traversedNode{appendNew(key, fldName), node})
					continue
				}
				// not of type Node -> do not go deeper

			case reflect.Slice:
				// check type of slice
				st := fldDrf.Type().Elem()
				if st.Implements(NodeType) {
					if fldDrf.Len() > 0 {
						// iterate over slice
						for i := 0; i < fldDrf.Len(); i++ {
							elm := fld.Index(i)

							node := elm.Interface().(Node)
							// add to queue for later processing
							queue.Push(node)
						}
					}
					continue
				}
				// not of type Node -> do not go deeper
			}

			// leaf -> call traverse function
			err := fn(key, fldDrf.Interface())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
