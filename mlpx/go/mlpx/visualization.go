package mlpx

import (
	"fmt"

	"github.com/kingishb/go-gnuplot"

	"github.com/montanaflynn/stats"
)

// AverageBias generates a list of average output layer biases across each
// snapshot for each of the MLPX objects.
//
// AverageBias[i][j] is the average output bias of the i-th MLPX during it's
// j-th snapshot.
//
// It is not guaranteed that all sub-lists will be the same length.
func AverageBias(mlps []*MLPX) [][]float64 {
	averageBiases := [][]float64{}

	for _, mlp := range mlps {
		averageBias := []float64{}
		for _, snapid := range mlp.SortedSnapshotIDs() {
			snapshot := mlp.Snapshots[snapid]
			layerIDs := snapshot.SortedLayerIDs()
			outputLayer := snapshot.Layers[layerIDs[len(layerIDs)-1]]
			if outputLayer.Biases == nil {
				continue
			}

			mean, err := stats.Mean(*outputLayer.Biases)
			if err != nil {
				continue
			}

			averageBias = append(averageBias, mean)
		}
		averageBiases = append(averageBiases, averageBias)
	}

	return averageBiases
}

// PlotAverageBias uses gnuplot to display the results from AverageBias
func PlotAverageBias(mlps []*MLPX) error {

	p, err := gnuplot.NewPlotter("", true, false)
	if err != nil {
		return err
	}

	err = p.Cmd("set terminal qt")
	if err != nil {
		return err
	}

	err = p.SetStyle("lines")
	if err != nil {
		return err
	}

	biases := AverageBias(mlps)

	for i := range mlps {
		err = p.PlotX(biases[i], fmt.Sprintf("MLPX %d", i))
		if err != nil {
			return err
		}
	}

	err = p.Cmd("pause 31540000")
	if err != nil {
		return err
	}

	return nil
}
