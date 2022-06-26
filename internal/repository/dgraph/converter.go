package dgraph

import (
	"fmt"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

// initializeConverters initializes the converters for the copier.
func (dr *DgraphRepository) initializeConverters() {
	// // convertersForGet are used to convert DGraph output models to regular models.
	// dr.convertersForGet = []copier.TypeConverter{
	// 	// Category Converter
	// 	copier.TypeConverter{
	// 		SrcType: &dgclient.CategoryFragment_Parent{},
	// 		DstType: &models.Category{},
	// 		Fn: func(src interface{}) (interface{}, error) {
	// 			return &models.Category{ID: src.(*dgclient.CategoryFragment_Parent).ID}, nil
	// 		},
	// 	},

	// 	// Host Converter
	// 	copier.TypeConverter{
	// 		SrcType: &dgclient.UserFragment_UserOrGroupFragment_Host{},
	// 		DstType: &models.Host{},
	// 		Fn: func(src interface{}) (interface{}, error) {
	// 			return &models.Host{ID: src.(*dgclient.UserFragment_UserOrGroupFragment_Host).ID}, nil
	// 		},
	// 	},

	// 	// Avatar Converter
	// 	copier.TypeConverter{
	// 		SrcType: &dgclient.UserFragment_UserOrGroupFragment_Avatar{},
	// 		DstType: &models.File{},
	// 		Fn: func(src interface{}) (interface{}, error) {
	// 			return &models.File{ID: src.(*dgclient.UserFragment_UserOrGroupFragment_Avatar).ID}, nil
	// 		},
	// 	},

	// 	// Product Converter
	// 	copier.TypeConverter{
	// 		SrcType: []*dgclient.UserFragment_UserOrGroupFragment_Products{},
	// 		DstType: []*models.Product{},
	// 		Fn: func(src interface{}) (interface{}, error) {
	// 			srcList := src.([]*dgclient.UserFragment_UserOrGroupFragment_Products)
	// 			dstList := make([]*models.Product, 0, len(srcList))
	// 			for _, x := range srcList {
	// 				dstList = append(dstList, &models.Product{ID: x.ID})
	// 			}
	// 			return dstList, nil
	// 		},
	// 	},

	// 	// MemberOf List
	// 	copier.TypeConverter{
	// 		SrcType: []*dgclient.UserFragment_UserOrGroupFragment_MemberOf{},
	// 		DstType: []*models.Group{},
	// 		Fn: func(src interface{}) (interface{}, error) {
	// 			srcList := src.([]*dgclient.UserFragment_UserOrGroupFragment_MemberOf)
	// 			dstList := make([]*models.Group, 0, len(srcList))
	// 			for _, x := range srcList {
	// 				dstList = append(dstList, &models.Group{ID: x.ID})
	// 			}
	// 			return dstList, nil
	// 		},
	// 	},

	// 	// Members List Converter
	// 	copier.TypeConverter{
	// 		SrcType: []*dgclient.GroupFragment_Members{},
	// 		DstType: []models.UserOrGroup{},
	// 		Fn: func(src interface{}) (interface{}, error) {
	// 			srcList := src.([]*dgclient.GroupFragment_Members)
	// 			dstList := make([]models.UserOrGroup, 0, len(srcList))
	// 			for _, x := range srcList {
	// 				if x.User.ID != "" {
	// 					dstList = append(dstList, &models.User{ID: x.User.ID})
	// 				} else if x.Group.ID != "" {
	// 					dstList = append(dstList, &models.Group{ID: x.Group.ID})
	// 				} else {
	// 					panic(fmt.Errorf("implementation error: unknown UserOrGroup type: %T", x))
	// 				}
	// 			}
	// 			return dstList, nil
	// 		},
	// 	},
	// }

	// -------------------------------------------------------------------------

	var userOrGroup models.UserOrGroup

	// convertersForSave are used to convert regular models to Dgraph input models.
	dr.convertersForSave = []copier.TypeConverter{
		// Category Converter
		copier.TypeConverter{
			SrcType: &models.Category{},
			DstType: &models.CategoryRef{},
			Fn: func(src interface{}) (interface{}, error) {
				return dr.saveCategoryIfNecessary(src.(*models.Category))
			},
		},

		// Component Converter
		copier.TypeConverter{
			SrcType: &models.Component{},
			DstType: &models.ComponentRef{},
			Fn: func(src interface{}) (interface{}, error) {
				return dr.saveComponentIfNecessary(src.(*models.Component))
			},
		},
		copier.TypeConverter{
			SrcType: []*models.Component{},
			DstType: []*models.ComponentRef{},
			Fn: func(src interface{}) (interface{}, error) {
				srcList := src.([]*models.Component)
				dstList := make([]*models.ComponentRef, 0, len(srcList))
				for _, x := range srcList {
					ref, err := dr.saveComponentIfNecessary(x)
					if err != nil {
						return nil, repository.NewRepoErrorWrap(err, errSaveComponentStr)
					}
					dstList = append(dstList, ref)
				}
				return dstList, nil
			},
		},

		// ComponentSource Converter
		copier.TypeConverter{
			SrcType: &models.ComponentSource{},
			DstType: &models.ComponentSourceRef{},
			Fn: func(src interface{}) (interface{}, error) {
				return dr.saveComponentSourceIfNecessary(src.(*models.ComponentSource))
			},
		},

		// File Converter
		copier.TypeConverter{
			SrcType: &models.File{},
			DstType: &models.FileRef{},
			Fn: func(src interface{}) (interface{}, error) {
				return dr.saveFileIfNecessary(src.(*models.File))
			},
		},

		// Host Converter
		copier.TypeConverter{
			SrcType: &models.Host{},
			DstType: &models.HostRef{},
			Fn: func(src interface{}) (interface{}, error) {
				return dr.saveHostIfNecessary(src.(*models.Host))
			},
		},

		// Product Converter
		copier.TypeConverter{
			SrcType: &models.Product{},
			DstType: &models.ProductRef{},
			Fn: func(src interface{}) (interface{}, error) {
				return dr.saveProductIfNecessary(src.(*models.Product))
			},
		},
		copier.TypeConverter{
			SrcType: []*models.Product{},
			DstType: []*models.ProductRef{},
			Fn: func(src interface{}) (interface{}, error) {
				srcList := src.([]*models.Product)
				dstList := make([]*models.ProductRef, 0, len(srcList))
				for _, x := range srcList {
					ref, err := dr.saveProductIfNecessary(x)
					if err != nil {
						return nil, repository.NewRepoErrorWrap(err, errSaveProductStr)
					}
					dstList = append(dstList, ref)
				}
				return dstList, nil
			},
		},

		// User Converter
		copier.TypeConverter{
			SrcType: &models.User{},
			DstType: &models.UserRef{},
			Fn: func(src interface{}) (interface{}, error) {
				return dr.saveUserIfNecessary(src.(*models.User))
			},
		},
		copier.TypeConverter{
			SrcType: []*models.User{},
			DstType: []*models.UserRef{},
			Fn: func(src interface{}) (interface{}, error) {
				srcList := src.([]*models.User)
				dstList := make([]*models.UserRef, 0, len(srcList))
				for _, x := range srcList {
					ref, err := dr.saveUserIfNecessary(x)
					if err != nil {
						return nil, repository.NewRepoErrorWrap(err, errSaveUserStr).
							AddIfNotNil("userId", x.ID).AddIfNotNil("userXid", x.Xid)
					}
					dstList = append(dstList, ref)
				}
				return dstList, nil
			},
		},

		// Group Converter
		copier.TypeConverter{
			SrcType: &models.Group{},
			DstType: &models.GroupRef{},
			Fn: func(src interface{}) (interface{}, error) {
				return dr.saveGroupIfNecessary(src.(*models.Group))
			},
		},
		copier.TypeConverter{
			SrcType: []*models.Group{},
			DstType: []*models.GroupRef{},
			Fn: func(src interface{}) (interface{}, error) {
				srcList := src.([]*models.Group)
				dstList := make([]*models.GroupRef, 0, len(srcList))
				for _, x := range srcList {
					ref, err := dr.saveGroupIfNecessary(x)
					if err != nil {
						return nil, repository.NewRepoErrorWrap(err, errSaveGroupStr).
							AddIfNotNil("groupId", x.ID).AddIfNotNil("groupXid", x.Xid)
					}
					dstList = append(dstList, ref)
				}
				return dstList, nil
			},
		},

		// UserOrGroup Converter
		copier.TypeConverter{
			SrcType: userOrGroup,
			DstType: &models.UserOrGroupRef{},
			Fn: func(src interface{}) (interface{}, error) {
				return dr.saveUserOrGroupIfNecessary(src.(models.UserOrGroup))
			},
		},
		copier.TypeConverter{
			SrcType: []models.UserOrGroup{},
			DstType: []*models.UserOrGroupRef{},
			Fn: func(src interface{}) (interface{}, error) {
				srcList := src.([]models.UserOrGroup)
				dstList := make([]*models.UserOrGroupRef, 0, len(srcList))
				for _, x := range srcList {
					ref, err := dr.saveUserOrGroupIfNecessary(x)
					if err != nil {
						switch x.(type) {
						case *models.User:
							return nil, repository.NewRepoErrorWrap(err, errSaveUserStr)
						case *models.Group:
							return nil, repository.NewRepoErrorWrap(err, errSaveGroupStr)
						}
						panic(fmt.Errorf("implementation error: unknown UserOrGroup type: %T", x))
					}
					dstList = append(dstList, ref)
				}
				return dstList, nil
			},
		},
	}
}
