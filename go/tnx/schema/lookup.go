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
		tnx.linkLookupCache = make(map[string][]*Link)
	}

	if tnx.nodeIOLookupCache == nil {
		tnx.nodeIOLookupCache = make(map[string]*Node)
	}

	if tnx.nodeLookupCache == nil {
		tnx.nodeLookupCache = make(map[string]*Node)
	}

}

// LookupNodeByID retrieves a node matching the given ID. It will return an
// error if either the ID does not exist, or there is no node with the matching
// ID.
func (tnx *TNX) LookupNodeByID(id string) (*Node, error) {
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

// LookupNodeByIOID retrieves a node with the matching input or output ID
func (tnx *TNX) LookupNodeByIOID(searchID string) (*Node, error) {
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

// IsInput returns true if and only if the searcID references an ID that
// exists within the given TNX and is an input to a node.
func (tnx *TNX) IsInput(searchID string) bool {
	node, err := tnx.LookupNodeByIOID(searchID)
	if err != nil {
		return false
	}

	for _, v := range node.Inputs {
		if searchID == v {
			return true
		}
	}

	return false
}

// IsOutput returns true if and only if the searcID references an ID that
// exists within the given TNX and is an output from a node.
func (tnx *TNX) IsOutput(searchID string) bool {
	node, err := tnx.LookupNodeByIOID(searchID)
	if err != nil {
		return false
	}

	for _, v := range node.Outputs {
		if searchID == v {
			return true
		}
	}

	return false
}

// IsIO returns true if and only if searchID references an ID that exists
// within the given TNX and is either an input or an output from a node.
func (tnx *TNX) IsIO(searchID string) bool {
	_, err := tnx.LookupNodeByIOID(searchID)
	return err == nil
}

// LookupAdjacent returns all nodes which are adjacent to the given ID.
//
// If the searchID is a node, then it will return all nodes which are connected
// to ANY output of the specified node, not including the node
// searchID specifies.
//
// If the searchID is an input or an output, then it will return all nodes
// which are reachable via a link that has the given ID as an endpoint, not
// including the node to which searchID is attached.
func (tnx *TNX) LookupAdjacent(searchID string) ([]*Node, error) {
	adjacent := make([]*Node, 0)

	if tnx.IsIO(searchID) {
		links, err := tnx.LookupLinkByEndpoint(searchID)
		if err != nil {
			return nil, err
		}

		for _, l := range links {
			if l.Source != searchID {
				n, err := tnx.LookupNodeByIOID(l.Source)
				if err != nil {
					return nil, err
				}
				adjacent = append(adjacent, n)
			}

			if l.Target != searchID {
				n, err := tnx.LookupNodeByIOID(l.Target)
				if err != nil {
					return nil, err
				}
				adjacent = append(adjacent, n)
			}
		}
	} else { // this is a node
		node, err := tnx.LookupNodeByID(searchID)
		if err != nil {
			return nil, err
		}

		for _, oid := range node.Outputs {
			links, err := tnx.LookupLinkByEndpoint(oid)
			if err != nil {
				return nil, err
			}

			for _, l := range links {
				// NOTE: we assume that we are the source,
				// and the other end is the target. The schema
				// validation will catch it if this is not
				// the case, but if you want to use this
				// function on it's own, you should probably
				// run the validation first.
				//
				// TL;DR: this makes assumptions that are only
				// safe IF the TNX is valid.

				node, err := tnx.LookupNodeByIOID(l.Target)
				if err != nil {
					return nil, err
				}

				adjacent = append(adjacent, node)
			}
		}
	}

	return adjacent, nil

}

// LookupLinkByEndpoint retrieves a link where either endpoint is exactly
// equal to the specified search ID
func (tnx *TNX) LookupLinkByEndpoint(searchID string) ([]*Link, error) {
	tnx.prepareLookupCaches()

	links, ok := tnx.linkLookupCache[searchID]
	if ok {
		return links, nil
	}

	links = make([]*Link, 0)

	for _, l := range tnx.Topology.Links {
		if (l.Source == searchID) || (l.Target == searchID) {
			links = append(links, &l)
		}
	}

	tnx.linkLookupCache[searchID] = links

	return links, nil
}
