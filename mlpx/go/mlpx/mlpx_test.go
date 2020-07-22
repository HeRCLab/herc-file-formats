package mlpx

import (
	"github.com/google/go-cmp/cmp"

	"testing"

	"github.com/kr/pretty"
)

func getTestJSON1() string {
	return `
	{
		"schema": ["mlpx", 0],
		"snapshots": {
			"0": {
				"layers": {
					"input": {
						"successor": "hidden0",
						"neurons": 2
					},
					"hidden0": {
						"successor": "output",
						"predecessor": "input",
						"neurons": 2,
						"weights": [1.5, 2.5, 3.5, 4]

					},
					"output": {
						"predecessor": "hidden0",
						"neurons": 2,
						"outputs": [0.5, 1.4]
					}
				}
			}
		}
	}`
}

func getTestMLPX1() *MLPX {
	m := MakeMLPX()
	m.MakeSnapshot("0")
	m.Snapshots["0"].MakeLayer("input", 2, "", "hidden0")
	m.Snapshots["0"].MakeLayer("hidden0", 2, "input", "output")
	m.Snapshots["0"].MakeLayer("output", 2, "hidden0", "")

	m.Snapshots["0"].Layers["hidden0"].Weights = &[]float64{1.5, 2.5, 3.5, 4}
	m.Snapshots["0"].Layers["output"].Outputs = &[]float64{0.5, 1.4}

	return m
}

func TestEndcodeDecode(t *testing.T) {

	m1, err := FromJSON([]byte(getTestJSON1()))
	if err != nil {
		t.Fatal(err)
	}

	m2 := getTestMLPX1()

	t.Logf("m1=%s", pretty.Sprintf("%#v", m1))
	t.Logf("m2=%s", pretty.Sprintf("%#v", m2))

	if !cmp.Equal(m1, m2) {
		t.Errorf("Decoded JSON does not match expected value")
		for _, v := range pretty.Diff(m1, m2) {
			t.Logf(v)
		}
	}

	b1, err := m1.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	m3, err := FromJSON(b1)
	if err != nil {
		t.Error(err)
	}

	t.Logf("m3=%s", pretty.Sprintf("%#v", m3))

	if !cmp.Equal(m1, m3) {
		t.Errorf("Decoded JSON does not match expected value")
		for _, v := range pretty.Diff(m1, m3) {
			t.Logf(v)
		}
	}
}

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
	m1.Snapshots["1"] = &Snapshot{
		Layers: map[string]*Layer{
			"input": &Layer{
				Successor: "output",
				Neurons:   2,
			},
			"output": &Layer{
				Predecessor: "input",
				Neurons:     2,
			},
		},
	}
	err = m1.Validate()
	if err == nil {
		t.Errorf("Should have error-ed with non-isomorphic layers")
	}

	// make sure we can detect mismatched neuron counts
	delete(m1.Snapshots, "1")
	m1.MakeIsomorphicSnapshot("1", "0")
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

func TestMakeIsomorphicSnapshot(t *testing.T) {
	m := getTestMLPX1()

	// this should result in a valid MLPX, and we have already verified
	// that the MLPX validation logic is correct in a separate test.
	m.MakeIsomorphicSnapshot("1", "0")
	err := m.Validate()
	if err != nil {
		t.Error(err)
	}
}