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
						"weights": [1.5, 2.5, 3.5, 4],
						"activation_function": "foobar"

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
	m.MustMakeSnapshot("0")
	m.Snapshots["0"].MustMakeLayer("input", 2, "", "hidden0")
	m.Snapshots["0"].MustMakeLayer("hidden0", 2, "input", "output")
	m.Snapshots["0"].MustMakeLayer("output", 2, "hidden0", "")

	m.Snapshots["0"].Layers["hidden0"].Weights = &[]float64{1.5, 2.5, 3.5, 4}
	m.Snapshots["0"].Layers["output"].Outputs = &[]float64{0.5, 1.4}
	m.Snapshots["0"].Layers["hidden0"].ActivationFunction = "foobar"

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
	}
}

func TestMakeIsomorphicSnapshot(t *testing.T) {
	m := getTestMLPX1()

	// this should result in a valid MLPX, and we have already verified
	// that the MLPX validation logic is correct in a separate test.
	err := m.MakeIsomorphicSnapshot("1", "0")
	if err != nil {
		t.Error(err)
	}

	err = m.Validate()
	if err != nil {
		t.Error(err)
	}
}

func TestMLPX1Valid(t *testing.T) {
	m := getTestMLPX1()
	err := m.Validate()
	if err != nil {
		t.Log(err)
		t.Fatalf("Test MLPX 1 is not valid, this means that other tests will almost certainly result in false negatives")
	}
}

func TestSortedSnapshotIDs(t *testing.T) {
	m := getTestMLPX1()

	m.MustMakeIsomorphicSnapshot("1", "0")
	m.MustMakeIsomorphicSnapshot("initializer", "0")
	m.MustMakeIsomorphicSnapshot("02", "0")
	m.MustMakeIsomorphicSnapshot("foo", "0")
	m.MustMakeIsomorphicSnapshot("10", "0")
	m.MustMakeIsomorphicSnapshot("aaa", "0")

	// apparently sort.Slice is sometimes non-deterministic?!

	for i := 0; i < 1000; i++ {
		expect := []string{"initializer", "0", "1", "02", "10", "aaa", "foo"}

		sorted := m.SortedSnapshotIDs()
		if !cmp.Equal(sorted, expect) {
			t.Logf("Expected: %v", expect)
			t.Logf("Sorted: %v", sorted)
			t.Fatalf("Sort order incorrect!")
		}
	}

}

func TestSnapshotSucessor(t *testing.T) {
	m := getTestMLPX1()
	cases := []struct {
		input     string
		expectID  string
		shoulderr bool
	}{
		{"0", "1", false},
		{"1", "", true},
		{"2", "", true},
	}

	err := m.MakeIsomorphicSnapshot("1", "0")
	if err != nil {
		t.Error(err)
	}

	for i, c := range cases {
		res, err := m.Snapshots["0"].Successor(c.input)
		if err != nil {
			if !c.shoulderr {
				t.Errorf("Test case %d: %v should not have errored but did: %v", i, c, err)
			}
		} else if c.shoulderr {
			t.Errorf("Test case %d: %v should have errored but did not", i, c)
		}

		if res == nil {
			continue
		}

		if res.ID != c.expectID {
			t.Errorf("Test case '%d: %v, output '%s' did not match expected '%s'", i, c, res.ID, c.expectID)
		}
	}
}

func TestSnapshotPredecessor(t *testing.T) {
	m := getTestMLPX1()
	cases := []struct {
		input     string
		expectID  string
		shoulderr bool
	}{
		{"0", "", true},
		{"1", "0", false},
		{"2", "", true},
	}

	err := m.MakeIsomorphicSnapshot("1", "0")
	if err != nil {
		t.Error(err)
	}

	for i, c := range cases {
		res, err := m.Snapshots["0"].Predecessor(c.input)
		if err != nil {
			if !c.shoulderr {
				t.Errorf("Test case %d: %v should not have errored but did: %v", i, c, err)
			}
		} else if c.shoulderr {
			t.Errorf("Test case %d: %v should have errored but did not", i, c)
		}

		if res == nil {
			continue
		}

		if res.ID != c.expectID {
			t.Errorf("Test case '%d: %v, output '%s' did not match expected '%s'", i, c, res.ID, c.expectID)
		}
	}
}

func TestSortedLayerIDs(t *testing.T) {
	m := MakeMLPX()
	m.MustMakeSnapshot("0")
	m.Snapshots["0"].MustMakeLayer("input", 2, "", "hidden0")
	m.Snapshots["0"].MustMakeLayer("hidden0", 2, "input", "aaaa")
	m.Snapshots["0"].MustMakeLayer("aaaa", 2, "hidden0", "0000")
	m.Snapshots["0"].MustMakeLayer("0000", 2, "aaaa", "output")
	m.Snapshots["0"].MustMakeLayer("output", 2, "hidden0", "")

	layerids := m.Snapshots["0"].SortedLayerIDs()
	expect := []string{"input", "hidden0", "aaaa", "0000", "output"}

	t.Logf("sorted layer IDs: %v", layerids)
	t.Logf("expected layer IDs: %v", expect)

	if !cmp.Equal(layerids, expect) {
		t.Errorf("Sorted layer IDs were incorrect!")
	}
}
