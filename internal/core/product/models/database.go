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

type Database struct {
	ID      *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Version *string `mandatory:"true" json:"version,omitempty" graphql:"version" dql:"Database.version"`
}

func (*Database) IsNode() {}
