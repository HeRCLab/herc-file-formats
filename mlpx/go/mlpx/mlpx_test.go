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
	err := m.MakeSnapshot("0")
	if err != nil {
		panic(err)
	}
	err = m.Snapshots["0"].MakeLayer("input", 2, "", "hidden0")
	if err != nil {
		panic(err)
	}
	err = m.Snapshots["0"].MakeLayer("hidden0", 2, "input", "output")
	if err != nil {
		panic(err)
	}
	err = m.Snapshots["0"].MakeLayer("output", 2, "hidden0", "")
	if err != nil {
		panic(err)
	}

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
