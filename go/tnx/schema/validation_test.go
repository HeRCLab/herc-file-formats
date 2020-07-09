package schema

import (
	"testing"

	"github.com/herclab/tnx/go/tnx/schema/samples"

	"github.com/google/go-cmp/cmp"
)

func TestValidateSchema(t *testing.T) {
	cases := []struct {
		input []string
		err   bool
	}{
		{[]string{"tnx", "0"}, false},
		{[]string{"tnx", "1"}, true},
		{[]string{"foo", "1"}, true},
		{[]string{"foo", "0"}, true},
		{[]string{"tnx"}, true},
		{[]string{"tnx", "0", "0"}, true},
	}

	for i, v := range cases {
		t.Logf("Test case %d: input=%v, err=%v", i, v.input, v.err)
		res := ValidateSchema(v.input)
		if v.err && res == nil {
			t.Errorf("Schema should have errored but did not")
		}
		if !v.err && res != nil {
			t.Errorf("Schema should not have errored but did")
		}
	}
}

func TestValidateTopology(t *testing.T) {
	cases := []struct {
		input Topology
		err   bool
	}{
		{
			// Most trivial case where there is just one node
			input: Topology{
				Nodes: []Node{
					Node{
						ID:        "foo",
						Operation: "input",
						Inputs:    []string{},
						Outputs:   []string{},
					},
				},
			},
			err: false,
		},
		{
			// Slightly less trivial case, two nodes and one link
			input: Topology{
				Nodes: []Node{
					Node{
						ID:        "foo",
						Operation: "input",
						Inputs:    []string{},
						Outputs:   []string{"foo->output0"},
					},
					Node{
						ID:        "bar",
						Operation: "output",
						Inputs:    []string{"bar->input0"},
						Outputs:   []string{},
					},
				},
				Links: []Link{
					Link{
						Source: "foo->output0",
						Target: "bar->input0",
					},
				},
			},
			err: false,
		},
		{
			// Two nodes with aliased IDs
			input: Topology{
				Nodes: []Node{
					Node{
						ID:        "foo",
						Operation: "input",
						Inputs:    []string{},
						Outputs:   []string{"foo->output0"},
					},
					Node{
						ID:        "foo",
						Operation: "output",
						Inputs:    []string{"bar->input0"},
						Outputs:   []string{},
					},
				},
				Links: []Link{
					Link{
						Source: "foo->output0",
						Target: "bar->input0",
					},
				},
			},
			err: true,
		},
		{
			// Input aliases node ID
			input: Topology{
				Nodes: []Node{
					Node{
						ID:        "foo",
						Operation: "input",
						Inputs:    []string{"foo"},
						Outputs:   []string{},
					},
					Node{
						ID:        "bar",
						Operation: "output",
						Inputs:    []string{"bar->input0"},
						Outputs:   []string{},
					},
				},
				Links: []Link{
					Link{
						Source: "foo->output0",
						Target: "bar->input0",
					},
				},
			},
			err: true,
		},
		{
			// Input aliases node ID
			input: Topology{
				Nodes: []Node{
					Node{
						ID:        "foo",
						Operation: "input",
						Inputs:    []string{"bar"},
						Outputs:   []string{},
					},
					Node{
						ID:        "bar",
						Operation: "output",
						Inputs:    []string{"bar->input0"},
						Outputs:   []string{},
					},
				},
				Links: []Link{
					Link{
						Source: "foo->output0",
						Target: "bar->input0",
					},
				},
			},
			err: true,
		},
		{
			// Output aliases node ID
			input: Topology{
				Nodes: []Node{
					Node{
						ID:        "foo",
						Operation: "input",
						Outputs:   []string{"foo"},
						Inputs:    []string{},
					},
					Node{
						ID:        "bar",
						Operation: "output",
						Inputs:    []string{"bar->input0"},
						Outputs:   []string{},
					},
				},
				Links: []Link{
					Link{
						Source: "foo->output0",
						Target: "bar->input0",
					},
				},
			},
			err: true,
		},
		{
			// Output aliases node ID
			input: Topology{
				Nodes: []Node{
					Node{
						ID:        "foo",
						Operation: "input",
						Outputs:   []string{"bar"},
						Inputs:    []string{},
					},
					Node{
						ID:        "bar",
						Operation: "output",
						Inputs:    []string{"bar->input0"},
						Outputs:   []string{},
					},
				},
				Links: []Link{
					Link{
						Source: "foo->output0",
						Target: "bar->input0",
					},
				},
			},
			err: true,
		}, {
			// Link references nonexistant endpoint
			input: Topology{
				Nodes: []Node{
					Node{
						ID:        "foo",
						Operation: "input",
						Inputs:    []string{},
						Outputs:   []string{"foo->output0"},
					},
					Node{
						ID:        "bar",
						Operation: "output",
						Inputs:    []string{"bar->input0"},
						Outputs:   []string{},
					},
				},
				Links: []Link{
					Link{
						Source: "baz->output0",
						Target: "bar->input0",
					},
				},
			},
			err: true,
		}, {
			// Link in the wrong direction
			input: Topology{
				Nodes: []Node{
					Node{
						ID:        "foo",
						Operation: "input",
						Inputs:    []string{},
						Outputs:   []string{"foo->output0"},
					},
					Node{
						ID:        "bar",
						Operation: "output",
						Inputs:    []string{"bar->input0"},
						Outputs:   []string{},
					},
				},
				Links: []Link{
					Link{
						Source: "bar->input0",
						Target: "foo->output0",
					},
				},
			},
			err: true,
		},
	}

	for i, v := range cases {
		t.Logf("Test case %d: err=%v input=%v", i, v.err, v.input)
		res := ValidateTopology(v.input)
		t.Logf("\tres=%v", res)
		if v.err && res == nil {
			t.Errorf("\tTopology should have errored but did not")
		}
		if !v.err && res != nil {
			t.Errorf("\tTopology should not have errored but did")
			t.Errorf("\tTopology should have errored but did not")
		}

	}
}

func TestGetEffectiveDimensions(t *testing.T) {
	tnx, err := FromJSON(samples.SampleMLP3Layer())

	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		id        string
		dim       *[]int
		shoulderr bool
	}{
		{"foo", nil, true},
		{"input->output0", &[]int{25}, false},
		{"hidden1<-input0", &[]int{25}, false},
		{"hidden1->output0", &[]int{25}, false},
	}

	for n, c := range cases {
		actual, err := tnx.GetEffectiveDimensions(c.id)

		t.Logf("Test case %d, IOID '%s':", n, c.id)

		if c.shoulderr {
			if err == nil {
				t.Errorf("\tshould have error-ed and did not")
			}
			continue
		}

		if err != nil {
			t.Errorf("\tunexpected error %v", err)
			continue
		}

		t.Logf("\texpected dimensions: %v", *c.dim)
		t.Logf("\tactual dimensions: %v", *actual)

		if !cmp.Equal(actual, c.dim) {
			t.Errorf("\tactual and expected dimensions did not match")
		}

	}
}
