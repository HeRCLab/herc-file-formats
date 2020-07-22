package mlpx

import (
	"testing"
)

func TestValidate(t *testing.T) {
	m1 := getTestMLPX1()

	err := m1.Validate()
	if err != nil {
		t.Error(err)
	}

	// test bad schema version
	m1.Schema[1] = 1
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with invalid schema number, but didn't")
	}
	m1.Schema[1] = 0

	// test bad schema type
	m1.Schema[0] = "foo"
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with invalid schema type, but didn't")
	}
	m1.Schema[0] = "mlpx"

	// test invalid reference to predecessor
	m1.Snapshots["0"].Layers["hidden0"].Predecessor = "foo"
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with invalid predecessor reference, but didn't")
	}
	m1.Snapshots["0"].Layers["hidden0"].Predecessor = "input"

	// test invalid reference to successor
	m1.Snapshots["0"].Layers["hidden0"].Successor = "foo"
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with invalid predecessor reference, but didn't")
	}
	m1.Snapshots["0"].Layers["hidden0"].Successor = "output"

	// test invalid weight list size
	m1.Snapshots["0"].Layers["hidden0"].Weights = &[]float64{1}
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with invalid weight length, but didn't")
	}
	m1.Snapshots["0"].Layers["hidden0"].Weights = &[]float64{1, 2, 3, 4}

	// test invalid output list size
	m1.Snapshots["0"].Layers["hidden0"].Outputs = &[]float64{1}
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with invalid outputs length, but didn't")
	}
	m1.Snapshots["0"].Layers["hidden0"].Outputs = &[]float64{1, 2}

	// test invalid activation list size
	m1.Snapshots["0"].Layers["hidden0"].Activations = &[]float64{1}
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with invalid activation length, but didn't")
	}
	m1.Snapshots["0"].Layers["hidden0"].Activations = &[]float64{1, 2}

	// test invalid deltas list size
	m1.Snapshots["0"].Layers["hidden0"].Deltas = &[]float64{1}
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with invalid deltas length, but didn't")
	}
	m1.Snapshots["0"].Layers["hidden0"].Deltas = &[]float64{1, 2}

	// test invalid bias list size
	m1.Snapshots["0"].Layers["hidden0"].Biases = &[]float64{1}
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with invalid bias length, but didn't")
	}
	m1.Snapshots["0"].Layers["hidden0"].Biases = &[]float64{1, 2}

	// test non-isomorphic snapshots -- different lengths of snapshot lists
	// case
	err = m1.MakeSnapshot("1")
	if err != nil {
		t.Error(err)
	}
	err = m1.Snapshots["1"].MakeLayer("input", 2, "", "output")
	if err != nil {
		t.Error(err)
	}
	err = m1.Snapshots["1"].MakeLayer("output", 2, "input", "")
	if err != nil {
		t.Error(err)
	}
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with non-isomorphic layers")
	}

	// make sure we can detect mismatched neuron counts
	delete(m1.Snapshots, "1")
	err = m1.MakeIsomorphicSnapshot("1", "0")
	if err != nil {
		t.Error(err)
	}
	m1.Snapshots["1"].Layers["hidden0"].Neurons = 5
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with isomorphic layers with non-matching neuron counts")
	}
	m1.Snapshots["1"].Layers["hidden0"].Neurons = 2

	// test an invalid topology
	m1.Snapshots["1"].Layers["hidden0"].Successor = "hidden0"
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with topology errors")
	}
	m1.Snapshots["1"].Layers["hidden0"].Predecessor = "hidden0"
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with topology errors")
	}
	m1.Snapshots["1"].Layers["input"].Successor = "input"
	m1.Snapshots["1"].Layers["input"].Predecessor = "input"
	m1.Snapshots["1"].Layers["output"].Successor = "output"
	m1.Snapshots["1"].Layers["output"].Predecessor = "output"
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with topology errors")
	}

}
