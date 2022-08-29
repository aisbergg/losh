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
