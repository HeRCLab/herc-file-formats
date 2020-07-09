package schema

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/herclab/tnx/go/tnx/schema/samples"
)

func TestLookupNodeByID(t *testing.T) {
	tnx := &TNX{
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

	// simple node ID lookup
	n1, err := tnx.LookupNodeByID("foo")
	if err != nil {
		t.Error(err)
	}
	if n1 == nil {
		t.Errorf("Lookup up node ID foo should not have returned nil")
	}
	if n1.ID != "foo" {
		t.Errorf("Lookup of node ID foo returned the wrong result")
	}

	// do it a second time to make sure there isn't any weirdness relating
	// to cacheing
	n2, err := tnx.LookupNodeByID("foo")
	if err != nil {
		t.Error(err)
	}
	if n2 == nil {
		t.Errorf("Lookup up node ID foo should not have returned nil")
	}
	if n2.ID != "foo" {
		t.Errorf("Lookup of node ID foo returned the wrong result")
	}
	if n1 != n2 {
		t.Errorf("Looking up the same node twice returned different results.")
	}

	// should not get any results for a nonexistant node
	n3, err := tnx.LookupNodeByID("baz")
	if (err == nil) || (n3 != nil) {
		t.Errorf("should have failed to lookup nonexistant node")
	}
}

func TestLookupNodeByIOID(t *testing.T) {
	tnx := &TNX{
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

	// simple node ID lookup
	n1, err := tnx.LookupNodeByIOID("foo->output0")
	if err != nil {
		t.Error(err)
	}
	if n1 == nil {
		t.Errorf("Lookup up node ID foo should not have returned nil")
	}
	if n1.ID != "foo" {
		t.Errorf("Lookup of node ID foo returned the wrong result")
	}

	// do it a second time to make sure there isn't any weirdness relating
	// to cacheing
	n2, err := tnx.LookupNodeByIOID("foo->output0")
	if err != nil {
		t.Error(err)
	}
	if n2 == nil {
		t.Errorf("Lookup up node ID foo should not have returned nil")
	}
	if n2.ID != "foo" {
		t.Errorf("Lookup of node ID foo returned the wrong result")
	}
	if n1 != n2 {
		t.Errorf("Looking up the same node twice returned different results.")
	}

	// should not get any results for a nonexistant node
	n3, err := tnx.LookupNodeByIOID("baz")
	if (err == nil) || (n3 != nil) {
		t.Errorf("should have failed to lookup nonexistant node")
	}

	// should not get any results for a node ID that isn't an IO
	n4, err := tnx.LookupNodeByIOID("foo")
	if (err == nil) || (n4 != nil) {
		t.Errorf("should have failed to lookup node ID (rather than IOID)")
	}
}

func TestLookupLinkByEndpoint(t *testing.T) {
	tnx, err := FromJSON(samples.SampleMLP3Layer())
	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		id     string
		expect []*Link
	}{
		{"input->output0", []*Link{&Link{Source: "input->output0", Target: "hidden1<-input0"}}},
		{"hidden1<-input0", []*Link{&Link{Source: "input->output0", Target: "hidden1<-input0"}}},
		{"hidden1->output0", []*Link{&Link{Source: "hidden1->output0", Target: "activation1->input0"}}},
		{"activation1<-input0", []*Link{&Link{Source: "hidden1->output0", Target: "activation1<-input0"}}},
	}

	for n, c := range cases {
		actual, err := tnx.LookupLinkByEndpoint(c.id)
		t.Logf("Case %d, IOID '%s':", n, c.id)
		if err != nil {
			t.Errorf("\tUnexpected error: %v", err)
			continue
		}

		// we sort everything, to guarantee that it will be safe
		// to in-order compare the actual and expected values.

		sort.Slice(actual, func(i, j int) bool {
			return actual[i].Source < actual[j].Source
		})

		sort.Slice(actual, func(i, j int) bool {
			return actual[i].Target < actual[j].Target
		})

		sort.Slice(c.expect, func(i, j int) bool {
			return c.expect[i].Source < c.expect[j].Source
		})

		sort.Slice(c.expect, func(i, j int) bool {
			return c.expect[i].Target < c.expect[j].Target
		})

		t.Logf("\tactual:")
		for _, v := range actual {
			t.Logf("\t\t%v\n", v)
		}
		t.Logf("\texpect:")
		for _, v := range c.expect {
			t.Logf("\t\t%v\n", v)
		}

		if !cmp.Equal(actual, c.expect) {
			t.Errorf("\tactual and expected values do not match")
		}
	}

}

func TestCheckingIDType(t *testing.T) {
	tnx := &TNX{
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

	r1 := tnx.IsInput("foo")
	if r1 != false {
		t.Errorf("IsInput found a node to be an input")
	}

	r2 := tnx.IsInput("foo->output0")
	if r2 != false {
		t.Errorf("IsInput found an output to be an input")
	}

	r3 := tnx.IsInput("bar<-input0")
	if r3 != true {
		t.Errorf("IsInput did not return true for an input")
	}

	r4 := tnx.IsOutput("bar<-input0")
	if r4 != false {
		t.Errorf("IsOutput returned true for an input")
	}

	r5 := tnx.IsOutput("bar")
	if r5 != false {
		t.Errorf("IsOutput returned true for a node")
	}

	r6 := tnx.IsOutput("foo->output0")
	if r6 != true {
		t.Errorf("IsOutput returned false for an output")
	}

	r7 := tnx.IsIO("foo")
	if r7 != false {
		t.Errorf("IsIO returned true for a node")
	}

	r8 := tnx.IsIO("foo->output0")
	if r8 != true {
		t.Errorf("IsIO returned false for an output")
	}

	r9 := tnx.IsIO("bar<-input0")
	if r9 != true {
		t.Errorf("IsIO returned false for an input")
	}

}

func TestLookupAdjacent(t *testing.T) {
	tnx := &TNX{
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

	a1, err := tnx.LookupAdjacent("foo")
	if err != nil {
		t.Error(err)
	}
	expect1 := []*Node{&tnx.Topology.Nodes[1]}
	if !cmp.Equal(a1, expect1) {
		t.Errorf("Expected LookupAdjacent to return %v, but got %v instead", expect1, a1)
	}

	a2, err := tnx.LookupAdjacent("foo->output0")
	if err != nil {
		t.Error(err)
	}
	expect2 := []*Node{&tnx.Topology.Nodes[1]}
	if !cmp.Equal(a2, expect2) {
		t.Errorf("Expected LookupAdjacent to return %v, but got %v instead", expect2, a2)
	}

	a3, err := tnx.LookupAdjacent("bar<-input0")
	if err != nil {
		t.Error(err)
	}
	expect3 := []*Node{&tnx.Topology.Nodes[0]}
	if !cmp.Equal(a3, expect3) {
		t.Errorf("Expected LookupAdjacent to return %v, but got %v instead", expect3, a3)
	}
}
