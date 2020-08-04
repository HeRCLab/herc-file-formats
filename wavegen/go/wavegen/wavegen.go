// Package wavegen is used for generating synthetic test data for the Dropbear
// LSTM project.
package wavegen

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/montanaflynn/stats"

	"github.com/guptarohit/asciigraph"

	"gonum.org/v1/plot/plotter"
)

// WaveParameters is used to store the parameters that generate a particular
// wave. See GenerateSyntheticData().
type WaveParameters struct {
	// SampleRate is the sample rate at which the wave should be generated,
	// in Hz
	SampleRate float64

	// Offset is the time at which samples should begin being collected, in
	// seconds
	Offset float64

	// Duration is the number of seconds of samples that should be generated
	Duration float64

	// Frequencies is the list of frequencies of Sin wave that should be
	// generated
	Frequencies []float64

	// Phases is the list of phases of the Sin waves that should be
	// generated
	Phases []float64

	// Amplitudes is the list of amplitudes of Sin that should be generated
	Amplitudes []float64

	// Noises stores a list of noise functions which are applied on a
	// per-signal basis.  if this field is left empty, then no noise will
	// be generated for the given signal. These noise functions are
	// understood:
	//
	// * "pseudo" (rand.Float64)
	// * "none" (no noise)
	// * "" (no noise)
	Noises []string

	// NoiseMagnitues is a list of coefficients to the given noise function
	// for a particular signal.  If empty, it is assumed that all
	// magnitudes are 1.0.
	NoiseMagnitudes []float64

	// GlobalNoise accepts the same values as Noises, but are applied
	// globally rather than to a specific signal.
	GlobalNoise string

	// GlobalNoiseMagnitude works similarly to NoiseMagnitudes, but applies
	// to the global noise.
	GlobalNoiseMagnitude float64
}

func (w *WaveParameters) Summarize() (string, error) {
	err := w.ValidateParameters()
	if err != nil {
		return "", err
	}

	s := "SYNTHETIC WAVE PARAMETERS SUMMARY:\n\n"
	s = fmt.Sprintf("%s\tSample Rate . . . . . %f\n", s, w.SampleRate)
	s = fmt.Sprintf("%s\tOffset  . . . . . . . %fs\n", s, w.Offset)
	s = fmt.Sprintf("%s\tDuration  . . . . . . %fs\n", s, w.Duration)
	s = fmt.Sprintf("%s\tGlobal Noise  . . . . %s\n", s, w.GlobalNoise)
	s = fmt.Sprintf("%s\t|Global Noise|  . . . %f\n", s, w.GlobalNoiseMagnitude)
	s = fmt.Sprintf("%s\n\tCOMPONENTS:\n", s)

	for i := range w.Frequencies {
		if w.Noises[i] == "" || w.Noises[i] == "none" {
			s = fmt.Sprintf("%s\t\t%f × Sin(2 × π × %f × t + %f)\n",
				s, w.Amplitudes[i], w.Frequencies[i], w.Phases[i])
		} else {
			s = fmt.Sprintf("%s\t\t%f × Sin(2 × π × %f × t + %f) + %f × %s()\n",
				s, w.Amplitudes[i], w.Frequencies[i], w.Phases[i], w.NoiseMagnitudes[i], w.Noises[i])
		}
	}

	return s, nil
}

// Signal represents a time-series signal
//
// When generating or modifying a signal, you must guarantee that the T and S
// values are sorted by T (e.g. earlier points in time have lower indices), and
// that T ans S are the same size.
type Signal struct {
	// signal values
	S []float64 `json:"samples"`

	T []float64 `json:"times"`

	// sample rate in Hz, note that this is set by the caller, this library
	// cannot guarantee the accuracy of this value
	SampleRate float64
}

// Declare conformity with plotter.XYer
var _ plotter.XYer = (*Signal)(nil)

// Len implements gonum's plotter.XYer interface
func (s *Signal) Len() int {
	return s.Size()
}

// XY implements gonum's plotter.XYer interface
func (s *Signal) XY(i int) (float64, float64) {
	return s.T[i], s.T[i]
}

func (s *Signal) Summarize() (string, error) {
	err := s.ValidateIndex(0)
	if err != nil {
		return "", err
	}

	mean, err := stats.Mean(s.S)
	if err != nil {
		return "", err
	}

	median, err := stats.Median(s.S)
	if err != nil {
		return "", err
	}

	min, err := stats.Min(s.S)
	if err != nil {
		return "", err
	}

	max, err := stats.Max(s.S)
	if err != nil {
		return "", err
	}

	stdev, err := stats.StandardDeviation(s.S)
	if err != nil {
		return "", err
	}

	str := "SIGNAL DATA SUMMARY:\n\n"
	str = fmt.Sprintf("%s\t# of Samples . . . . . . %d\n", str, len(s.T))
	str = fmt.Sprintf("%s\tReported Sample Rate . . %f\n", str, s.SampleRate)
	str = fmt.Sprintf("%s\tAverage Sample Rate  . . %f\n", str, s.AverageSampleRate())
	str = fmt.Sprintf("%s\tDuration . . . . . . . . %fs\n", str, s.Duration())
	str = fmt.Sprintf("%s\tMean . . . . . . . . . . %f\n", str, mean)
	str = fmt.Sprintf("%s\tMedian . . . . . . . . . %f\n", str, median)
	str = fmt.Sprintf("%s\tStandard Deviation . . . %f\n", str, stdev)
	str = fmt.Sprintf("%s\tMin  . . . . . . . . . . %f\n", str, min)
	str = fmt.Sprintf("%s\tMax  . . . . . . . . . . %f\n", str, max)

	str = fmt.Sprintf("%s\n\tSIGNAL DATA OVERVIEW:\n\n", str)

	graph := asciigraph.Plot(s.S, asciigraph.Height(30), asciigraph.Width(80-16), asciigraph.Offset(16))

	str = fmt.Sprintf("%s%s\n", str, graph)

	return str, nil
}

// Sample represents a single sample from a Signal
type Sample struct {

	// The time component of the sample
	T float64

	// The value component of the sample
	S float64
}

type SampleList []Sample

// Declare conformity with plotter.XYer
var _ plotter.XYer = (SampleList)(nil)

// Len implements gonum's plotter.XYer interface
func (l SampleList) Len() int {
	return len(l)
}

// XY implements gonum's plotter.XYer interface
func (l SampleList) XY(i int) (float64, float64) {
	return l[i].T, l[i].S
}

// ValidateIndex ensures that the specified index can be retrieved from the
// signal, and generates an error if not.
func (s *Signal) ValidateIndex(i int) error {
	if len(s.S) != len(s.T) {
		return fmt.Errorf("Signal is corrupt, S array is length %d, but T array is length %T",
			len(s.S), len(s.T))
	}

	if i < 0 || i >= s.Size() {
		return fmt.Errorf("Index %d out of bounds for signal of length %d",
			i, s.Size())
	}

	return nil
}

// AverageSampleRate calculates the average sample rate in Hz of the signal.
// This can be useful when dealing with real signals, rather than those
// generated using wavegen.
func (s *Signal) AverageSampleRate() float64 {
	dt := 0.0
	count := 0
	for i := range s.T {
		if i == 0 {
			continue
		}

		dt += s.T[i] - s.T[i-1]
		count++
	}

	return 1.0 / (dt / float64(count))

}

// Size returns the number of samples which a Signal contains.
func (s *Signal) Size() int {
	return len(s.S)
}

// Duration returns the number seconds of data in the signal.
func (s *Signal) Duration() float64 {
	return s.T[s.Size()-1] - s.T[0]
}

// Index retrieves the ith sample value.
func (s *Signal) Index(i int) (Sample, error) {
	err := s.ValidateIndex(i)
	if err != nil {
		return Sample{}, err
	}

	return Sample{T: s.T[i], S: s.S[i]}, nil
}

// MustIndex works identically to Index(), but calls panic() if an error occurs
func (s *Signal) MustIndex(i int) Sample {
	sample, err := s.Index(i)
	if err != nil {
		panic(err)
	}
	return sample
}

// NearestIndex will return the index within a signal which has a time value
// as close as possible to the specified time argument. It will return
// an index with a greater or equal value if overshoot is true, and a lesser
// or equal value if overshoot is false. It will return 0 if the time value
// is before the beginning of the signal, and s.Size()-1 if the time value
// is after the end of the signal.
//
// Note that time includes the offset, so time=0 when the offset is 1 second
// will return 0, since it is before the beginning of the signal data.
func (s *Signal) NearestIndex(time float64, overshoot bool) int {
	if s.Size() == 0 {
		return 0
	}

	if time < s.T[0] {
		return 0
	}

	i := 0
	for {
		t := s.T[i]

		if t == time {
			return i
		} else if overshoot && (t >= time) {
			return i
		} else if (t != time) && (i+1 >= s.Size()) {
			// prevent the next case from going out of bounds
			return s.Size() - 1
		} else if (!overshoot) && ((s.T[i+1]) > time) {
			return i
		}

		i++
	}
}

// Interpolate can be used to perform linear interpolation. It will compute a
// separate linear interpolation for each point in times, and return an
// appropriate sample for each.
func (s *Signal) Interpolate(times ...float64) []Sample {
	interpolated := make([]Sample, len(times))
	for i, t := range times {
		i0 := s.NearestIndex(t, false) // undershoot
		i1 := s.NearestIndex(t, true)  // overshoot

		s0 := s.MustIndex(i0)
		s1 := s.MustIndex(i1)

		// In the case where the requested value is exactly equal to a
		// known value, we can omit the interpolation and return it
		// directly. This would work equivalently with t1 and s1.
		if s0.T == t && s1.T == t {
			interpolated[i] = Sample{T: t, S: s0.S}
			continue
		}

		dist := math.Abs(s0.T - s1.T)
		d0 := math.Abs(s0.T - t)
		d1 := math.Abs(s1.T - t)

		if dist == 0 {
			interpolated[i] = Sample{T: t, S: s0.S}
			continue
		}

		interpolated[i] = Sample{T: t, S: (d0*s1.S)/(d0+d1) + (d1*s0.S)/(d0+d1)}
	}

	return interpolated
}

// ValidateParameters will ensure that the parameters are valid.
//
// If the noises or noise magnitudes are empty, then they will be filled to
// an appropriate length with default values.
func (w *WaveParameters) ValidateParameters() error {
	// if omitted, assume no noise is desired
	if len(w.Noises) == 0 {
		w.Noises = make([]string, len(w.Frequencies))
		for i := 0; i < len(w.Frequencies); i++ {
			w.Noises[i] = "none"
		}
	}

	// if omitted, assume magnitudes of 1.0
	if len(w.NoiseMagnitudes) == 0 {
		w.NoiseMagnitudes = make([]float64, len(w.Frequencies))
		for i := 0; i < len(w.Frequencies); i++ {
			w.NoiseMagnitudes[i] = 1.0
		}
	}

	if len(w.Noises) != len(w.Frequencies) {
		return fmt.Errorf("Invalid parameters: length of noises does not match length of frequencies")
	}

	if len(w.Phases) != len(w.Frequencies) {
		return fmt.Errorf("Invalid parameters: length of phases does not match length of frequencies")
	}

	if len(w.Amplitudes) != len(w.Frequencies) {
		return fmt.Errorf("Invalid parameters: length of amplitudes does not match length of frequencies")
	}

	if len(w.NoiseMagnitudes) != len(w.Frequencies) {
		return fmt.Errorf("Invalid parameters: length of noise magnitudes does not match length of frequencies")
	}

	return nil
}

// Noise generates a randomized noise value for the index-th component of
// the wave parameter. An index of -1 indicates that the global noise should
// be generated instead.
func (w *WaveParameters) Noise(index int) (float64, error) {
	err := w.ValidateParameters()
	if err != nil {
		return 0, err
	}

	if index < -1 || index >= len(w.Frequencies) {
		return 0, fmt.Errorf("Index %d out of bound for parameters of %d components", index, len(w.Frequencies))
	}

	kind := ""
	if index == -1 {
		kind = w.GlobalNoise
	} else {
		kind = w.Noises[index]
	}

	v := 0.0
	switch kind {
	case "pseudo":
		v = rand.Float64()
	case "none":
		v = 0
	case "":
		v = 0
	default:
		return 0, fmt.Errorf("Unknown noise kind '%s'", kind)

	}

	if index == -1 {
		return w.GlobalNoiseMagnitude * v, nil
	}
	return w.NoiseMagnitudes[index] * v, nil
}

// GenerateSyntheticData generates a signal which is a composition of several
// Sin functions of the given frequencies, phases, and amplitudes, with noise
// optionally applied to each signal, and optionally applied to the data
// overall.
func (w *WaveParameters) GenerateSyntheticData() (*Signal, error) {
	// number of points to generate
	points := int(math.Ceil(w.SampleRate * w.Duration))
	samplePeriod := 1.0 / w.SampleRate

	err := w.ValidateParameters()
	if err != nil {
		return nil, err
	}

	sig := &Signal{
		T:          make([]float64, points),
		S:          make([]float64, points),
		SampleRate: w.SampleRate,
	}

	for i := 0; i < points; i++ {
		sig.S[i] = 0
		sig.T[i] = samplePeriod * float64(i)
		for j, freq := range w.Frequencies {

			// component noise
			n, err := w.Noise(j)
			if err != nil {
				return nil, err
			}

			sig.S[i] += w.Amplitudes[j]*math.Sin(2*math.Pi*freq*sig.T[i]+w.Phases[j]) + n
		}

		// global noise
		n, err := w.Noise(-1)
		if err != nil {
			return nil, err
		}
		sig.S[i] = sig.S[i] + n
	}

	return sig, nil
}
