package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/kong"

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

// getDashDir exists because we use the convention of '-' to mean standard in
// or standard out for file argument. We want to use kong "path" type for
// these, so that path expansions will be applied. However since '-' is a
// non-relative path, Kong expands it to <cwd>/-. This function returns what
// that path will be in that event.
//
// This precludes the user from specifying a file named '-' in their CWD.
func getDashDir() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Sprintf("failed to get working directory: %v\n", err)
		os.Exit(1)
	}
	return fmt.Sprintf("%s/-", wd)
}

var CLI struct {
	New struct {
		Sizes             []int     `name:"sizes" short:"s" default:"5,10,5" help:"List of layer sizes in in neurons."`
		Activations       []string  `name:"activations" short:"a" help:"List of activation function strings for each layer. Takes precedence over --default_activation."`
		DefaultActivation string    `name:"default_activation" short:"A" default:"sigmoid" help:"Activation function to be used for all layers. If used in conjunction with --activations, this argument is ignored."`
		BiasRange         []float64 `name:"bias_range" short:"b" default:"0.0,1.0" help:"Lower and upper bound (from left to right) for the random bias values."`
		WeightRange       []float64 `name:"weight_range" short:"w" default:"0.0,1.0" help:"Lower and upper bound (from left to right) for the random weight values. (Default: 0 1)"`
		Alpha             float64   `name:"alpha" short:"p" default:"0.001" help:"Learning rate."`
		Output            string    `name:"output" short:"o" type:"path" default:"-" help:"Output file to which the generated MLPX will be written. Specify '-' for standard output."`
	} `cmd help:"Generate a new MLPX file."`

	Validate struct {
		Input string `arg name:"input" short:"i" type:"path" default:"-" help:"Input MLPX file to validate, or '-' for standard input."`
	} `cmd help:"Validate an existing MLPX file."`

	Summarize struct {
		Input  string `arg name:"input" short:"i" type:"path" default:"-" help:"Input MLPX file to summarize, or '-' for standard input."`
		Indent string `name:"indent" default:"\t" short"I" help:"Specify the indent that should be used to show hierarchy."`
	} `cmd help:"Summarize an existing MLPX file."`

	Diff struct {
		Base    string  `arg required type:"path" help:"Path to an MLPX file which is used as the baseline for the comparison."`
		Other   string  `arg required type:"path" help:"Path to an MLPX file which will be compared to the baseline."`
		Indent  string  `name:"indent" default:"\t" short"I" help:"Specify the indent that should be used to show hierarchy."`
		Epsilon float64 `name:"epsilon" short:"e" default:"0.00001" help:"Epsilon value to use when comparing floating point numbers."`
	} `cmd help:"Compare two existing MLPX files."`

	PlotBias struct {
		Inputs []string `arg type:"existingpath" help:"Input files to plot."`
	} `cmd help:"Plot average bias across one or more MLPX files over time."`

	Seed string `name:"seed" short:"S" default:"-" help:"Seed for random number generator, as an integer. You may wish to set this if you want to reproducible generate the same MLPX multiple times. Use '-' for the current system time."`

	Version bool `name:"version" short:"V" default:"false" help:"Display version and exit"`
}

func main() {

	ctx := kong.Parse(&CLI)

	if CLI.Version {
		fmt.Printf("mlpx v0.0.2\n")
		os.Exit(0)
	}

	if CLI.Seed == "-" {
		rand.Seed(time.Now().UTC().UnixNano())
	} else {
		s, err := strconv.Atoi(CLI.Seed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid random seed '%s': %v\n", CLI.Seed, err)
			os.Exit(1)
		}
		rand.Seed(int64(s))
	}

	if ctx.Command() == "new" {
		activations := []string{}
		if len(CLI.New.Activations) == 0 {
			for i := 0; i < len(CLI.New.Sizes); i++ {
				activations = append(activations, CLI.New.DefaultActivation)
			}
		} else {
			activations = append(activations, CLI.New.Activations...)
		}

		if len(CLI.New.Sizes) != len(activations) {
			fmt.Fprintf(os.Stderr, "Must specify same number of layer sizes and activation functions\n")
			os.Exit(1)
		}

		if len(CLI.New.Sizes) < 2 {
			fmt.Fprintf(os.Stderr, "Must specify at least two layers\n")
			os.Exit(1)
		}

		if len(CLI.New.BiasRange) != 2 {
			fmt.Fprintf(os.Stderr, "Must specify an upper and lower bias range\n")
			os.Exit(1)
		}

		if len(CLI.New.WeightRange) != 2 {
			fmt.Fprintf(os.Stderr, "Must specify an upper and lower weight range\n")
			os.Exit(1)
		}

		m := mlpx.MakeMLPX()

		// create our initializer snapshot
		snapid := m.NextSnapshotID()
		m.MustMakeSnapshot(snapid, CLI.New.Alpha)
		snap := m.Snapshots[snapid]

		for index, neurons := range CLI.New.Sizes {
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

			} else if index == len(CLI.New.Sizes)-1 {
				layerID = "output"
				succ = ""

			}
			if index == 1 {
				pred = "input"
			}

			if index == len(CLI.New.Sizes)-2 {
				succ = "output"
			}

			snap.MustMakeLayer(layerID, neurons, pred, succ)
			layer := snap.Layers[layerID]

			// initialize weights and biases
			layer.EnsureWeights()
			layer.EnsureBiases()

			for i := 0; i < len(*layer.Weights); i++ {
				(*layer.Weights)[i] = randrange(CLI.New.WeightRange)
			}

			for i := 0; i < len(*layer.Biases); i++ {
				(*layer.Biases)[i] = randrange(CLI.New.BiasRange)
			}

			layer.ActivationFunction = activations[index]
		}

		if CLI.New.Output == getDashDir() {
			data, err := m.ToJSON()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to generate JSON: %v\n", err)
				os.Exit(1)
			}

			fmt.Print(string(data))
			fmt.Print("")
		} else {
			err := m.WriteJSON(CLI.New.Output)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write output: %v\n", err)
				os.Exit(1)
			}
		}

	} else if (ctx.Command() == "validate") || (ctx.Command() == "validate <input>") {
		data := []byte{}

		if CLI.Validate.Input == getDashDir() {
			if isatty.IsTerminal(os.Stdin.Fd()) {
				fmt.Fprintf(os.Stderr, "Reading input from standard in.\n")
			}

			var err error
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
				os.Exit(1)
			}

		} else {
			var err error
			data, err = ioutil.ReadFile(CLI.Validate.Input)
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

	} else if (ctx.Command() == "summarize") || (ctx.Command() == "summarize <input>") {
		data := []byte{}

		if CLI.Summarize.Input == getDashDir() {
			if isatty.IsTerminal(os.Stdin.Fd()) {
				fmt.Fprintf(os.Stderr, "Reading input from standard in.\n")
			}

			var err error
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
				os.Exit(1)
			}

		} else {
			var err error
			data, err = ioutil.ReadFile(CLI.Summarize.Input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
				os.Exit(1)
			}
		}

		// This makes \t actually turn into a tab character
		expanded, err := strconv.Unquote(fmt.Sprintf("\"%s\"", CLI.Summarize.Indent))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to expand indent: '%s': %v\n", CLI.Summarize.Indent, err)
			os.Exit(1)
		}

		m, err := mlpx.FromJSON(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse input: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(m.Summarize(expanded))

		os.Exit(0)

	} else if ctx.Command() == "diff <base> <other>" {
		m1, err := mlpx.ReadJSON(CLI.Diff.Base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read MLPX file '%s': %v\n", CLI.Diff.Base, err)
			os.Exit(1)
		}

		m2, err := mlpx.ReadJSON(CLI.Diff.Other)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read MLPX file '%s': %v\n", CLI.Diff.Other, err)
			os.Exit(1)
		}

		// This makes \t actually turn into a tab character
		expanded, err := strconv.Unquote(fmt.Sprintf("\"%s\"", CLI.Diff.Indent))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to expand indent: '%s': %v\n", CLI.Diff.Indent, err)
			os.Exit(1)
		}

		diffs := m1.Diff(m2, expanded, CLI.Diff.Epsilon)
		for _, d := range diffs {
			fmt.Println(d)
		}

		os.Exit(len(diffs))

	} else if ctx.Command() == "plot-bias" || ctx.Command() == "plot-bias <inputs>" {
		mlps := []*mlpx.MLPX{}

		for i, ipath := range CLI.PlotBias.Inputs {
			m, err := mlpx.ReadJSON(ipath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load %d-th input file '%s'\n", i, ipath)
				os.Exit(1)
			}

			mlps = append(mlps, m)
		}

		err := mlpx.PlotAverageBias(mlps)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error plotting biases: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)

	} else {
		fmt.Fprintf(os.Stderr, "Don't understand how to parse that command.\n")
		panic(ctx.Command())
	}

}
