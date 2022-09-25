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

package services

import (
	"context"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
)

// LicenseProvider is the interface for providing SPDX licenses.
type LicenseProvider interface {
	// GetLicense returns a license by its SPDX ID.
	GetLicense(ctx context.Context, id, xid *string) (*models.License, error)
	GetLicenseID(ctx context.Context, xid *string) (*string, error)

	// GetAllLicenses returns a list of all licenses.
	GetAllLicenses(ctx context.Context) ([]*models.License, error)
}

// Repository is the interface for getting and saving licenses to a repository.
type Repository interface {
	NodeRepository
	ProductRepository
	ComponentRepository
	SoftwareRepository
	RepositoryRepository
	TechnologySpecificDocumentationCriteriaRepository
	TechnicalStandardRepository
	UserRepository
	GroupRepository
	FileRepository
	KeyValueRepository
	StringVRepository
	FloatVRepository
	MaterialRepository
	ManufacturingProcessRepository
	BoundingBoxDimensionsRepository
	OpenSCADDimensionsRepository
	CategoryRepository
	TagRepository
	LicenseRepository
	HostRepository
}

// NodeRepository is an interface for getting `Node` objects from a repository.
type NodeRepository interface {
	GetNode(ctx context.Context, id string) (interface{}, error)
}

// ProductRepository is an interface for getting and saving `Product` objects to a repository.
type ProductRepository interface {
	GetProduct(ctx context.Context, id, xid *string) (*models.Product, error)
	GetProductID(ctx context.Context, xid *string) (*string, error)
	GetProducts(ctx context.Context, filter *dgclient.ProductFilter, order *dgclient.ProductOrder, first *int64, offset *int64) ([]*models.Product, int64, error)
	GetAllProducts(ctx context.Context) ([]*models.Product, int64, error)
	CreateProduct(ctx context.Context, input *models.Product) error
	UpdateProduct(ctx context.Context, input *models.Product) error
	DeleteProduct(ctx context.Context, id, xid *string) error
	DeleteAllProducts(ctx context.Context) error
}

// ComponentRepository is an interface for getting and saving `Component` objects to a repository.
type ComponentRepository interface {
	GetComponent(ctx context.Context, id, xid *string) (*models.Component, error)
	GetComponentID(ctx context.Context, xid *string) (*string, error)
	GetComponents(ctx context.Context, filter *dgclient.ComponentFilter, order *dgclient.ComponentOrder, first *int64, offset *int64) ([]*models.Component, int64, error)
	GetAllComponents(ctx context.Context) ([]*models.Component, int64, error)
	CreateComponent(ctx context.Context, input *models.Component) error
	UpdateComponent(ctx context.Context, input *models.Component) error
	DeleteComponent(ctx context.Context, id, xid *string) error
	DeleteAllComponents(ctx context.Context) error
}

// SoftwareRepository is an interface for getting and saving `Software` objects to a repository.
type SoftwareRepository interface {
	GetSoftware(ctx context.Context, id *string) (*models.Software, error)
	GetSoftwares(ctx context.Context, filter *dgclient.SoftwareFilter, order *dgclient.SoftwareOrder, first *int64, offset *int64) ([]*models.Software, int64, error)
	GetAllSoftwares(ctx context.Context) ([]*models.Software, int64, error)
	CreateSoftware(ctx context.Context, input *models.Software) error
	UpdateSoftware(ctx context.Context, input *models.Software) error
	DeleteSoftware(ctx context.Context, id *string) error
	DeleteAllSoftwares(ctx context.Context) error
}

// RepositoryRepository is an interface for getting and saving `Repository` objects to a repository.
type RepositoryRepository interface {
	GetRepository(ctx context.Context, id, xid *string) (*models.Repository, error)
	GetRepositoryID(ctx context.Context, xid *string) (*string, error)
	GetRepositories(ctx context.Context, filter *dgclient.RepositoryFilter, order *dgclient.RepositoryOrder, first *int64, offset *int64) ([]*models.Repository, int64, error)
	GetAllRepositories(ctx context.Context) ([]*models.Repository, int64, error)
	CreateRepository(ctx context.Context, input *models.Repository) error
	UpdateRepository(ctx context.Context, input *models.Repository) error
	DeleteRepository(ctx context.Context, id, xid *string) error
	DeleteAllRepositories(ctx context.Context) error
}

// TechnologySpecificDocumentationCriteriaRepository is an interface for getting and saving `TechnologySpecificDocumentationCriteria` objects to a repository.
type TechnologySpecificDocumentationCriteriaRepository interface {
	GetTechnologySpecificDocumentationCriteria(ctx context.Context, id, xid *string) (*models.TechnologySpecificDocumentationCriteria, error)
	GetTechnologySpecificDocumentationCriteriaID(ctx context.Context, xid *string) (*string, error)
	GetTechnologySpecificDocumentationCriterias(ctx context.Context, filter *dgclient.TechnologySpecificDocumentationCriteriaFilter, order *dgclient.TechnologySpecificDocumentationCriteriaOrder, first *int64, offset *int64) ([]*models.TechnologySpecificDocumentationCriteria, int64, error)
	GetAllTechnologySpecificDocumentationCriterias(ctx context.Context) ([]*models.TechnologySpecificDocumentationCriteria, int64, error)
	CreateTechnologySpecificDocumentationCriteria(ctx context.Context, input *models.TechnologySpecificDocumentationCriteria) error
	UpdateTechnologySpecificDocumentationCriteria(ctx context.Context, input *models.TechnologySpecificDocumentationCriteria) error
	DeleteTechnologySpecificDocumentationCriteria(ctx context.Context, id, xid *string) error
	DeleteAllTechnologySpecificDocumentationCriterias(ctx context.Context) error
}

// TechnicalStandardRepository is an interface for getting and saving `TechnicalStandard` objects to a repository.
type TechnicalStandardRepository interface {
	GetTechnicalStandard(ctx context.Context, id, xid *string) (*models.TechnicalStandard, error)
	GetTechnicalStandardID(ctx context.Context, xid *string) (*string, error)
	GetTechnicalStandards(ctx context.Context, filter *dgclient.TechnicalStandardFilter, order *dgclient.TechnicalStandardOrder, first *int64, offset *int64) ([]*models.TechnicalStandard, int64, error)
	GetAllTechnicalStandards(ctx context.Context) ([]*models.TechnicalStandard, int64, error)
	CreateTechnicalStandard(ctx context.Context, input *models.TechnicalStandard) error
	UpdateTechnicalStandard(ctx context.Context, input *models.TechnicalStandard) error
	DeleteTechnicalStandard(ctx context.Context, id, xid *string) error
	DeleteAllTechnicalStandards(ctx context.Context) error
}

// UserRepository is an interface for getting and saving `User` objects to a repository.
type UserRepository interface {
	GetUser(ctx context.Context, id, xid *string) (*models.User, error)
	GetUserID(ctx context.Context, xid *string) (*string, error)
	GetUsers(ctx context.Context, filter *dgclient.UserFilter, order *dgclient.UserOrder, first *int64, offset *int64) ([]*models.User, int64, error)
	GetAllUsers(ctx context.Context) ([]*models.User, int64, error)
	CreateUser(ctx context.Context, input *models.User) error
	UpdateUser(ctx context.Context, input *models.User) error
	DeleteUser(ctx context.Context, id, xid *string) error
	DeleteAllUsers(ctx context.Context) error
}

// GroupRepository is an interface for getting and saving `Group` objects to a repository.
type GroupRepository interface {
	GetGroup(ctx context.Context, id, xid *string) (*models.Group, error)
	GetGroupID(ctx context.Context, xid *string) (*string, error)
	GetGroups(ctx context.Context, filter *dgclient.GroupFilter, order *dgclient.GroupOrder, first *int64, offset *int64) ([]*models.Group, int64, error)
	GetAllGroups(ctx context.Context) ([]*models.Group, int64, error)
	CreateGroup(ctx context.Context, input *models.Group) error
	UpdateGroup(ctx context.Context, input *models.Group) error
	DeleteGroup(ctx context.Context, id, xid *string) error
	DeleteAllGroups(ctx context.Context) error
}

// FileRepository is an interface for getting and saving `File` objects to a repository.
type FileRepository interface {
	GetFile(ctx context.Context, id, xid *string) (*models.File, error)
	GetFileID(ctx context.Context, xid *string) (*string, error)
	GetFiles(ctx context.Context, filter *dgclient.FileFilter, order *dgclient.FileOrder, first *int64, offset *int64) ([]*models.File, int64, error)
	GetAllFiles(ctx context.Context) ([]*models.File, int64, error)
	CreateFile(ctx context.Context, input *models.File) error
	UpdateFile(ctx context.Context, input *models.File) error
	DeleteFile(ctx context.Context, id, xid *string) error
	DeleteAllFiles(ctx context.Context) error
}

// KeyValueRepository is an interface for getting and saving `KeyValue` objects to a repository.
type KeyValueRepository interface {
	GetKeyValue(ctx context.Context, id *string) (*models.KeyValue, error)
	GetKeyValues(ctx context.Context, filter *dgclient.KeyValueFilter, order *dgclient.KeyValueOrder, first *int64, offset *int64) ([]*models.KeyValue, int64, error)
	GetAllKeyValues(ctx context.Context) ([]*models.KeyValue, int64, error)
	CreateKeyValue(ctx context.Context, input *models.KeyValue) error
	UpdateKeyValue(ctx context.Context, input *models.KeyValue) error
	DeleteKeyValue(ctx context.Context, id *string) error
	DeleteAllKeyValues(ctx context.Context) error
}

// FloatVRepository is an interface for getting and saving `FloatV` objects to a repository.
type FloatVRepository interface {
	GetFloatV(ctx context.Context, id *string) (*models.FloatV, error)
	GetFloatVs(ctx context.Context, filter *dgclient.FloatVFilter, order *dgclient.FloatVOrder, first *int64, offset *int64) ([]*models.FloatV, int64, error)
	GetAllFloatVs(ctx context.Context) ([]*models.FloatV, int64, error)
	CreateFloatV(ctx context.Context, input *models.FloatV) error
	UpdateFloatV(ctx context.Context, input *models.FloatV) error
	DeleteFloatV(ctx context.Context, id *string) error
	DeleteAllFloatVs(ctx context.Context) error
}

// StringVRepository is an interface for getting and saving `StringV` objects to a repository.
type StringVRepository interface {
	GetStringV(ctx context.Context, id *string) (*models.StringV, error)
	GetStringVs(ctx context.Context, filter *dgclient.StringVFilter, order *dgclient.StringVOrder, first *int64, offset *int64) ([]*models.StringV, int64, error)
	GetAllStringVs(ctx context.Context) ([]*models.StringV, int64, error)
	CreateStringV(ctx context.Context, input *models.StringV) error
	UpdateStringV(ctx context.Context, input *models.StringV) error
	DeleteStringV(ctx context.Context, id *string) error
	DeleteAllStringVs(ctx context.Context) error
}

// MaterialRepository is an interface for getting and saving `Material` objects to a repository.
type MaterialRepository interface {
	GetMaterial(ctx context.Context, id *string) (*models.Material, error)
	GetMaterials(ctx context.Context, filter *dgclient.MaterialFilter, order *dgclient.MaterialOrder, first *int64, offset *int64) ([]*models.Material, int64, error)
	GetAllMaterials(ctx context.Context) ([]*models.Material, int64, error)
	CreateMaterial(ctx context.Context, input *models.Material) error
	UpdateMaterial(ctx context.Context, input *models.Material) error
	DeleteMaterial(ctx context.Context, id *string) error
	DeleteAllMaterials(ctx context.Context) error
}

// ManufacturingProcessRepository is an interface for getting and saving `ManufacturingProcess` objects to a repository.
type ManufacturingProcessRepository interface {
	GetManufacturingProcess(ctx context.Context, id *string) (*models.ManufacturingProcess, error)
	GetManufacturingProcesses(ctx context.Context, filter *dgclient.ManufacturingProcessFilter, order *dgclient.ManufacturingProcessOrder, first *int64, offset *int64) ([]*models.ManufacturingProcess, int64, error)
	GetAllManufacturingProcesses(ctx context.Context) ([]*models.ManufacturingProcess, int64, error)
	CreateManufacturingProcess(ctx context.Context, input *models.ManufacturingProcess) error
	UpdateManufacturingProcess(ctx context.Context, input *models.ManufacturingProcess) error
	DeleteManufacturingProcess(ctx context.Context, id *string) error
	DeleteAllManufacturingProcesses(ctx context.Context) error
}

// BoundingBoxDimensionsRepository is an interface for getting and saving `BoundingBoxDimensions` objects to a repository.
type BoundingBoxDimensionsRepository interface {
	GetBoundingBoxDimensions(ctx context.Context, id *string) (*models.BoundingBoxDimensions, error)
	GetBoundingBoxDimensionss(ctx context.Context, filter *dgclient.BoundingBoxDimensionsFilter, order *dgclient.BoundingBoxDimensionsOrder, first *int64, offset *int64) ([]*models.BoundingBoxDimensions, int64, error)
	GetAllBoundingBoxDimensionss(ctx context.Context) ([]*models.BoundingBoxDimensions, int64, error)
	CreateBoundingBoxDimensions(ctx context.Context, input *models.BoundingBoxDimensions) error
	UpdateBoundingBoxDimensions(ctx context.Context, input *models.BoundingBoxDimensions) error
	DeleteBoundingBoxDimensions(ctx context.Context, id *string) error
	DeleteAllBoundingBoxDimensionss(ctx context.Context) error
}

// OpenSCADDimensionsRepository is an interface for getting and saving `OpenSCADDimensions` objects to a repository.
type OpenSCADDimensionsRepository interface {
	GetOpenSCADDimensions(ctx context.Context, id *string) (*models.OpenSCADDimensions, error)
	GetOpenSCADDimensionss(ctx context.Context, filter *dgclient.OpenSCADDimensionsFilter, order *dgclient.OpenSCADDimensionsOrder, first *int64, offset *int64) ([]*models.OpenSCADDimensions, int64, error)
	GetAllOpenSCADDimensionss(ctx context.Context) ([]*models.OpenSCADDimensions, int64, error)
	CreateOpenSCADDimensions(ctx context.Context, input *models.OpenSCADDimensions) error
	UpdateOpenSCADDimensions(ctx context.Context, input *models.OpenSCADDimensions) error
	DeleteOpenSCADDimensions(ctx context.Context, id *string) error
	DeleteAllOpenSCADDimensionss(ctx context.Context) error
}

// CategoryRepository is an interface for getting and saving `Category` objects to a repository.
type CategoryRepository interface {
	GetCategory(ctx context.Context, id, xid *string) (*models.Category, error)
	GetCategoryID(ctx context.Context, xid *string) (*string, error)
	GetCategories(ctx context.Context, filter *dgclient.CategoryFilter, order *dgclient.CategoryOrder, first *int64, offset *int64) ([]*models.Category, int64, error)
	GetAllCategories(ctx context.Context) ([]*models.Category, int64, error)
	CreateCategory(ctx context.Context, input *models.Category) error
	UpdateCategory(ctx context.Context, input *models.Category) error
	DeleteCategory(ctx context.Context, id, xid *string) error
	DeleteAllCategories(ctx context.Context) error
}

// TagRepository is an interface for getting and saving `Tag` objects to a repository.
type TagRepository interface {
	GetTag(ctx context.Context, id, xid *string) (*models.Tag, error)
	GetTagID(ctx context.Context, name *string) (*string, error)
	GetTags(ctx context.Context, filter *dgclient.TagFilter, order *dgclient.TagOrder, first *int64, offset *int64) ([]*models.Tag, int64, error)
	GetAllTags(ctx context.Context) ([]*models.Tag, int64, error)
	CreateTag(ctx context.Context, input *models.Tag) error
	UpdateTag(ctx context.Context, input *models.Tag) error
	DeleteTag(ctx context.Context, id, xid *string) error
	DeleteAllTags(ctx context.Context) error
}

// LicenseRepository is an interface for getting and saving `License` objects to a repository.
type LicenseRepository interface {
	GetLicense(ctx context.Context, id, xid *string) (*models.License, error)
	GetLicenseID(ctx context.Context, xid *string) (*string, error)
	GetLicenses(ctx context.Context, filter *dgclient.LicenseFilter, order *dgclient.LicenseOrder, first *int64, offset *int64) ([]*models.License, int64, error)
	GetAllLicenses(ctx context.Context) ([]*models.License, int64, error)
	CreateLicense(ctx context.Context, input *models.License) error
	UpdateLicense(ctx context.Context, input *models.License) error
	DeleteLicense(ctx context.Context, id, xid *string) error
	DeleteAllLicenses(ctx context.Context) error
}

// HostRepository is an interface for getting and saving `Host` objects to a repository.
type HostRepository interface {
	GetHost(ctx context.Context, id, domain *string) (*models.Host, error)
	GetHostID(ctx context.Context, domain *string) (*string, error)
	GetHosts(ctx context.Context, filter *dgclient.HostFilter, order *dgclient.HostOrder, first *int64, offset *int64) ([]*models.Host, int64, error)
	GetAllHosts(ctx context.Context) ([]*models.Host, int64, error)
	GetAllLicensesBasic(ctx context.Context) ([]*models.License, error)
	CreateHost(ctx context.Context, input *models.Host) error
	UpdateHost(ctx context.Context, input *models.Host) error
	DeleteHost(ctx context.Context, id, domain *string) error
	DeleteAllHosts(ctx context.Context) error
}
