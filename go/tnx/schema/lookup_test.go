package schema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
	n1, err := tnx.lookupNodeByID("foo")
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
	n2, err := tnx.lookupNodeByID("foo")
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
	n3, err := tnx.lookupNodeByID("baz")
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
	n1, err := tnx.lookupNodeByIOID("foo->output0")
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
	n2, err := tnx.lookupNodeByIOID("foo->output0")
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
	n3, err := tnx.lookupNodeByIOID("baz")
	if (err == nil) || (n3 != nil) {
		t.Errorf("should have failed to lookup nonexistant node")
	}

	// should not get any results for a node ID that isn't an IO
	n4, err := tnx.lookupNodeByIOID("foo")
	if (err == nil) || (n4 != nil) {
		t.Errorf("should have failed to lookup node ID (rather than IOID)")
	}
}

func TestLookupLinkByEndpoint(t *testing.T) {
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

	l1, err := tnx.lookupLinkByEndpoint("foo->output0")
	if (err != nil) || (l1 == nil) {
		t.Errorf("Should not have failed to lookup foo->output0")
	}
	t.Logf("l1: %v", l1)

	l2, err := tnx.lookupLinkByEndpoint("foo->output0")
	if (err != nil) || (l2 == nil) {
		t.Errorf("Should not have failed to lookup foo->output0")
	}
	t.Logf("l2: %v", l2)

	l3, err := tnx.lookupLinkByEndpoint("bar<-input0")
	if (err != nil) || (l3 == nil) {
		t.Errorf("Should not have failed to lookup bar<-input0")
	}
	t.Logf("l3: %v", l3)

	l4, err := tnx.lookupLinkByEndpoint("bar<-input0")
	if (err != nil) || (l4 == nil) {
		t.Errorf("Should not have failed to lookup bar<-input0")
	}
	t.Logf("l4: %v", l4)

	if !cmp.Equal(l1, l2) || !cmp.Equal(l2, l3) || !cmp.Equal(l3, l4) {
		t.Errorf("Looking up the link should not return different results")
	}

	l5, err := tnx.lookupLinkByEndpoint("foo")
	if (l5 != nil) || (err == nil) {
		t.Errorf("Should not be able to look up link by node ID using lookupLinkByEndpoint")
	}
}
