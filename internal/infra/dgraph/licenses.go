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

	"losh/internal/core/product/models"
)

// GetAllLicensesBasic returns a list of all `License` objects with basic
// information.
func (dr *DgraphRepository) GetAllLicensesBasic(ctx context.Context) ([]*models.License, error) {
	rsp := struct {
		Licenses []*models.License "json:\"queryLicense\" graphql:\"queryLicense\""
	}{}
	if err := dr.client.GetAllLicensesBasicWithResponse(ctx, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetLicenseStr)
	}
	return rsp.Licenses, nil
}
