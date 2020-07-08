package schema

import (
	"fmt"
)

// This file contains lookup logic for querying a TNX schema object. Because
// validation requires many lookups by various criteria, we implement simple
// caching to make sure the average lookup time is O(1).

// prepareLookupCaches ensures that the lookup cache fields of a tnx object
// have been properly initialized.
func (tnx *TNX) prepareLookupCaches() {
	if tnx.linkLookupCache == nil {
		tnx.linkLookupCache = make(map[string]*Link)
	}

	if tnx.nodeIOLookupCache == nil {
		tnx.nodeIOLookupCache = make(map[string]*Node)
	}

	if tnx.nodeLookupCache == nil {
		tnx.nodeLookupCache = make(map[string]*Node)
	}

}

// lookupNodeByID retrieves a node matching the given ID. It will return an
// error if either the ID does not exist, or there is no node with the matching
// ID.
func (tnx *TNX) lookupNodeByID(id string) (*Node, error) {
	tnx.prepareLookupCaches()

	n, ok := tnx.nodeLookupCache[id]
	if ok {
		return n, nil
	}

	for _, n := range tnx.Topology.Nodes {
		if n.ID == id {
			tnx.nodeLookupCache[id] = &n
			return &n, nil
		}
	}

	return nil, fmt.Errorf("No such node with id '%s', either the ID does not exist, or does not refer to a node", id)
}

// lookupNodeByIOID retrieves a node with the matching input or output ID
func (tnx *TNX) lookupNodeByIOID(searchID string) (*Node, error) {
	tnx.prepareLookupCaches()

	n, ok := tnx.nodeIOLookupCache[searchID]
	if ok {
		return n, nil
	}

	for _, n := range tnx.Topology.Nodes {
		for _, id := range n.Inputs {
			if id == searchID {
				tnx.nodeIOLookupCache[searchID] = &n
				return &n, nil
			}
		}

		for _, id := range n.Outputs {
			if id == searchID {
				tnx.nodeIOLookupCache[searchID] = &n
				return &n, nil
			}
		}
	}
	return nil, fmt.Errorf("No node with an input or output with ID '%s' found", searchID)
}

// lookupLinkByEndpoint retrieves a link where either endpoint is exactly
// equal to the specified search ID
func (tnx *TNX) lookupLinkByEndpoint(searchID string) (*Link, error) {
	tnx.prepareLookupCaches()

	l, ok := tnx.linkLookupCache[searchID]
	if ok {
		return l, nil
	}

	for _, l := range tnx.Topology.Links {
		if (l.Source == searchID) || (l.Target == searchID) {
			tnx.linkLookupCache[searchID] = &l
			return &l, nil
		}
	}

	return nil, fmt.Errorf("No link with source or target ID '%s' found", searchID)
}
