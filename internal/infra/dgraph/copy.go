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
	"reflect"
	"time"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/util/reflectutil"

	"github.com/aisbergg/go-copier/pkg/copier"
)

func copierConverters() []copier.TypeConverter {
	return []copier.TypeConverter{
		copier.TypeConverter{
			SrcType: dgclient.UserOrGroupBasicFragment{},
			DstType: reflect.TypeOf((*models.UserOrGroup)(nil)).Elem(),
			Fn: func(src interface{}, _ *copier.Copier) (interface{}, bool, error) {
				o := src.(dgclient.UserOrGroupBasicFragment)
				switch *o.Typename {
				case "User":
					return &models.User{}, false, nil
				case "Group":
					return &models.Group{}, false, nil
				}
				panic("unexpected type")
			},
		},
		copier.TypeConverter{
			SrcType: dgclient.UserOrGroupFragment{},
			DstType: reflect.TypeOf((*models.UserOrGroup)(nil)).Elem(),
			Fn: func(src interface{}, _ *copier.Copier) (interface{}, bool, error) {
				o := src.(dgclient.UserOrGroupFragment)
				switch *o.Typename {
				case "User":
					return &models.User{}, false, nil
				case "Group":
					return &models.Group{}, false, nil
				}
				panic("unexpected type")
			},
		},
		copier.TypeConverter{
			SrcType: dgclient.UserOrGroupFullFragment{},
			DstType: reflect.TypeOf((*models.UserOrGroup)(nil)).Elem(),
			Fn: func(src interface{}, copier *copier.Copier) (interface{}, bool, error) {
				o := src.(dgclient.UserOrGroupFullFragment)
				switch *o.Typename {
				case "User":
					u := &models.User{}
					copier.CopyTo(o.User, u)
					return u, true, nil
				case "Group":
					g := &models.Group{}
					copier.CopyTo(o.Group, g)
					return g, true, nil
				}
				panic("unexpected type")
			},
		},
	}
}

func stringListContains(list []interface{}, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func dqlCopierConverters() []copier.TypeConverter {
	return []copier.TypeConverter{
		copier.TypeConverter{
			SrcType: map[string]interface{}{},
			DstType: reflect.TypeOf((*models.UserOrGroup)(nil)).Elem(),
			Fn: func(src interface{}, _ *copier.Copier) (interface{}, bool, error) {
				o := src.(map[string]interface{})
				types := o["dgraph.type"].([]interface{})
				if stringListContains(types, "User") {
					return &models.User{}, false, nil
				} else if stringListContains(types, "Group") {
					return &models.Group{}, false, nil
				}
				panic("unexpected type")
			},
		},
	}
}

// initMandatoryFields initializes mandatory fields of type slice. Slices are
// nullable and therefore must be initialized.
func (dr *DgraphRepository) initMandatorySlices(node models.Node) {
	if node == nil {
		return
	}

	nodeVal := reflectutil.Indirect(reflect.ValueOf(node))
	flds := reflectutil.GetStructFields(nodeVal)
	for _, fld := range flds {
		if fld.IsMandatory {
			continue
		}
		fldVal := nodeVal.Field(fld.Index)
		if fldVal.Kind() != reflect.Slice {
			continue
		}
		if fldVal.IsNil() {
			fldVal.Set(reflect.MakeSlice(fldVal.Type(), 0, 0))
		}
	}
}

func (dr *DgraphRepository) copyORMStruct(src, dst interface{}) {
	dstVal := reflect.ValueOf(dst)
	dstValInd := reflectutil.Indirect(dstVal)
	if dstValInd.Kind() != reflect.Struct || !dstValInd.IsValid() {
		panic("copy destination is invalid")
	}

	srcVal := reflect.ValueOf(src)
	srcValInd := reflectutil.Indirect(srcVal)
	if srcValInd.Kind() != reflect.Struct || !srcValInd.IsValid() {
		panic("copy source is invalid")
	}

	srcFlds := reflectutil.GetStructFields(srcValInd)
	dstFlds := reflectutil.GetStructFields(dstValInd)

	for k, srcFld := range srcFlds {
		dstFld, ok := dstFlds[k]
		if !ok {
			continue
		}
		srcFldVal := srcValInd.Field(srcFld.Index)
		dstFldVal := dstValInd.Field(dstFld.Index)
		dr.copyValue(srcFldVal, dstFldVal)
	}
}

var timeType = reflect.TypeOf(time.Time{})

func (dr *DgraphRepository) copyValue(srcVal, dstVal reflect.Value) {
	if !srcVal.IsValid() {
		return
	}
	// resolve indirection
	srcValInd := reflectutil.Indirect(srcVal)

	switch srcValInd.Kind() {
	case reflect.Invalid:
		return

	case reflect.Struct:
		srcNode, ok := srcVal.Interface().(models.Node)
		if ok {
			// copy only ID
			id := srcNode.GetID()
			if id == nil {
				return
			}

			dstVal.Set(reflectutil.ZeroFromType(dstVal.Type()))
			dstVal = reflectutil.Indirect(dstVal)
			dstFlds := reflectutil.GetStructFields(dstVal)
			idFld, ok := dstFlds["ID"]
			if !ok {
				return
			}
			dstFld := dstVal.Field(idFld.Index)
			if dstFld.Kind() == reflect.Pointer {
				dstFld.Set(reflect.New(dstFld.Type().Elem()))
			}
			dstFld = reflect.Indirect(dstFld)
			dstFld.Set(reflect.ValueOf(*id))
			return
		}
		// not of type Node -> copy as is

		// everything else copy as is

		// src is a struct -> copy ID

		// // treat time.Time as a special case
		// if srcVal.Type() == reflect.TypeOf(time.Time{}) {
		// 	if dstVal.Kind() == reflect.Pointer && dstVal.IsNil() {
		// 		dstVal.Set(reflect.New(dstVal.Type().Elem()))
		// 	}
		// 	dstVal = reflect.Indirect(dstVal)
		// 	dstVal.Set(srcVal)
		// }

		// if dstVal.Kind() == reflect.Pointer {
		// 	dstVal.Set(reflect.New(dstVal.Type().Elem()))
		// 	dstVal = dstVal.Elem()
		// }
		// srcVal = reflectutil.Indirect(srcVal)
		// srcFlds := reflectutil.GetStructFields(srcVal)
		// idFld, ok := srcFlds["ID"]
		// if !ok {
		// 	return
		// }
		// srcFld := srcVal.Field(idFld.Index)
		// dstVal = reflectutil.Indirect(dstVal)
		// dstFlds := reflectutil.GetStructFields(dstVal)
		// idFld, ok = dstFlds["ID"]
		// if !ok {
		// 	return
		// }
		// dstFld := dstVal.Field(idFld.Index)
		// dr.copyValue(srcFld, dstFld)
		// return

	case reflect.Slice:
		if srcValInd.IsNil() {
			return
		}

		nslv := reflectutil.ZeroFromType(dstVal.Type()) // new slice value
		nslvd := reflectutil.Indirect(nslv)             // new slice value dereferenced
		nslvd.Set(reflect.MakeSlice(dstVal.Type(), 0, srcValInd.Len()))

		// check type of slice
		st := srcValInd.Type().Elem()
		if st.Implements(models.NodeType) {
			// dstValInd := reflectutil.Indirect(dstVal)
			dt := dstVal.Type().Elem()
			for i := 0; i < srcVal.Len(); i++ {
				srcNode := srcVal.Index(i).Interface().(models.Node)
				// copy only ID
				id := srcNode.GetID()
				if id == nil {
					continue
				}

				dstElm := reflectutil.ZeroFromType(dt)
				dstElmInd := reflectutil.Indirect(dstElm)
				dstFlds := reflectutil.GetStructFields(dstElmInd)
				idFld, ok := dstFlds["ID"]
				if !ok {
					continue
				}
				dstFld := dstElmInd.Field(idFld.Index)
				if dstFld.Kind() == reflect.Pointer {
					dstFld.Set(reflect.New(dstFld.Type().Elem()))
					dstFld = dstFld.Elem()
				}
				dstFld.Set(reflect.ValueOf(*id))
				nslvd.Set(reflect.Append(nslvd, dstElm))
			}

		} else {
			// not of type Node -> copy as is
			for i := 0; i < srcVal.Len(); i++ {
				nslvd.Set(reflect.Append(nslvd, srcVal.Index(i)))
			}
		}

		dstVal.Set(nslv)
		return

		// // check type of slice
		// st := srcValInd.Type().Elem()
		// if st.Implements(models.NodeType) {
		// 	// dstValInd := reflectutil.Indirect(dstVal)
		// 	dt := dstVal.Type().Elem()
		// 	if srcValInd.Len() > 0 {
		// 		nslv := reflectutil.ZeroFromType(dstVal.Type()) // new slice value
		// 		nslvd := reflectutil.Indirect(nslv)             // new slice value dereferenced
		// 		nslvd.Set(reflect.MakeSlice(dstVal.Type(), 0, srcValInd.Len()))

		// 		for i := 0; i < srcVal.Len(); i++ {
		// 			srcNode := srcVal.Index(i).Interface().(models.Node)
		// 			// copy only ID
		// 			id := srcNode.GetID()
		// 			if id == nil {
		// 				continue
		// 			}

		// 			dstElm := reflectutil.ZeroFromType(dt)
		// 			dstElmInd := reflectutil.Indirect(dstElm)
		// 			dstFlds := reflectutil.GetStructFields(dstElmInd)
		// 			idFld, ok := dstFlds["ID"]
		// 			if !ok {
		// 				continue
		// 			}
		// 			dstFld := dstElmInd.Field(idFld.Index)
		// 			if dstFld.Kind() == reflect.Pointer {
		// 				dstFld.Set(reflect.New(dstFld.Type().Elem()))
		// 				dstFld = dstFld.Elem()
		// 			}
		// 			dstFld.Set(reflect.ValueOf(*id))
		// 			nslvd.Set(reflect.Append(nslvd, dstElm))
		// 		}
		// 		dstVal.Set(nslv)
		// 	}
		// 	return
		// }
		// // not of type Node -> copy as is
	}

	// src is basically the same type as dst -> copy value directly
	if dstVal.Kind() == reflect.Pointer && dstVal.IsNil() {
		dstVal.Set(reflect.New(dstVal.Type().Elem()))
	}
	dstVal = reflect.Indirect(dstVal)
	dstVal.Set(srcValInd)
}
