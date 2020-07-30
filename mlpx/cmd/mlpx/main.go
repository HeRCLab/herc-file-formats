package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/akamensky/argparse"

	"github.com/herclab/herc-file-formats/mlpx/go/mlpx"

	"github.com/mattn/go-isatty"
)

func randrange(l []float64) float64 {
	min := l[0]
	max := l[1]
	if max < min {
		min, max = max, min
	}
	return min + rand.Float64()*(max-min)
}

func main() {
	parser := argparse.NewParser("mlpx", "MLP eXchange utility")

	versionFlag := parser.Flag("v", "version", &argparse.Options{Help: "Display version number and exit."})

	randomSeed := parser.String("S", "seed", &argparse.Options{
		Help:    "Specify a random seed as an integer, or '-' for the current system time",
		Default: "-",
	})

	//// new command //////////////////////////////////////////////////////

	newCommand := parser.NewCommand("new", "Generate a new MLPX with random weights and biases.")

	newLayerSizes := newCommand.IntList("s", "sizes",
		&argparse.Options{
			Help:    "Specify a list of integer layer sizes, in neurons",
			Default: []int{5, 10, 5},
			Validate: func(args []string) error {
				for i, v := range args {
					s, err := strconv.Atoi(v)
					if err != nil {
						return fmt.Errorf("Could not parse layer size %d ('%s') due to error: %v", i, v, err)
					}

					if s < 1 {
						return fmt.Errorf("Layer size %d ('%s') is too small, must have at least one neuron", i, v)
					}
				}

				return nil
			},
		})

	defaultLayerActivations := true

	newLayerActivations := newCommand.StringList("a", "activations",
		&argparse.Options{
			Help: "Specify a list of activation function strings for each layer. Takes precedence over --default_activation",
			Validate: func(args []string) error {
				defaultLayerActivations = true
				return nil
			},
		})

	newLayerDefaultActivation := newCommand.String("A", "default_activation",
		&argparse.Options{
			Help:    "Specify an activation function to be used for all layers. If used in conjunction with --activations, this argument is ignored",
			Default: "sigmoid",
		})

	newBiasRange := newCommand.FloatList("b", "bias_range",
		&argparse.Options{
			Help:    "Specify the upper and lower bounds for the random bias values.",
			Default: []float64{0, 1.0},
		})

	newWeightRange := newCommand.FloatList("w", "weight_range",
		&argparse.Options{
			Help:    "Specify the upper and lower bounds for the random weight values.",
			Default: []float64{0, 1.0},
		})

	newAlpha := newCommand.Float("p", "alpha",
		&argparse.Options{
			Help:    "Specify the learning rate.",
			Default: 0.1,
		})

	newOutput := newCommand.String("o", "output",
		&argparse.Options{
			Default: "-",
			Help:    "Specify the output at which the generated MLPX should be stored, or '-' for standard out.",
		})

	//// validate command /////////////////////////////////////////////////

	validateCommand := parser.NewCommand("validate", "Validate an existing MLPX file from disk.")

	validateInput := validateCommand.String("i", "input",
		&argparse.Options{
			Help:    "Specify the MLPX file to validate, or '-' for standard input.",
			Default: "-",
		})

	//// diff command /////////////////////////////////////////////////////

	diffCommand := parser.NewCommand("diff", "Compare two existing MLPX files")

	diffInput := diffCommand.StringList("i", "input",
		&argparse.Options{
			Help:     "Specify two MLPX files to diff.",
			Required: true,
		})

	diffIndent := diffCommand.String("I", "indent",
		&argparse.Options{
			Help:    "Specify the indent that should be used to show hierarchy.",
			Default: "\t",
		})

	diffEpsilon := diffCommand.Float("e", "epsilon",
		&argparse.Options{
			Help:    "Epsilon value used when comparing floating point numbers.",
			Default: 0.00001,
		})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Fprint(os.Stderr, parser.Usage(err))
		os.Exit(1)
	}

	if *versionFlag {
		fmt.Printf("mlpx v0.0.1-git\n")
		os.Exit(0)
	}

	if *randomSeed == "-" {
		rand.Seed(time.Now().UTC().UnixNano())
	} else {
		s, err := strconv.Atoi(*randomSeed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid random seed '%s': %v\n", *randomSeed, err)
			os.Exit(1)
		}
		rand.Seed(int64(s))
	}

	if newCommand.Happened() {
		activations := []string{}
		if defaultLayerActivations {
			for i := 0; i < len(*newLayerSizes); i++ {
				activations = append(activations, *newLayerDefaultActivation)
			}
		} else {
			activations = append(activations, *newLayerActivations...)
		}

		if len(*newLayerSizes) != len(activations) {
			fmt.Fprintf(os.Stderr, "Must specify same number of layer sizes and activation functions")
			os.Exit(1)
		}

		if len(*newLayerSizes) < 2 {
			fmt.Fprintf(os.Stderr, "Must specify at least two layers")
			os.Exit(1)
		}

		if len(*newBiasRange) != 2 {
			fmt.Fprintf(os.Stderr, "Must specify an upper and lower bias range")
			os.Exit(1)
		}

		if len(*newWeightRange) != 2 {
			fmt.Fprintf(os.Stderr, "Must specify an upper and lower weight range")
			os.Exit(1)
		}

		m := mlpx.MakeMLPX()

		// create our initializer snapshot
		snapid := m.NextSnapshotID()
		m.MustMakeSnapshot(snapid, *newAlpha)
		snap := m.Snapshots[snapid]

		for index, neurons := range *newLayerSizes {
			if neurons < 1 {
				fmt.Fprintf(os.Stderr, "Layer %d must have at least one neuron, requested %d", index, neurons)
				os.Exit(1)
			}

			layerID := fmt.Sprintf("%d", index)
			pred := fmt.Sprintf("%d", index-1)
			succ := fmt.Sprintf("%d", index+1)

			if index == 0 {
				layerID = "input"
				pred = ""
				succ = "1"

			} else if index == len(*newLayerSizes)-1 {
				layerID = "output"
				succ = ""

			}
			if index == 1 {
				pred = "input"
			}

			if index == len(*newLayerSizes)-2 {
				succ = "output"
			}

			snap.MustMakeLayer(layerID, neurons, pred, succ)
			layer := snap.Layers[layerID]

			// initialize weights and biases
			layer.EnsureWeights()
			layer.EnsureBiases()

			for i := 0; i < len(*layer.Weights); i++ {
				(*layer.Weights)[i] = randrange(*newWeightRange)
			}

			for i := 0; i < len(*layer.Biases); i++ {
				(*layer.Biases)[i] = randrange(*newBiasRange)
			}

			layer.ActivationFunction = activations[index]
		}

		if *newOutput == "-" {
			data, err := m.ToJSON()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to generate JSON: %v\n", err)
				os.Exit(1)
			}

			fmt.Print(string(data))
			fmt.Print("")
		} else {
			err := m.WriteJSON(*newOutput)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write output: %v\n", err)
				os.Exit(1)
			}
		}

	} else if validateCommand.Happened() {
		data := []byte{}
		if *validateInput == "-" {
			if isatty.IsTerminal(os.Stdin.Fd()) {
				fmt.Fprintf(os.Stderr, "Reading input from standard in.\n")
			}

			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
				os.Exit(1)
			}

		} else {
			data, err = ioutil.ReadFile(*validateInput)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
				os.Exit(1)
			}
		}

		m, err := mlpx.FromJSON(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse input: %v\n", err)
			os.Exit(1)
		}

		err = m.Validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Validation failed, error was: %v\n", err)
			os.Exit(2)
		}

		os.Exit(0)

	} else if diffCommand.Happened() {
		if len(*diffInput) != 2 {
			fmt.Fprintf(os.Stderr, "Must specify exactly two files to diff")
			os.Exit(1)
		}

		m1, err := mlpx.ReadJSON((*diffInput)[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read MLPX file '%s': %v", (*diffInput)[0], err)
			os.Exit(1)
		}

		m2, err := mlpx.ReadJSON((*diffInput)[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read MLPX file '%s': %v", (*diffInput)[1], err)
			os.Exit(1)
		}

		diffs := m1.Diff(m2, *diffIndent, *diffEpsilon)
		for _, d := range diffs {
			fmt.Println(d)
		}

		os.Exit(len(diffs))
	}

}