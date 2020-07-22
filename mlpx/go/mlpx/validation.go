package mlpx

import (
	"fmt"
)

// validateReferences checks the MLP for referential integrity
func (mlp *MLPX) validateReferences() error {

	for snapid, snapshot := range mlp.Snapshots {
		if snapid != snapshot.ID {
			return fmt.Errorf("snapshot '%s': snapshot ID in struct and in MLPX table do not match, your MLPX implementation has a bug ('%s' =/= '%s')", snapid, snapid, snapshot.ID)
		}
		if snapshot.Parent != mlp {
			return fmt.Errorf("snapshot '%s': parent pointer does not reference the parent MLPX object, your MLPX implementation has bugs", snapid)
		}

		for layerid, layer := range snapshot.Layers {
			if layer.Parent != snapshot {
				return fmt.Errorf("snapshot '%s', layer '%s': parent pointer does not reference the parent snapshot object, your MLPX implementation has bugs", snapid, layerid)
			}

			if layerid != layer.ID {
				return fmt.Errorf("snapshot '%s', layer '%s': snapshot ID in struct and in MLPX table do not match, your MLPX implementation has a bug ('%s' =/= '%s')",
					snapid, layerid, layerid, layer.ID)
			}

			// verify integrity of predecessor references
			if layerid != "input" {
				// input layers don't have predecessors
				if _, ok := snapshot.Layers[layer.Predecessor]; !ok {
					return fmt.Errorf("snapshot '%s', layer '%s': predecessor '%s' references nonexistent layer",
						snapid, layerid, layer.Predecessor)
				}
			}

			// verify integrity of successor references
			if layerid != "output" {
				// output layers don't have successors
				if _, ok := snapshot.Layers[layer.Successor]; !ok {
					return fmt.Errorf("snapshot '%s', layer '%s': successor '%s' references nonexistent layer",
						snapid, layerid, layer.Predecessor)
				}
			}

			// verify size of weights list
			if layerid != "input" && layer.Weights != nil {
				expect := layer.Neurons * snapshot.Layers[layer.Predecessor].Neurons
				if len(*layer.Weights) != expect {
					return fmt.Errorf("snapshot '%s', layer '%s': weights array of length %d, should be %d",
						snapid, layerid, len(*layer.Weights), expect)
				}
			}

			// verify size of outputs list
			if layer.Outputs != nil {
				if len(*layer.Outputs) != layer.Neurons {
					return fmt.Errorf("snapshot '%s', layer '%s': output array of length %d, should be %d",
						snapid, layerid, len(*layer.Outputs), layer.Neurons)
				}
			}

			// verify size of activation list
			if layer.Activations != nil {
				if len(*layer.Activations) != layer.Neurons {
					return fmt.Errorf("snapshot '%s', layer '%s': activation array of length %d, should be %d",
						snapid, layerid, len(*layer.Activations), layer.Neurons)
				}
			}

			// verify size of deltas list
			if layer.Deltas != nil {
				if len(*layer.Deltas) != layer.Neurons {
					return fmt.Errorf("snapshot '%s', layer '%s': delta array of length %d, should be %d",
						snapid, layerid, len(*layer.Deltas), layer.Neurons)
				}
			}

			// verify size of bias list
			if layer.Biases != nil {
				if len(*layer.Biases) != layer.Neurons {
					return fmt.Errorf("snapshot '%s', layer '%s': bias array of length %d, should be %d",
						snapid, layerid, len(*layer.Biases), layer.Neurons)
				}
			}
		}
	}

	return nil
}

// validateIsomorphism checks that all snapshots in the MLP are isomorphic
// and have the same neuron counts.
func (mlp *MLPX) validateIsomorphism() error {
	// if we have made it this far, then we know that each snapshot is
	// internally valid, so all we have left to do is verify that all
	// layers have the same topology.
	if len(mlp.Snapshots) < 2 {
		// if there are no snapshots, there is nothing to do, likewise
		// if there is only one snapshot, we can safely assume it
		// is isomorphic to itself
		return nil
	}

	// we choose a "key" snapshot which we will compare with all the others
	keyid := ""
	for k := range mlp.Snapshots {
		keyid = k
		break
	}
	key := mlp.Snapshots[keyid]

	// we now compare each other snapshot against our chosen key
	for snapid, snapshot := range mlp.Snapshots {

		// make sure all layers in the key are in the snapshot
		for layerid := range key.Layers {
			if _, ok := snapshot.Layers[layerid]; !ok {
				return fmt.Errorf("snapshot '%s' and '%s' are not isomorphic, snapshot '%s' has layer ID '%s, but snapshot '%s' does not",
					keyid, snapid, keyid, layerid, snapid)
			}
		}

		// and the converse
		for layerid := range snapshot.Layers {
			if _, ok := key.Layers[layerid]; !ok {
				return fmt.Errorf("snapshot '%s' and '%s' are not isomorphic, snapshot '%s' has layer ID '%s, but snapshot '%s' does not",
					snapid, keyid, snapid, layerid, keyid)
			}
		}

		// finally, check the neuron counts, which also imply the
		// other member fields, and we have already validated they
		// are sized in a way appropriate for their neuron counts
		for layerid, layer := range key.Layers {
			if layer.Neurons != snapshot.Layers[layerid].Neurons {
				return fmt.Errorf("snapshot '%s' and '%s' have different numbers of neurons (%d, and %d respectively) for layer '%s'",
					keyid, snapid, layer.Neurons, snapshot.Layers[layerid].Neurons, layerid)
			}
		}
	}

	return nil
}

// validateTopology checks that the topology of each snapshot is such that
// the following facts hold true:
//
// * The number of in-edges for each non-input/output nodes is 1
// * The number of out-edges for each non-input/output nodes is 1
// * The number of in-edges for each input node is 0
// * The number of out-edges for each input node is 1
// * The number of in-edges for each output node is 1
// * The number of out-edges for each output node is 0
//
// NOTE: this function assumes that the MLP has correct referential integrity.
func (mlp *MLPX) validateTopology() error {

	for snapid, snapshot := range mlp.Snapshots {

		inEdges := make(map[string]int)
		outEdges := make(map[string]int)

		for layerid, layer := range snapshot.Layers {
			// make sure that the referenced predecessor and
			// successor also reference us
			if layerid != "output" { // successor case
				succpred := snapshot.Layers[layer.Successor].Predecessor
				if succpred != layerid {
					return fmt.Errorf("snapshot '%s', layer '%s': successor layer '%s' has a different predecessor '%s'",
						snapid, layerid, layer.Successor, succpred)
				}
			}

			if layerid != "input" { // predecessor case
				predsucc := snapshot.Layers[layer.Predecessor].Successor
				if predsucc != layerid {
					return fmt.Errorf("snapshot '%s', layer '%s': predecessor layer '%s' has a different successor '%s'",
						snapid, layerid, layer.Predecessor, predsucc)
				}
			}

			if layerid != "output" {
				inEdges[layer.Successor]++
			}

			if layerid != "input" {
				outEdges[layer.Predecessor]++
			}
		}

		for k, v := range inEdges {
			if k == "input" {
				if v != 0 {
					return fmt.Errorf("snapshot '%s', layer '%s': has wrong number of in-edges %d (expected 0)",
						snapid, k, v)
				}
				continue
			}
			if v != 1 {
				return fmt.Errorf("snapshot '%s', layer '%s': has wrong number of in-edges %d (expected 1)",
					snapid, k, v)
			}
		}

		for k, v := range outEdges {
			if k == "output" {
				if v != 0 {
					return fmt.Errorf("snapshot '%s', layer '%s': has wrong number of out-edges %d (expected 0)",
						snapid, k, v)
				}
				continue
			}
			if v != 1 {
				return fmt.Errorf("snapshot '%s', layer '%s': has wrong number of out-edges %d (expected 1)",
					snapid, k, v)
			}
		}
	}

	return nil
}

// Validate checks the MLPX file for any errors. If none are found, it returns
// nil.
func (mlp *MLPX) Validate() error {

	version, err := mlp.Version()
	if err != nil {
		return err
	}

	if version != 0 {
		return fmt.Errorf("unknown version number %d", version)
	}

	err = mlp.validateReferences()
	if err != nil {
		return err
	}

	err = mlp.validateIsomorphism()
	if err != nil {
		return err
	}

	err = mlp.validateTopology()
	if err != nil {
		return err
	}

	return nil

}
