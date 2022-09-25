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
	"encoding/json"
	"strings"
)

// NodeDumpJSON marshals a Product to JSON.
func NodeDumpJSON(node Node, indent ...uint8) ([]byte, error) {
	// remove potential circular reference
	node = AsTree(node)
	if len(indent) > 0 {
		return json.MarshalIndent(node, "", strings.Repeat(" ", int(indent[0])))
	}
	return json.Marshal(node)
}

// ProductDumpJSON marshals a Product to JSON.
func ProductDumpJSON(prd *Product, indent ...uint8) ([]byte, error) {
	// remove potential circular reference
	p := AsTree(prd)
	if len(indent) > 0 {
		return json.MarshalIndent(p, "", strings.Repeat(" ", int(indent[0])))
	}
	return json.Marshal(p)
}

// ProductLoadJSON unmarshals a Product from JSON.
func ProductLoadJSON(data []byte) (*Product, error) {
	prd := &Product{}
	err := json.Unmarshal(data, prd)
	if err != nil {
		return nil, err
	}
	// restore circular reference
	prd = AsGraph(prd).(*Product)
	return prd, nil
}
