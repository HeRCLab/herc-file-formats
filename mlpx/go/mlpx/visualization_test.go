package mlpx

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAverageBias(t *testing.T) {
	m := MakeMLPX()
	m.MustMakeSnapshot("0", 0.1)
	m.Snapshots["0"].MustMakeLayer("input", 2, "", "hidden0")
	m.Snapshots["0"].MustMakeLayer("hidden0", 2, "input", "output")
	m.Snapshots["0"].MustMakeLayer("output", 2, "hidden0", "")

	m.Snapshots["0"].Layers["hidden0"].Weights = &[]float64{1.5, 2.5, 3.5, 4}
	m.Snapshots["0"].Layers["output"].Outputs = &[]float64{0.5, 1.4}
	m.Snapshots["0"].Layers["hidden0"].ActivationFunction = "foobar"
	m.Snapshots["0"].Layers["output"].Biases = &[]float64{1, 2}

	expect := [][]float64{[]float64{1.5}}

	if !cmp.Equal(expect, AverageBias([]*MLPX{m})) {
		t.Logf("expect=%v", expect)
		t.Logf("actual=%v", AverageBias([]*MLPX{m}))
		t.Errorf("Expected and actual AverageBias values do not match")
	}
}
