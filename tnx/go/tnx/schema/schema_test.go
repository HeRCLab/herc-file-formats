package schema

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/kr/pretty"
)

func compareJSON(s1, s2 string) (bool, error) {
	var i1 interface{}
	var i2 interface{}

	if err := json.Unmarshal([]byte(s1), &i1); err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(s2), &i2); err != nil {
		return false, err
	}

	return cmp.Equal(i1, i2, cmpopts.IgnoreUnexported(TNX{})), nil
}

func TestFromJSONTopology(t *testing.T) {

	text := `
{
	"schema": ["tnx", "0"],
	"topology": {
		"nodes": [
			{
				"id": "foo",
				"operation": "input",
				"outputs": ["foo->output0"]
			},
			{
				"id": "bar",
				"operation": "output",
				"inputs": ["bar<-input0"]
			}
		],
		"links": [
			{ "source": "foo->output0", "target": "bar<-input0"}
		]
	}
}
`

	tnx, err := FromJSON([]byte(text))
	if err != nil {
		t.Error(err)
	}

	expect := &TNX{
		Schema: []string{"tnx", "0"},
		Topology: Topology{
			Nodes: []Node{
				Node{
					ID:        "foo",
					Operation: "input",
					Inputs:    nil,
					Outputs:   []string{"foo->output0"},
				},
				Node{
					ID:        "bar",
					Operation: "output",
					Inputs:    []string{"bar<-input0"},
					Outputs:   nil,
				},
			},
			Links: []Link{
				Link{
					Source: "foo->output0",
					Target: "bar<-input0",
				},
			},
		},
	}

	t.Logf(pretty.Sprintf("Actual decoded TNX: %#v\n", tnx))
	t.Logf(pretty.Sprintf("Expected: %#v\n", expect))
	if !cmp.Equal(tnx, expect, cmpopts.IgnoreUnexported(TNX{})) {
		t.Errorf("Actual and expected TNX values differ")
		t.Logf("\n\nDifferences: \n\n")
		for _, v := range pretty.Diff(tnx, expect) {
			t.Logf(v)
		}
	}
}

func TestFromJSONParameters(t *testing.T) {

	text := `
{
	"schema": ["tnx", "0"],
	"topology": {
		"nodes": [
			{
				"id": "foo",
				"operation": "input",
				"outputs": ["foo->output0"]
			},
			{
				"id": "bar",
				"operation": "output",
				"inputs": ["bar<-input0"]
			}
		],
		"links": [
			{ "source": "foo->output0", "target": "bar<-input0"}
		]
	},
	"parameters": {
		"foo": {
			"dimensions": [10, 10],
			"deltas": [1.0, 2.0, 3.5],
			"weights": [1.0, 1.5],
			"biases": [0.5, 1.0],
			"activation": "test string"
		}
	}
}
`

	tnx, err := FromJSON([]byte(text))
	if err != nil {
		t.Error(err)
	}

	expect := &TNX{
		Schema: []string{"tnx", "0"},
		Topology: Topology{
			Nodes: []Node{
				Node{
					ID:        "foo",
					Operation: "input",
					Inputs:    nil,
					Outputs:   []string{"foo->output0"},
				},
				Node{
					ID:        "bar",
					Operation: "output",
					Inputs:    []string{"bar<-input0"},
					Outputs:   nil,
				},
			},
			Links: []Link{
				Link{
					Source: "foo->output0",
					Target: "bar<-input0",
				},
			},
		},
		Parameters: map[string]*Parameter{
			"foo": &Parameter{
				Dimensions: &[]int{10, 10},
				Deltas:     &[]float64{1.0, 2.0, 3.5},
				Weights:    &[]float64{1.0, 1.5},
				Biases:     &[]float64{0.5, 1.0},
			},
		},
	}

	s1 := "test string"
	expect.Parameters["foo"].Activation = &s1

	t.Logf(pretty.Sprintf("Actual decoded TNX: %#v\n", tnx))
	t.Logf(pretty.Sprintf("Expected: %#v\n", expect))
	if !cmp.Equal(tnx, expect, cmpopts.IgnoreUnexported(TNX{})) {
		t.Errorf("Actual and expected TNX values differ")
		t.Logf("\n\nDifferences: \n\n")
		for _, v := range pretty.Diff(tnx, expect) {
			t.Logf(v)
		}
	}
}

func TestFromJSONSnapshot(t *testing.T) {

	text := `
{
	"schema": ["tnx", "0"],
	"topology": {
		"nodes": [
			{
				"id": "foo",
				"operation": "input",
				"outputs": ["foo->output0"]
			},
			{
				"id": "bar",
				"operation": "output",
				"inputs": ["bar<-input0"]
			}
		],
		"links": [
			{ "source": "foo->output0", "target": "bar<-input0"}
		]
	},
	"parameters": {
		"foo": {
			"dimensions": [10, 10],
			"deltas": [1.0, 2.0, 3.5],
			"weights": [1.0, 1.5],
			"biases": [0.5, 1.0],
			"activation": "test string"
		}
	},
	"snapshots": {
		"foo->output0": {
			"matrix": {
				"output": {
					"name": "",
					"dimensions": [10, 10],
					"data": [1.0, 2.0, 3.5]
				}
			}
		}
	}
}
`

	tnx, err := FromJSON([]byte(text))
	if err != nil {
		t.Error(err)
	}

	expect := &TNX{
		Schema: []string{"tnx", "0"},
		Topology: Topology{
			Nodes: []Node{
				Node{
					ID:        "foo",
					Operation: "input",
					Inputs:    nil,
					Outputs:   []string{"foo->output0"},
				},
				Node{
					ID:        "bar",
					Operation: "output",
					Inputs:    []string{"bar<-input0"},
					Outputs:   nil,
				},
			},
			Links: []Link{
				Link{
					Source: "foo->output0",
					Target: "bar<-input0",
				},
			},
		},
		Parameters: map[string]*Parameter{
			"foo": &Parameter{
				Dimensions: &[]int{10, 10},
				Deltas:     &[]float64{1.0, 2.0, 3.5},
				Weights:    &[]float64{1.0, 1.5},
				Biases:     &[]float64{0.5, 1.0},
			},
		},
		Snapshots: map[string]*Snapshot{
			"foo->output0": &Snapshot{
				Matrix: map[string]*Matrix{
					"output": &Matrix{
						Dimensions: []int{10, 10},
						Name:       "",
						Data:       []float64{1.0, 2.0, 3.5},
					},
				},
			},
		},
	}

	s1 := "test string"
	expect.Parameters["foo"].Activation = &s1

	t.Logf(pretty.Sprintf("Actual decoded TNX: %#v\n", tnx))
	t.Logf(pretty.Sprintf("Expected: %#v\n", expect))
	if !cmp.Equal(tnx, expect, cmpopts.IgnoreUnexported(TNX{})) {
		t.Errorf("Actual and expected TNX values differ")
		t.Logf("\n\nDifferences: \n\n")
		for _, v := range pretty.Diff(tnx, expect) {
			t.Logf(v)
		}
	}
}
