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

// ProductID serves as an identifier for projects, that can be used by the appropriate fetcher to fetch the projects metadata.
type ProductID struct {
	Platform string `json:"platform"`
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	Path     string `json:"path"`
}

// NewProductID creates a new ProductID from the given platform, owner, repo and path.
func NewProductID(platform, owner, repo, path string) ProductID {
	return ProductID{
		Platform: platform,
		Owner:    owner,
		Repo:     repo,
		Path:     path,
	}
}

// NewProductIDFromURL creates a new ProductID from the given URL.
func NewProductIDFromURL(url string) (ProductID, error) {
	pu, err := NewProductURLFromURL(url)
	if err != nil {
		return ProductID{}, err
	}
	return NewProductID(pu.Domain, pu.Owner, pu.Repo, pu.Path), nil
}

// String returns the string representation of the ProductID.
func (p ProductID) String() string {
	result := p.Platform
	elms := []string{p.Owner, p.Repo, p.Path}
	for _, elm := range elms {
		if elm != "" {
			result += "/" + elm
		}
	}
	return result
}
